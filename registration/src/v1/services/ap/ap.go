package ap

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

type APService interface {
	// Invoices
	CreateInvoice(ctx context.Context, req *models.CreateAPInvoiceRequest, createdBy int) (int, error)
	GetInvoice(ctx context.Context, id int) (*models.APInvoiceWithDetails, error)
	GetInvoiceByNumber(ctx context.Context, vendorID int, number string) (*models.APInvoiceWithDetails, error)
	ListInvoices(ctx context.Context, filters *models.APInvoiceListFilters) ([]models.APInvoiceWithDetails, int64, error)
	ApproveInvoice(ctx context.Context, id int, approvedBy int) error
	VoidInvoice(ctx context.Context, id int) error
	CreateFromReceiving(ctx context.Context, receivingID int, createdBy int) (int, error)

	// Payments
	CreatePayment(ctx context.Context, req *models.CreateAPPaymentRequest, preparedBy int) (int, error)
	GetPayment(ctx context.Context, id int) (*models.APPaymentWithDetails, error)
	ListPayments(ctx context.Context, vendorID *int, limit int) ([]models.APPaymentWithDetails, error)
	VoidPayment(ctx context.Context, id int) error

	// Vendor Balance
	GetVendorBalance(ctx context.Context, vendorID int) (*models.VendorBalance, error)

	// Aging
	GetVendorAging(ctx context.Context, vendorID int) (*models.VendorAging, error)
	GetAgingReport(ctx context.Context) ([]models.VendorAging, error)

	// Due Bills
	GetDueInvoices(ctx context.Context, withinDays int) ([]models.APInvoiceWithDetails, error)
	GetOverdueInvoices(ctx context.Context, daysOverdue int) ([]models.APInvoiceWithDetails, error)
}

// ============================================
// Service Implementation
// ============================================

type apServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) APService {
	return &apServiceImpl{db: db}
}

// ============================================
// Invoices
// ============================================

func (s *apServiceImpl) CreateInvoice(ctx context.Context, req *models.CreateAPInvoiceRequest, createdBy int) (int, error) {
	invDate, _ := time.Parse("2006-01-02", req.InvoiceDate)

	// Get vendor payment terms for due date
	var paymentTerms int
	s.db.QueryRow(ctx, `SELECT COALESCE(payment_terms_days, 30) FROM vendors WHERE id = $1`, req.VendorID).Scan(&paymentTerms)

	var dueDate time.Time
	if req.DueDate != "" {
		dueDate, _ = time.Parse("2006-01-02", req.DueDate)
	} else {
		dueDate = invDate.AddDate(0, 0, paymentTerms)
	}

	// Calculate totals
	var subtotal float64
	for _, line := range req.Lines {
		subtotal += line.Quantity * line.UnitCost
	}
	totalAmount := subtotal + req.TaxAmount + req.FreightAmount

	query := `
		INSERT INTO ap_invoices (
			invoice_number, vendor_id, po_id, receiving_id, invoice_date, due_date, status,
			subtotal, tax_amount, freight_amount, total_amount, balance_due,
			currency, notes, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, 'PENDING', $7, $8, $9, $10, $10, 'USD', $11, $12)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.InvoiceNumber, req.VendorID, req.POID, req.ReceivingID, invDate, dueDate,
		subtotal, req.TaxAmount, req.FreightAmount, totalAmount, req.Notes, createdBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create AP invoice: %w", err)
	}

	// Insert lines
	for i, line := range req.Lines {
		lineTotal := line.Quantity * line.UnitCost * (1 + line.TaxPercent/100)

		lineQuery := `
			INSERT INTO ap_invoice_lines (
				invoice_id, line_number, product_id, description, quantity,
				unit_cost, tax_percent, line_total, gl_account_id
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

		_, err := s.db.Exec(ctx, lineQuery,
			id, i+1, line.ProductID, line.Description, line.Quantity,
			line.UnitCost, line.TaxPercent, lineTotal, line.GLAccountID,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to create invoice line: %w", err)
		}
	}

	return id, nil
}

