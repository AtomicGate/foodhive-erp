package ar

import (
	"context"
	"fmt"
	"time"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

// ============================================
// Service Interface
// ============================================

type ARService interface {
	// Invoices
	CreateInvoice(ctx context.Context, req *models.CreateARInvoiceRequest, createdBy int) (int, error)
	GetInvoice(ctx context.Context, id int) (*models.ARInvoiceWithDetails, error)
	GetInvoiceByNumber(ctx context.Context, number string) (*models.ARInvoiceWithDetails, error)
	ListInvoices(ctx context.Context, filters *models.ARInvoiceListFilters) ([]models.ARInvoiceWithDetails, int64, error)
	PostInvoice(ctx context.Context, id int, postedBy int) error
	VoidInvoice(ctx context.Context, id int) error
	CreateFromOrder(ctx context.Context, orderID int, createdBy int) (int, error)

	// Payments
	CreatePayment(ctx context.Context, req *models.CreateARPaymentRequest, receivedBy int) (int, error)
	GetPayment(ctx context.Context, id int) (*models.ARPaymentWithDetails, error)
	ListPayments(ctx context.Context, customerID *int, limit int) ([]models.ARPaymentWithDetails, error)

	// Credit Management
	GetCustomerCredit(ctx context.Context, customerID int) (*models.CustomerCredit, error)
	CheckCreditAvailable(ctx context.Context, customerID int, amount float64) (bool, float64, error)
	UpdateCreditLimit(ctx context.Context, customerID int, newLimit float64) error

	// Aging
	GetCustomerAging(ctx context.Context, customerID int) (*models.CustomerAging, error)
	GetAgingReport(ctx context.Context) ([]models.CustomerAging, error)

	// Statement
	GetStatement(ctx context.Context, customerID int, fromDate, toDate string) (*models.CustomerStatement, error)

	// Overdue
	GetOverdueInvoices(ctx context.Context, daysOverdue int) ([]models.ARInvoiceWithDetails, error)
}

// ============================================
// Service Implementation
// ============================================

type arServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) ARService {
	return &arServiceImpl{db: db}
}

// ============================================
// Invoices
// ============================================

func (s *arServiceImpl) CreateInvoice(ctx context.Context, req *models.CreateARInvoiceRequest, createdBy int) (int, error) {
	invoiceNumber := s.generateInvoiceNumber(ctx)

	invDate, _ := time.Parse("2006-01-02", req.InvoiceDate)

	// Get customer payment terms for due date
	var paymentTerms int
	s.db.QueryRow(ctx, `SELECT COALESCE(payment_terms_days, 30) FROM customers WHERE id = $1`, req.CustomerID).Scan(&paymentTerms)

	var dueDate time.Time
	if req.DueDate != "" {
		dueDate, _ = time.Parse("2006-01-02", req.DueDate)
	} else {
		dueDate = invDate.AddDate(0, 0, paymentTerms)
	}

	// Calculate totals
	var subtotal float64
	for _, line := range req.Lines {
		subtotal += line.Quantity * line.UnitPrice
	}
	totalAmount := subtotal + req.TaxAmount + req.FreightAmount

	query := `
		INSERT INTO ar_invoices (
			invoice_number, customer_id, order_id, invoice_date, due_date, status,
			subtotal, tax_amount, freight_amount, total_amount, balance_due,
			currency, notes, created_by
		) VALUES ($1, $2, $3, $4, $5, 'DRAFT', $6, $7, $8, $9, $9, 'USD', $10, $11)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		invoiceNumber, req.CustomerID, req.OrderID, invDate, dueDate,
		subtotal, req.TaxAmount, req.FreightAmount, totalAmount, req.Notes, createdBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create AR invoice: %w", err)
	}

	// Insert lines
	for i, line := range req.Lines {
		lineTotal := line.Quantity * line.UnitPrice * (1 + line.TaxPercent/100)

		lineQuery := `
			INSERT INTO ar_invoice_lines (
				invoice_id, line_number, product_id, description, quantity,
				unit_price, tax_percent, line_total
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err := s.db.Exec(ctx, lineQuery,
			id, i+1, line.ProductID, line.Description, line.Quantity,
			line.UnitPrice, line.TaxPercent, lineTotal,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to create invoice line: %w", err)
		}
	}

	return id, nil
}