func (s *apServiceImpl) GetInvoice(ctx context.Context, id int) (*models.APInvoiceWithDetails, error) {
	return s.getInvoice(ctx, "i.id = $1", id)
}

func (s *apServiceImpl) GetInvoiceByNumber(ctx context.Context, vendorID int, number string) (*models.APInvoiceWithDetails, error) {
	return s.getInvoice(ctx, "i.vendor_id = $1 AND i.invoice_number = $2", vendorID, number)
}

func (s *apServiceImpl) getInvoice(ctx context.Context, whereClause string, args ...interface{}) (*models.APInvoiceWithDetails, error) {
	query := fmt.Sprintf(`
		SELECT i.id, i.invoice_number, i.vendor_id, i.po_id, i.receiving_id, i.invoice_date, i.due_date,
			   i.status, i.subtotal, i.tax_amount, i.freight_amount, i.discount_amount,
			   i.total_amount, i.amount_paid, i.balance_due, i.currency, i.notes,
			   i.approved_by, i.approved_at, i.created_by, i.created_at, i.updated_at,
			   v.name as vendor_name, v.vendor_code,
			   COALESCE(po.po_number, '') as po_number,
			   GREATEST(0, CURRENT_DATE - i.due_date) as days_overdue
		FROM ap_invoices i
		JOIN vendors v ON i.vendor_id = v.id
		LEFT JOIN purchase_orders po ON i.po_id = po.id
		WHERE %s`, whereClause)

	var inv models.APInvoiceWithDetails
	var notes *string
	var approvedAt *time.Time

	err := s.db.QueryRow(ctx, query, args...).Scan(
		&inv.Invoice.ID, &inv.Invoice.InvoiceNumber, &inv.Invoice.VendorID, &inv.Invoice.POID,
		&inv.Invoice.ReceivingID, &inv.Invoice.InvoiceDate, &inv.Invoice.DueDate, &inv.Invoice.Status,
		&inv.Invoice.Subtotal, &inv.Invoice.TaxAmount, &inv.Invoice.FreightAmount,
		&inv.Invoice.DiscountAmount, &inv.Invoice.TotalAmount, &inv.Invoice.AmountPaid,
		&inv.Invoice.BalanceDue, &inv.Invoice.Currency, &notes, &inv.Invoice.ApprovedBy,
		&approvedAt, &inv.Invoice.CreatedBy, &inv.Invoice.CreatedAt, &inv.Invoice.UpdatedAt,
		&inv.VendorName, &inv.VendorCode, &inv.PONumber, &inv.DaysOverdue,
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
			   unit_cost, tax_percent, line_total, po_line_id, receiving_line_id, gl_account_id
		FROM ap_invoice_lines
		WHERE invoice_id = $1
		ORDER BY line_number`

	rows := s.db.Query(ctx, linesQuery, inv.Invoice.ID)
	defer rows.Close()

	for rows.Next() {
		var line models.APInvoiceLine
		err := rows.Scan(&line.ID, &line.InvoiceID, &line.LineNumber, &line.ProductID,
			&line.Description, &line.Quantity, &line.UnitCost, &line.TaxPercent,
			&line.LineTotal, &line.POLineID, &line.ReceivingLineID, &line.GLAccountID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invoice line: %w", err)
		}
		inv.Lines = append(inv.Lines, line)
	}

	return &inv, nil
}

func (s *apServiceImpl) ListInvoices(ctx context.Context, filters *models.APInvoiceListFilters) ([]models.APInvoiceWithDetails, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 20
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if filters.VendorID != nil {
		whereClause += fmt.Sprintf(" AND i.vendor_id = $%d", argNum)
		args = append(args, *filters.VendorID)
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
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM ap_invoices i %s`, whereClause)
	var total int64
	s.db.QueryRow(ctx, countQuery, args...).Scan(&total)

	// Query
	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT i.id, i.invoice_number, i.vendor_id, i.po_id, i.receiving_id, i.invoice_date, i.due_date,
			   i.status, i.subtotal, i.tax_amount, i.freight_amount, i.discount_amount,
			   i.total_amount, i.amount_paid, i.balance_due, i.currency, i.notes,
			   i.approved_by, i.approved_at, i.created_by, i.created_at, i.updated_at,
			   v.name as vendor_name, v.vendor_code,
			   COALESCE(po.po_number, '') as po_number,
			   GREATEST(0, CURRENT_DATE - i.due_date) as days_overdue
		FROM ap_invoices i
		JOIN vendors v ON i.vendor_id = v.id
		LEFT JOIN purchase_orders po ON i.po_id = po.id
		%s
		ORDER BY i.due_date ASC, i.id DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var invoices []models.APInvoiceWithDetails
	for rows.Next() {
		var inv models.APInvoiceWithDetails
		var notes *string
		var approvedAt *time.Time

		err := rows.Scan(
			&inv.Invoice.ID, &inv.Invoice.InvoiceNumber, &inv.Invoice.VendorID, &inv.Invoice.POID,
			&inv.Invoice.ReceivingID, &inv.Invoice.InvoiceDate, &inv.Invoice.DueDate, &inv.Invoice.Status,
			&inv.Invoice.Subtotal, &inv.Invoice.TaxAmount, &inv.Invoice.FreightAmount,
			&inv.Invoice.DiscountAmount, &inv.Invoice.TotalAmount, &inv.Invoice.AmountPaid,
			&inv.Invoice.BalanceDue, &inv.Invoice.Currency, &notes, &inv.Invoice.ApprovedBy,
			&approvedAt, &inv.Invoice.CreatedBy, &inv.Invoice.CreatedAt, &inv.Invoice.UpdatedAt,
			&inv.VendorName, &inv.VendorCode, &inv.PONumber, &inv.DaysOverdue,
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

func (s *apServiceImpl) ApproveInvoice(ctx context.Context, id int, approvedBy int) error {
	result, err := s.db.Exec(ctx, `
		UPDATE ap_invoices SET status = 'APPROVED', approved_by = $1, approved_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND status = 'PENDING'`, approvedBy, id)
	if err != nil {
		return fmt.Errorf("failed to approve invoice: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("invoice not found or already approved")
	}
	return nil
}

func (s *apServiceImpl) VoidInvoice(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx, `
		UPDATE ap_invoices SET status = 'VOID', updated_at = NOW()
		WHERE id = $1 AND status NOT IN ('PAID', 'VOID')`, id)
	if err != nil {
		return fmt.Errorf("failed to void invoice: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("invoice not found or cannot be voided")
	}
	return nil
}

func (s *apServiceImpl) CreateFromReceiving(ctx context.Context, receivingID int, createdBy int) (int, error) {
	// Get receiving details
	var vendorID, poID int
	var receivingDate time.Time
	err := s.db.QueryRow(ctx, `
		SELECT r.vendor_id, r.po_id, r.receiving_date
		FROM receiving r
		WHERE r.id = $1`, receivingID).Scan(&vendorID, &poID, &receivingDate)
	if err != nil {
		return 0, fmt.Errorf("receiving not found")
	}

	// Generate invoice number
	invoiceNumber := fmt.Sprintf("RCV-%d", receivingID)

	// Create invoice request
	req := &models.CreateAPInvoiceRequest{
		VendorID:      vendorID,
		InvoiceNumber: invoiceNumber,
		POID:          &poID,
		ReceivingID:   &receivingID,
		InvoiceDate:   receivingDate.Format("2006-01-02"),
	}

	// Get receiving lines
	rows := s.db.Query(ctx, `
		SELECT rl.product_id, COALESCE(p.name, 'Product'), rl.quantity_received, rl.unit_cost
		FROM receiving_lines rl
		LEFT JOIN products p ON rl.product_id = p.id
		WHERE rl.receiving_id = $1`, receivingID)
	defer rows.Close()

	for rows.Next() {
		var line models.CreateAPInvoiceLineReq
		rows.Scan(&line.ProductID, &line.Description, &line.Quantity, &line.UnitCost)
		req.Lines = append(req.Lines, line)
	}

	return s.CreateInvoice(ctx, req, createdBy)
}

// ============================================
// Payments
// ============================================

func (s *apServiceImpl) CreatePayment(ctx context.Context, req *models.CreateAPPaymentRequest, preparedBy int) (int, error) {
	paymentNumber := s.generatePaymentNumber(ctx)
	paymentDate, _ := time.Parse("2006-01-02", req.PaymentDate)

	query := `
		INSERT INTO ap_payments (
			payment_number, vendor_id, payment_date, payment_method, amount,
			currency, check_number, bank_account_id, reference_no, notes, prepared_by
		) VALUES ($1, $2, $3, $4, $5, 'USD', $6, $7, $8, $9, $10)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		paymentNumber, req.VendorID, paymentDate, req.PaymentMethod, req.Amount,
		req.CheckNumber, req.BankAccountID, req.ReferenceNo, req.Notes, preparedBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create payment: %w", err)
	}

	// Apply to invoices
	for _, app := range req.Applications {
		_, err := s.db.Exec(ctx, `
			INSERT INTO ap_payment_applications (payment_id, invoice_id, amount)
			VALUES ($1, $2, $3)`, id, app.InvoiceID, app.Amount)
		if err != nil {
			return 0, fmt.Errorf("failed to apply payment: %w", err)
		}

		// Update invoice
		s.db.Exec(ctx, `
			UPDATE ap_invoices SET
				amount_paid = amount_paid + $1,
				balance_due = balance_due - $1,
				status = CASE 
					WHEN balance_due - $1 <= 0 THEN 'PAID'::ap_invoice_status
					ELSE 'PARTIAL'::ap_invoice_status
				END,
				updated_at = NOW()
			WHERE id = $2`, app.Amount, app.InvoiceID)
	}

	return id, nil
}

func (s *apServiceImpl) GetPayment(ctx context.Context, id int) (*models.APPaymentWithDetails, error) {
	query := `
		SELECT p.id, p.payment_number, p.vendor_id, p.payment_date, p.payment_method,
			   p.amount, p.currency, p.check_number, p.bank_account_id, p.reference_no,
			   p.notes, p.prepared_by, p.approved_by, p.is_voided, p.created_at,
			   v.name as vendor_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as prepared_by_name
		FROM ap_payments p
		JOIN vendors v ON p.vendor_id = v.id
		LEFT JOIN employees e ON p.prepared_by = e.id
		WHERE p.id = $1`

	var pay models.APPaymentWithDetails
	var checkNo, refNo, notes *string

	err := s.db.QueryRow(ctx, query, id).Scan(
		&pay.Payment.ID, &pay.Payment.PaymentNumber, &pay.Payment.VendorID, &pay.Payment.PaymentDate,
		&pay.Payment.PaymentMethod, &pay.Payment.Amount, &pay.Payment.Currency,
		&checkNo, &pay.Payment.BankAccountID, &refNo, &notes, &pay.Payment.PreparedBy,
		&pay.Payment.ApprovedBy, &pay.Payment.IsVoided, &pay.Payment.CreatedAt,
		&pay.VendorName, &pay.PreparedByName,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	if checkNo != nil {
		pay.Payment.CheckNumber = *checkNo
	}
	if refNo != nil {
		pay.Payment.ReferenceNo = *refNo
	}
	if notes != nil {
		pay.Payment.Notes = *notes
	}

	// Get applications
	appRows := s.db.Query(ctx, `SELECT id, payment_id, invoice_id, amount FROM ap_payment_applications WHERE payment_id = $1`, id)
	defer appRows.Close()

	for appRows.Next() {
		var app models.APPaymentApplication
		appRows.Scan(&app.ID, &app.PaymentID, &app.InvoiceID, &app.Amount)
		pay.Applications = append(pay.Applications, app)
	}

	return &pay, nil
}

func (s *apServiceImpl) ListPayments(ctx context.Context, vendorID *int, limit int) ([]models.APPaymentWithDetails, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	whereClause := "WHERE p.is_voided = false"
	args := []interface{}{}
	argNum := 1

	if vendorID != nil {
		whereClause += fmt.Sprintf(" AND p.vendor_id = $%d", argNum)
		args = append(args, *vendorID)
		argNum++
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.payment_number, p.vendor_id, p.payment_date, p.payment_method,
			   p.amount, p.currency, p.check_number, p.bank_account_id, p.reference_no,
			   p.notes, p.prepared_by, p.approved_by, p.is_voided, p.created_at,
			   v.name as vendor_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as prepared_by_name
		FROM ap_payments p
		JOIN vendors v ON p.vendor_id = v.id
		LEFT JOIN employees e ON p.prepared_by = e.id
		%s
		ORDER BY p.payment_date DESC
		LIMIT $%d`, whereClause, argNum)

	args = append(args, limit)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var payments []models.APPaymentWithDetails
	for rows.Next() {
		var pay models.APPaymentWithDetails
		var checkNo, refNo, notes *string

		err := rows.Scan(
			&pay.Payment.ID, &pay.Payment.PaymentNumber, &pay.Payment.VendorID, &pay.Payment.PaymentDate,
			&pay.Payment.PaymentMethod, &pay.Payment.Amount, &pay.Payment.Currency,
			&checkNo, &pay.Payment.BankAccountID, &refNo, &notes, &pay.Payment.PreparedBy,
			&pay.Payment.ApprovedBy, &pay.Payment.IsVoided, &pay.Payment.CreatedAt,
			&pay.VendorName, &pay.PreparedByName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
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

func (s *apServiceImpl) VoidPayment(ctx context.Context, id int) error {
	// Get payment applications first to reverse
	appRows := s.db.Query(ctx, `SELECT invoice_id, amount FROM ap_payment_applications WHERE payment_id = $1`, id)
	defer appRows.Close()

	for appRows.Next() {
		var invoiceID int
		var amount float64
		appRows.Scan(&invoiceID, &amount)

		// Reverse invoice payment
		s.db.Exec(ctx, `
			UPDATE ap_invoices SET
				amount_paid = amount_paid - $1,
				balance_due = balance_due + $1,
				status = CASE 
					WHEN amount_paid - $1 <= 0 THEN 'APPROVED'::ap_invoice_status
					ELSE 'PARTIAL'::ap_invoice_status
				END,
				updated_at = NOW()
			WHERE id = $2`, amount, invoiceID)
	}

	result, err := s.db.Exec(ctx, `UPDATE ap_payments SET is_voided = true WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to void payment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment not found")
	}
	return nil
}

// ============================================
// Vendor Balance
// ============================================

func (s *apServiceImpl) GetVendorBalance(ctx context.Context, vendorID int) (*models.VendorBalance, error) {
	query := `
		SELECT v.id, v.name, v.payment_terms_days,
			   COALESCE((SELECT SUM(balance_due) FROM ap_invoices WHERE vendor_id = v.id AND status NOT IN ('VOID', 'PAID')), 0) as current_balance,
			   COALESCE((SELECT SUM(balance_due) FROM ap_invoices WHERE vendor_id = v.id AND due_date < CURRENT_DATE AND balance_due > 0), 0) as total_overdue,
			   COALESCE((SELECT MAX(CURRENT_DATE - due_date) FROM ap_invoices WHERE vendor_id = v.id AND due_date < CURRENT_DATE AND balance_due > 0), 0) as oldest_overdue
		FROM vendors v
		WHERE v.id = $1`

	var bal models.VendorBalance
	err := s.db.QueryRow(ctx, query, vendorID).Scan(
		&bal.VendorID, &bal.VendorName, &bal.PaymentTerms,
		&bal.CurrentBalance, &bal.TotalOverdue, &bal.OldestOverdue,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("vendor not found")
		}
		return nil, fmt.Errorf("failed to get vendor balance: %w", err)
	}

	return &bal, nil
}

// ============================================
// Aging
// ============================================

func (s *apServiceImpl) GetVendorAging(ctx context.Context, vendorID int) (*models.VendorAging, error) {
	query := `
		SELECT v.id, v.name, v.vendor_code,
			   COALESCE(SUM(CASE WHEN i.due_date >= CURRENT_DATE THEN i.balance_due ELSE 0 END), 0) as current_amt,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 1 AND 30 THEN i.balance_due ELSE 0 END), 0) as days_1_30,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 31 AND 60 THEN i.balance_due ELSE 0 END), 0) as days_31_60,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 61 AND 90 THEN i.balance_due ELSE 0 END), 0) as days_61_90,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date > 90 THEN i.balance_due ELSE 0 END), 0) as over_90,
			   COALESCE(SUM(i.balance_due), 0) as total
		FROM vendors v
		LEFT JOIN ap_invoices i ON v.id = i.vendor_id AND i.status NOT IN ('VOID', 'PAID')
		WHERE v.id = $1
		GROUP BY v.id, v.name, v.vendor_code`

	var aging models.VendorAging
	err := s.db.QueryRow(ctx, query, vendorID).Scan(
		&aging.VendorID, &aging.VendorName, &aging.VendorCode,
		&aging.Aging.Current, &aging.Aging.Days1_30, &aging.Aging.Days31_60,
		&aging.Aging.Days61_90, &aging.Aging.Over90, &aging.Aging.Total,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("vendor not found")
		}
		return nil, fmt.Errorf("failed to get vendor aging: %w", err)
	}

	return &aging, nil
}

func (s *apServiceImpl) GetAgingReport(ctx context.Context) ([]models.VendorAging, error) {
	query := `
		SELECT v.id, v.name, v.vendor_code,
			   COALESCE(SUM(CASE WHEN i.due_date >= CURRENT_DATE THEN i.balance_due ELSE 0 END), 0) as current_amt,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 1 AND 30 THEN i.balance_due ELSE 0 END), 0) as days_1_30,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 31 AND 60 THEN i.balance_due ELSE 0 END), 0) as days_31_60,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date BETWEEN 61 AND 90 THEN i.balance_due ELSE 0 END), 0) as days_61_90,
			   COALESCE(SUM(CASE WHEN CURRENT_DATE - i.due_date > 90 THEN i.balance_due ELSE 0 END), 0) as over_90,
			   COALESCE(SUM(i.balance_due), 0) as total
		FROM vendors v
		LEFT JOIN ap_invoices i ON v.id = i.vendor_id AND i.status NOT IN ('VOID', 'PAID')
		GROUP BY v.id, v.name, v.vendor_code
		HAVING COALESCE(SUM(i.balance_due), 0) > 0
		ORDER BY total DESC`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var report []models.VendorAging
	for rows.Next() {
		var aging models.VendorAging
		err := rows.Scan(
			&aging.VendorID, &aging.VendorName, &aging.VendorCode,
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
// Due & Overdue
// ============================================

func (s *apServiceImpl) GetDueInvoices(ctx context.Context, withinDays int) ([]models.APInvoiceWithDetails, error) {
	filters := &models.APInvoiceListFilters{
		Page:     1,
		PageSize: 100,
	}
	invoices, _, err := s.ListInvoices(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Filter by due within days
	dueDate := time.Now().AddDate(0, 0, withinDays)
	var result []models.APInvoiceWithDetails
	for _, inv := range invoices {
		if time.Time(inv.Invoice.DueDate).Before(dueDate) && inv.Invoice.BalanceDue > 0 {
			result = append(result, inv)
		}
	}

	return result, nil
}

func (s *apServiceImpl) GetOverdueInvoices(ctx context.Context, daysOverdue int) ([]models.APInvoiceWithDetails, error) {
	filters := &models.APInvoiceListFilters{
		Overdue:  true,
		Page:     1,
		PageSize: 100,
	}
	invoices, _, err := s.ListInvoices(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Filter by days overdue
	var result []models.APInvoiceWithDetails
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

func (s *apServiceImpl) generatePaymentNumber(ctx context.Context) string {
	var count int64
	s.db.QueryRow(ctx, `SELECT COUNT(*) FROM ap_payments WHERE DATE(created_at) = CURRENT_DATE`).Scan(&count)
	return fmt.Sprintf("PMT%s%04d", time.Now().Format("20060102"), count+1)
}