func (s *arServiceImpl) GetInvoice(ctx context.Context, id int) (*models.ARInvoiceWithDetails, error) {
	return s.getInvoice(ctx, "i.id = $1", id)
}

func (s *arServiceImpl) GetInvoiceByNumber(ctx context.Context, number string) (*models.ARInvoiceWithDetails, error) {
	return s.getInvoice(ctx, "i.invoice_number = $1", number)
}

func (s *arServiceImpl) getInvoice(ctx context.Context, whereClause string, arg interface{}) (*models.ARInvoiceWithDetails, error) {
	query := fmt.Sprintf(`
		SELECT i.id, i.invoice_number, i.customer_id, i.order_id, i.invoice_date, i.due_date,
			   i.status, i.subtotal, i.tax_amount, i.freight_amount, i.discount_amount,
			   i.total_amount, i.amount_paid, i.balance_due, i.currency, i.notes,
			   i.posted_by, i.posted_at, i.created_by, i.created_at, i.updated_at,
			   c.name as customer_name, c.customer_code,
			   COALESCE(so.order_number, '') as order_number,
			   GREATEST(0, CURRENT_DATE - i.due_date) as days_overdue
		FROM ar_invoices i
		JOIN customers c ON i.customer_id = c.id
		LEFT JOIN sales_orders so ON i.order_id = so.id
		WHERE %s`, whereClause)

	var inv models.ARInvoiceWithDetails
	var notes *string
	var postedAt *time.Time

	err := s.db.QueryRow(ctx, query, arg).Scan(
		&inv.Invoice.ID, &inv.Invoice.InvoiceNumber, &inv.Invoice.CustomerID, &inv.Invoice.OrderID,
		&inv.Invoice.InvoiceDate, &inv.Invoice.DueDate, &inv.Invoice.Status,
		&inv.Invoice.Subtotal, &inv.Invoice.TaxAmount, &inv.Invoice.FreightAmount,
		&inv.Invoice.DiscountAmount, &inv.Invoice.TotalAmount, &inv.Invoice.AmountPaid,
		&inv.Invoice.BalanceDue, &inv.Invoice.Currency, &notes, &inv.Invoice.PostedBy,
		&postedAt, &inv.Invoice.CreatedBy, &inv.Invoice.CreatedAt, &inv.Invoice.UpdatedAt,
		&inv.CustomerName, &inv.CustomerCode, &inv.OrderNumber, &inv.DaysOverdue,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	if notes != nil {
		inv.Invoice.Notes = *notes
	}

	// Get lines
	linesQuery := `
		SELECT id, invoice_id, line_number, product_id, description, quantity,
			   unit_price, tax_percent, line_total, order_line_id
		FROM ar_invoice_lines
		WHERE invoice_id = $1
		ORDER BY line_number`

	rows := s.db.Query(ctx, linesQuery, inv.Invoice.ID)
	defer rows.Close()

	for rows.Next() {
		var line models.ARInvoiceLine
		err := rows.Scan(&line.ID, &line.InvoiceID, &line.LineNumber, &line.ProductID,
			&line.Description, &line.Quantity, &line.UnitPrice, &line.TaxPercent,
			&line.LineTotal, &line.OrderLineID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invoice line: %w", err)
		}
		inv.Lines = append(inv.Lines, line)
	}

	return &inv, nil
}

func (s *arServiceImpl) ListInvoices(ctx context.Context, filters *models.ARInvoiceListFilters) ([]models.ARInvoiceWithDetails, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 20
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if filters.CustomerID != nil {
		whereClause += fmt.Sprintf(" AND i.customer_id = $%d", argNum)
		args = append(args, *filters.CustomerID)
		argNum++
	}
	if filters.Status != nil {
		whereClause += fmt.Sprintf(" AND i.status = $%d", argNum)
		args = append(args, *filters.Status)
		argNum++
	}
	if filters.DateFrom != "" {
		whereClause += fmt.Sprintf(" AND i.invoice_date >= $%d", argNum)
		args = append(args, filters.DateFrom)
		argNum++
	}
	if filters.DateTo != "" {
		whereClause += fmt.Sprintf(" AND i.invoice_date <= $%d", argNum)
		args = append(args, filters.DateTo)
		argNum++
	}
	if filters.Overdue {
		whereClause += " AND i.due_date < CURRENT_DATE AND i.balance_due > 0"
	}

	// Count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM ar_invoices i %s`, whereClause)
	var total int64
	s.db.QueryRow(ctx, countQuery, args...).Scan(&total)

	// Query
	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT i.id, i.invoice_number, i.customer_id, i.order_id, i.invoice_date, i.due_date,
			   i.status, i.subtotal, i.tax_amount, i.freight_amount, i.discount_amount,
			   i.total_amount, i.amount_paid, i.balance_due, i.currency, i.notes,
			   i.posted_by, i.posted_at, i.created_by, i.created_at, i.updated_at,
			   c.name as customer_name, c.customer_code,
			   COALESCE(so.order_number, '') as order_number,
			   GREATEST(0, CURRENT_DATE - i.due_date) as days_overdue
		FROM ar_invoices i
		JOIN customers c ON i.customer_id = c.id
		LEFT JOIN sales_orders so ON i.order_id = so.id
		%s
		ORDER BY i.invoice_date DESC, i.id DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var invoices []models.ARInvoiceWithDetails
	for rows.Next() {
		var inv models.ARInvoiceWithDetails
		var notes *string
		var postedAt *time.Time

		err := rows.Scan(
			&inv.Invoice.ID, &inv.Invoice.InvoiceNumber, &inv.Invoice.CustomerID, &inv.Invoice.OrderID,
			&inv.Invoice.InvoiceDate, &inv.Invoice.DueDate, &inv.Invoice.Status,
			&inv.Invoice.Subtotal, &inv.Invoice.TaxAmount, &inv.Invoice.FreightAmount,
			&inv.Invoice.DiscountAmount, &inv.Invoice.TotalAmount, &inv.Invoice.AmountPaid,
			&inv.Invoice.BalanceDue, &inv.Invoice.Currency, &notes, &inv.Invoice.PostedBy,
			&postedAt, &inv.Invoice.CreatedBy, &inv.Invoice.CreatedAt, &inv.Invoice.UpdatedAt,
			&inv.CustomerName, &inv.CustomerCode, &inv.OrderNumber, &inv.DaysOverdue,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan invoice: %w", err)
		}
		if notes != nil {
			inv.Invoice.Notes = *notes
		}
		invoices = append(invoices, inv)
	}

	return invoices, total, nil
}

func (s *arServiceImpl) PostInvoice(ctx context.Context, id int, postedBy int) error {
	result, err := s.db.Exec(ctx, `
		UPDATE ar_invoices SET status = 'POSTED', posted_by = $1, posted_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND status = 'DRAFT'`, postedBy, id)
	if err != nil {
		return fmt.Errorf("failed to post invoice: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("invoice not found or already posted")
	}

	// Update customer balance
	s.updateCustomerBalance(ctx, id, true)
	return nil
}

func (s *arServiceImpl) VoidInvoice(ctx context.Context, id int) error {
	// Get invoice info first
	var status string
	var customerID int
	s.db.QueryRow(ctx, `SELECT status, customer_id FROM ar_invoices WHERE id = $1`, id).Scan(&status, &customerID)

	if status == "VOID" {
		return fmt.Errorf("invoice already voided")
	}

	result, err := s.db.Exec(ctx, `
		UPDATE ar_invoices SET status = 'VOID', updated_at = NOW()
		WHERE id = $1 AND status != 'VOID'`, id)
	if err != nil {
		return fmt.Errorf("failed to void invoice: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("invoice not found")
	}

	// Reverse customer balance if was posted
	if status == "POSTED" || status == "PARTIAL" {
		s.updateCustomerBalance(ctx, id, false)
	}
	return nil
}

func (s *arServiceImpl) CreateFromOrder(ctx context.Context, orderID int, createdBy int) (int, error) {
	// Get order details
	var customerID int
	var subtotal, taxAmount, freightAmount float64
	err := s.db.QueryRow(ctx, `
		SELECT customer_id, subtotal, tax_amount, freight_amount
		FROM sales_orders WHERE id = $1 AND status IN ('SHIPPED', 'DELIVERED')`, orderID).Scan(
		&customerID, &subtotal, &taxAmount, &freightAmount)
	if err != nil {
		return 0, fmt.Errorf("order not found or not ready for invoicing")
	}

	// Create invoice request
	req := &models.CreateARInvoiceRequest{
		CustomerID:    customerID,
		OrderID:       &orderID,
		InvoiceDate:   time.Now().Format("2006-01-02"),
		TaxAmount:     taxAmount,
		FreightAmount: freightAmount,
	}

	// Get order lines
	rows := s.db.Query(ctx, `
		SELECT product_id, description, quantity_shipped, unit_price, 0 as tax_percent
		FROM sales_order_lines WHERE order_id = $1`, orderID)
	defer rows.Close()

	for rows.Next() {
		var line models.CreateARInvoiceLineReq
		rows.Scan(&line.ProductID, &line.Description, &line.Quantity, &line.UnitPrice, &line.TaxPercent)
		req.Lines = append(req.Lines, line)
	}

	return s.CreateInvoice(ctx, req, createdBy)
}

// ============================================
// Payments
// ============================================

func (s *arServiceImpl) CreatePayment(ctx context.Context, req *models.CreateARPaymentRequest, receivedBy int) (int, error) {
	receiptNumber := s.generateReceiptNumber(ctx)
	paymentDate, _ := time.Parse("2006-01-02", req.PaymentDate)

	query := `
		INSERT INTO ar_payments (
			receipt_number, customer_id, payment_date, payment_method, amount,
			currency, reference_no, check_number, notes, received_by
		) VALUES ($1, $2, $3, $4, $5, 'USD', $6, $7, $8, $9)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		receiptNumber, req.CustomerID, paymentDate, req.PaymentMethod, req.Amount,
		req.ReferenceNo, req.CheckNumber, req.Notes, receivedBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create payment: %w", err)
	}

	// Apply to invoices
	for _, app := range req.Applications {
		_, err := s.db.Exec(ctx, `
			INSERT INTO ar_payment_applications (payment_id, invoice_id, amount)
			VALUES ($1, $2, $3)`, id, app.InvoiceID, app.Amount)
		if err != nil {
			return 0, fmt.Errorf("failed to apply payment: %w", err)
		}

		// Update invoice
		s.db.Exec(ctx, `
			UPDATE ar_invoices SET
				amount_paid = amount_paid + $1,
				balance_due = balance_due - $1,
				status = CASE 
					WHEN balance_due - $1 <= 0 THEN 'PAID'::ar_invoice_status
					ELSE 'PARTIAL'::ar_invoice_status
				END,
				updated_at = NOW()
			WHERE id = $2`, app.Amount, app.InvoiceID)
	}

	// Update customer balance
	s.db.Exec(ctx, `UPDATE customers SET current_balance = current_balance - $1 WHERE id = $2`,
		req.Amount, req.CustomerID)

	return id, nil
}

func (s *arServiceImpl) GetPayment(ctx context.Context, id int) (*models.ARPaymentWithDetails, error) {
	query := `
		SELECT p.id, p.receipt_number, p.customer_id, p.payment_date, p.payment_method,
			   p.amount, p.currency, p.reference_no, p.check_number, p.bank_account,
			   p.notes, p.received_by, p.posted_at, p.created_at,
			   c.name as customer_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as receiver_name
		FROM ar_payments p
		JOIN customers c ON p.customer_id = c.id
		LEFT JOIN employees e ON p.received_by = e.id
		WHERE p.id = $1`

	var pay models.ARPaymentWithDetails
	var refNo, checkNo, bankAcc, notes *string
	var postedAt *time.Time

	err := s.db.QueryRow(ctx, query, id).Scan(
		&pay.Payment.ID, &pay.Payment.ReceiptNumber, &pay.Payment.CustomerID, &pay.Payment.PaymentDate,
		&pay.Payment.PaymentMethod, &pay.Payment.Amount, &pay.Payment.Currency,
		&refNo, &checkNo, &bankAcc, &notes, &pay.Payment.ReceivedBy, &postedAt, &pay.Payment.CreatedAt,
		&pay.CustomerName, &pay.ReceiverName,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	if refNo != nil {
		pay.Payment.ReferenceNo = *refNo
	}
	if checkNo != nil {
		pay.Payment.CheckNumber = *checkNo
	}
	if bankAcc != nil {
		pay.Payment.BankAccount = *bankAcc
	}
	if notes != nil {
		pay.Payment.Notes = *notes
	}

	// Get applications
	appRows := s.db.Query(ctx, `SELECT id, payment_id, invoice_id, amount FROM ar_payment_applications WHERE payment_id = $1`, id)
	defer appRows.Close()

	for appRows.Next() {
		var app models.ARPaymentApplication
		appRows.Scan(&app.ID, &app.PaymentID, &app.InvoiceID, &app.Amount)
		pay.Applications = append(pay.Applications, app)
	}

	return &pay, nil
}

func (s *arServiceImpl) ListPayments(ctx context.Context, customerID *int, limit int) ([]models.ARPaymentWithDetails, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if customerID != nil {
		whereClause += fmt.Sprintf(" AND p.customer_id = $%d", argNum)
		args = append(args, *customerID)
		argNum++
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.receipt_number, p.customer_id, p.payment_date, p.payment_method,
			   p.amount, p.currency, p.reference_no, p.check_number, p.bank_account,
			   p.notes, p.received_by, p.posted_at, p.created_at,
			   c.name as customer_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as receiver_name
		FROM ar_payments p
		JOIN customers c ON p.customer_id = c.id
		LEFT JOIN employees e ON p.received_by = e.id
		%s
		ORDER BY p.payment_date DESC
		LIMIT $%d`, whereClause, argNum)

	args = append(args, limit)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var payments []models.ARPaymentWithDetails
	for rows.Next() {
		var pay models.ARPaymentWithDetails
		var refNo, checkNo, bankAcc, notes *string
		var postedAt *time.Time

		err := rows.Scan(
			&pay.Payment.ID, &pay.Payment.ReceiptNumber, &pay.Payment.CustomerID, &pay.Payment.PaymentDate,
			&pay.Payment.PaymentMethod, &pay.Payment.Amount, &pay.Payment.Currency,
			&refNo, &checkNo, &bankAcc, &notes, &pay.Payment.ReceivedBy, &postedAt, &pay.Payment.CreatedAt,
			&pay.CustomerName, &pay.ReceiverName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		if refNo != nil {
			pay.Payment.ReferenceNo = *refNo
		}
		if checkNo != nil {
			pay.Payment.CheckNumber = *checkNo
		}
		if notes != nil {
			pay.Payment.Notes = *notes
		}
		payments = append(payments, pay)
	}

	return payments, nil
}

// ============================================
// Credit Management
// ============================================

func (s *arServiceImpl) GetCustomerCredit(ctx context.Context, customerID int) (*models.CustomerCredit, error) {
	query := `
		SELECT c.id, c.name, c.credit_limit, c.current_balance, c.payment_terms_days,
			   c.credit_limit - c.current_balance as available_credit,
			   COALESCE((SELECT SUM(balance_due) FROM ar_invoices WHERE customer_id = c.id AND due_date < CURRENT_DATE AND balance_due > 0), 0) as total_overdue,
			   COALESCE((SELECT MAX(CURRENT_DATE - due_date) FROM ar_invoices WHERE customer_id = c.id AND due_date < CURRENT_DATE AND balance_due > 0), 0) as oldest_overdue
		FROM customers c
		WHERE c.id = $1`

	var credit models.CustomerCredit
	err := s.db.QueryRow(ctx, query, customerID).Scan(
		&credit.CustomerID, &credit.CustomerName, &credit.CreditLimit, &credit.CurrentBalance,
		&credit.PaymentTermsDays, &credit.AvailableCredit, &credit.TotalOverdue, &credit.OldestOverdue,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer credit: %w", err)
	}

	// Determine credit status
	if credit.OldestOverdue > 90 {
		credit.CreditStatus = "HOLD"
	} else if credit.AvailableCredit <= 0 {
		credit.CreditStatus = "OVER_LIMIT"
	} else if credit.TotalOverdue > 0 {
		credit.CreditStatus = "OVERDUE"
	} else {
		credit.CreditStatus = "GOOD"
	}

	return &credit, nil
}

func (s *arServiceImpl) CheckCreditAvailable(ctx context.Context, customerID int, amount float64) (bool, float64, error) {
	credit, err := s.GetCustomerCredit(ctx, customerID)
	if err != nil {
		return false, 0, err
	}
	return credit.AvailableCredit >= amount, credit.AvailableCredit, nil
}

func (s *arServiceImpl) UpdateCreditLimit(ctx context.Context, customerID int, newLimit float64) error {
	result, err := s.db.Exec(ctx, `UPDATE customers SET credit_limit = $1 WHERE id = $2`, newLimit, customerID)
	if err != nil {
		return fmt.Errorf("failed to update credit limit: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("customer not found")
	}
	return nil
}

// ============================================
// Aging
// ============================================

func (s *arServiceImpl) GetCustomerAging(ctx context.Context, customerID int) (*models.CustomerAging, error) {
	query := `
		SELECT c.id, c.name, c.customer_code,
			   COALESCE(SUM(CASE WHEN i.due_date >= CURRENT_DATE THEN i.balance_due ELSE 0 END), 0) as current_amt,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 1 AND 30 THEN i.balance_due ELSE 0 END), 0) as days_1_30,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 31 AND 60 THEN i.balance_due ELSE 0 END), 0) as days_31_60,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 61 AND 90 THEN i.balance_due ELSE 0 END), 0) as days_61_90,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date > 90 THEN i.balance_due ELSE 0 END), 0) as over_90,
			   COALESCE(SUM(i.balance_due), 0) as total
		FROM customers c
		LEFT JOIN ar_invoices i ON c.id = i.customer_id AND i.status NOT IN ('DRAFT', 'VOID', 'PAID')
		WHERE c.id = $1
		GROUP BY c.id, c.name, c.customer_code`

	var aging models.CustomerAging
	err := s.db.QueryRow(ctx, query, customerID).Scan(
		&aging.CustomerID, &aging.CustomerName, &aging.CustomerCode,
		&aging.Aging.Current, &aging.Aging.Days1_30, &aging.Aging.Days31_60,
		&aging.Aging.Days61_90, &aging.Aging.Over90, &aging.Aging.Total,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer aging: %w", err)
	}

	return &aging, nil
}

func (s *arServiceImpl) GetAgingReport(ctx context.Context) ([]models.CustomerAging, error) {
	query := `
		SELECT c.id, c.name, c.customer_code,
			   COALESCE(SUM(CASE WHEN i.due_date >= CURRENT_DATE THEN i.balance_due ELSE 0 END), 0) as current_amt,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 1 AND 30 THEN i.balance_due ELSE 0 END), 0) as days_1_30,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 31 AND 60 THEN i.balance_due ELSE 0 END), 0) as days_31_60,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 61 AND 90 THEN i.balance_due ELSE 0 END), 0) as days_61_90,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date > 90 THEN i.balance_due ELSE 0 END), 0) as over_90,
			   COALESCE(SUM(i.balance_due), 0) as total
		FROM customers c
		LEFT JOIN ar_invoices i ON c.id = i.customer_id AND i.status NOT IN ('DRAFT', 'VOID', 'PAID')
		GROUP BY c.id, c.name, c.customer_code
		HAVING COALESCE(SUM(i.balance_due), 0) > 0
		ORDER BY total DESC`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var report []models.CustomerAging
	for rows.Next() {
		var aging models.CustomerAging
		err := rows.Scan(
			&aging.CustomerID, &aging.CustomerName, &aging.CustomerCode,
			&aging.Aging.Current, &aging.Aging.Days1_30, &aging.Aging.Days31_60,
			&aging.Aging.Days61_90, &aging.Aging.Over90, &aging.Aging.Total,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan aging: %w", err)
		}
		report = append(report, aging)
	}

	return report, nil
}

// ============================================
// Statement
// ============================================

func (s *arServiceImpl) GetStatement(ctx context.Context, customerID int, fromDate, toDate string) (*models.CustomerStatement, error) {
	// Get customer info
	var stmt models.CustomerStatement
	err := s.db.QueryRow(ctx, `SELECT id, name, customer_code FROM customers WHERE id = $1`, customerID).Scan(
		&stmt.CustomerID, &stmt.CustomerName, &stmt.CustomerCode,
	)
	if err != nil {
		return nil, fmt.Errorf("customer not found")
	}

	// Get opening balance
	s.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(CASE WHEN type = 'INVOICE' THEN amount ELSE -amount END), 0)
		FROM (
			SELECT 'INVOICE' as type, total_amount as amount FROM ar_invoices 
			WHERE customer_id = $1 AND invoice_date < $2 AND status NOT IN ('DRAFT', 'VOID')
			UNION ALL
			SELECT 'PAYMENT' as type, amount FROM ar_payments
			WHERE customer_id = $1 AND payment_date < $2
		) t`, customerID, fromDate).Scan(&stmt.OpeningBalance)

	// Get transactions
	query := `
		SELECT date, type, reference, description, debit, credit
		FROM (
			SELECT invoice_date as date, 'Invoice' as type, invoice_number as reference,
				   'Sales Invoice' as description, total_amount as debit, 0 as credit
			FROM ar_invoices
			WHERE customer_id = $1 AND invoice_date BETWEEN $2 AND $3 AND status NOT IN ('DRAFT', 'VOID')
			UNION ALL
			SELECT payment_date as date, 'Payment' as type, receipt_number as reference,
				   'Payment Received' as description, 0 as debit, amount as credit
			FROM ar_payments
			WHERE customer_id = $1 AND payment_date BETWEEN $2 AND $3
		) t
		ORDER BY date, type`

	rows := s.db.Query(ctx, query, customerID, fromDate, toDate)
	defer rows.Close()

	balance := stmt.OpeningBalance
	for rows.Next() {
		var line models.StatementLine
		rows.Scan(&line.Date, &line.Type, &line.Reference, &line.Description, &line.Debit, &line.Credit)
		balance = balance + line.Debit - line.Credit
		line.Balance = balance
		stmt.Lines = append(stmt.Lines, line)
	}

	stmt.ClosingBalance = balance
	return &stmt, nil
}

// ============================================
// Overdue
// ============================================

func (s *arServiceImpl) GetOverdueInvoices(ctx context.Context, daysOverdue int) ([]models.ARInvoiceWithDetails, error) {
	filters := &models.ARInvoiceListFilters{
		Overdue:  true,
		Page:     1,
		PageSize: 100,
	}
	invoices, _, err := s.ListInvoices(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Filter by days overdue
	var result []models.ARInvoiceWithDetails
	for _, inv := range invoices {
		if inv.DaysOverdue >= daysOverdue {
			result = append(result, inv)
		}
	}

	return result, nil
}

// ============================================
// Helpers
// ============================================

func (s *arServiceImpl) generateInvoiceNumber(ctx context.Context) string {
	var count int64
	s.db.QueryRow(ctx, `SELECT COUNT(*) FROM ar_invoices WHERE DATE(created_at) = CURRENT_DATE`).Scan(&count)
	return fmt.Sprintf("INV%s%04d", time.Now().Format("20060102"), count+1)
}

func (s *arServiceImpl) generateReceiptNumber(ctx context.Context) string {
	var count int64
	s.db.QueryRow(ctx, `SELECT COUNT(*) FROM ar_payments WHERE DATE(created_at) = CURRENT_DATE`).Scan(&count)
	return fmt.Sprintf("RCP%s%04d", time.Now().Format("20060102"), count+1)
}

func (s *arServiceImpl) updateCustomerBalance(ctx context.Context, invoiceID int, add bool) {
	var customerID int
	var amount float64
	s.db.QueryRow(ctx, `SELECT customer_id, total_amount FROM ar_invoices WHERE id = $1`, invoiceID).Scan(&customerID, &amount)

	if add {
		s.db.Exec(ctx, `UPDATE customers SET current_balance = current_balance + $1 WHERE id = $2`, amount, customerID)
	} else {
		s.db.Exec(ctx, `UPDATE customers SET current_balance = current_balance - $1 WHERE id = $2`, amount, customerID)
	}
}
