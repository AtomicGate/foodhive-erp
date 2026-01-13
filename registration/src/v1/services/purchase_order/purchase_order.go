package purchase_order

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

type PurchaseOrderService interface {
	// Purchase Orders
	Create(ctx context.Context, req *models.CreatePurchaseOrderRequest, createdBy int) (int, error)
	GetByID(ctx context.Context, id int) (*models.PurchaseOrderWithDetails, error)
	GetByPONumber(ctx context.Context, poNumber string) (*models.PurchaseOrderWithDetails, error)
	Update(ctx context.Context, id int, req *models.UpdatePurchaseOrderRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters *models.PurchaseOrderListFilters) ([]models.PurchaseOrderWithDetails, int64, error)
	Submit(ctx context.Context, id int) error
	Cancel(ctx context.Context, id int) error

	// Purchase Order Lines
	AddLine(ctx context.Context, poID int, req *models.CreatePurchaseOrderLineRequest) (int, error)
	UpdateLine(ctx context.Context, lineID int, req *models.CreatePurchaseOrderLineRequest) error
	DeleteLine(ctx context.Context, lineID int) error

	// Receiving
	CreateReceiving(ctx context.Context, req *models.CreateReceivingRequest, receivedBy int) (int, error)
	GetReceiving(ctx context.Context, id int) (*models.ReceivingWithDetails, error)
	ListReceivings(ctx context.Context, poID *int, warehouseID *int, limit int) ([]models.ReceivingWithDetails, error)
}

// ============================================
// Service Implementation
// ============================================

type purchaseOrderServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) PurchaseOrderService {
	return &purchaseOrderServiceImpl{db: db}
}

// ============================================
// Purchase Order CRUD
// ============================================

func (s *purchaseOrderServiceImpl) Create(ctx context.Context, req *models.CreatePurchaseOrderRequest, createdBy int) (int, error) {
	// Generate PO number
	poNumber := s.generatePONumber(ctx)

	// Calculate totals
	var subtotal float64
	for _, line := range req.Lines {
		subtotal += line.Quantity * line.UnitCost
	}

	// Parse expected date
	var expectedDate *time.Time
	if req.ExpectedDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpectedDate)
		if err == nil {
			expectedDate = &t
		}
	}

	// Insert header
	query := `
		INSERT INTO purchase_orders (
			po_number, vendor_id, warehouse_id, expected_date, status,
			subtotal, total_amount, notes, buyer_id, created_by
		) VALUES ($1, $2, $3, $4, 'DRAFT', $5, $5, $6, $7, $8)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		poNumber, req.VendorID, req.WarehouseID, expectedDate,
		subtotal, req.Notes, req.BuyerID, createdBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create purchase order: %w", err)
	}

	// Insert lines
	for i, line := range req.Lines {
		lineTotal := line.Quantity * line.UnitCost
		var lineExpDate *time.Time
		if line.ExpectedDate != "" {
			t, _ := time.Parse("2006-01-02", line.ExpectedDate)
			lineExpDate = &t
		}

		lineQuery := `
			INSERT INTO purchase_order_lines (
				po_id, line_number, product_id, description, quantity_ordered,
				unit_of_measure, unit_cost, line_total, expected_date
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

		_, err := s.db.Exec(ctx, lineQuery,
			id, i+1, line.ProductID, line.Description, line.Quantity,
			line.UnitOfMeasure, line.UnitCost, lineTotal, lineExpDate,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to create PO line: %w", err)
		}
	}

	return id, nil
}

func (s *purchaseOrderServiceImpl) GetByID(ctx context.Context, id int) (*models.PurchaseOrderWithDetails, error) {
	return s.getPurchaseOrder(ctx, "po.id = $1", id)
}

func (s *purchaseOrderServiceImpl) GetByPONumber(ctx context.Context, poNumber string) (*models.PurchaseOrderWithDetails, error) {
	return s.getPurchaseOrder(ctx, "po.po_number = $1", poNumber)
}

func (s *purchaseOrderServiceImpl) getPurchaseOrder(ctx context.Context, whereClause string, arg interface{}) (*models.PurchaseOrderWithDetails, error) {
	query := fmt.Sprintf(`
		SELECT po.id, po.po_number, po.vendor_id, po.warehouse_id, po.order_date,
			   po.expected_date, po.received_date, po.status, po.subtotal, po.tax_amount,
			   po.freight_amount, po.total_amount, po.notes, po.buyer_id, po.created_by,
			   po.created_at, po.updated_at,
			   v.name as vendor_name, v.code as vendor_code,
			   w.name as warehouse_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as buyer_name
		FROM purchase_orders po
		JOIN vendors v ON po.vendor_id = v.id
		JOIN warehouses w ON po.warehouse_id = w.id
		LEFT JOIN employees e ON po.buyer_id = e.id
		WHERE %s`, whereClause)

	var po models.PurchaseOrderWithDetails
	var expectedDate, receivedDate *time.Time
	var notes *string

	err := s.db.QueryRow(ctx, query, arg).Scan(
		&po.Order.ID, &po.Order.PONumber, &po.Order.VendorID, &po.Order.WarehouseID, &po.Order.OrderDate,
		&expectedDate, &receivedDate, &po.Order.Status, &po.Order.Subtotal, &po.Order.TaxAmount,
		&po.Order.FreightAmount, &po.Order.TotalAmount, &notes, &po.Order.BuyerID, &po.Order.CreatedBy,
		&po.Order.CreatedAt, &po.Order.UpdatedAt,
		&po.VendorName, &po.VendorCode, &po.WarehouseName, &po.BuyerName,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("purchase order not found")
		}
		return nil, fmt.Errorf("failed to get purchase order: %w", err)
	}

	if notes != nil {
		po.Order.Notes = *notes
	}

	// Get lines
	linesQuery := `
		SELECT id, po_id, line_number, product_id, description, quantity_ordered,
			   quantity_received, unit_of_measure, unit_cost, line_total, expected_date
		FROM purchase_order_lines
		WHERE po_id = $1
		ORDER BY line_number`

	rows := s.db.Query(ctx, linesQuery, po.Order.ID)
	defer rows.Close()

	for rows.Next() {
		var line models.PurchaseOrderLine
		var desc *string
		var expDate *time.Time

		err := rows.Scan(
			&line.ID, &line.POID, &line.LineNumber, &line.ProductID, &desc, &line.QuantityOrdered,
			&line.QuantityReceived, &line.UnitOfMeasure, &line.UnitCost, &line.LineTotal, &expDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan PO line: %w", err)
		}
		if desc != nil {
			line.Description = *desc
		}
		po.Lines = append(po.Lines, line)
	}

	return &po, nil
}

func (s *purchaseOrderServiceImpl) Update(ctx context.Context, id int, req *models.UpdatePurchaseOrderRequest) error {
	query := `UPDATE purchase_orders SET updated_at = NOW()`
	args := []interface{}{}
	argNum := 1

	if req.ExpectedDate != nil {
		t, _ := time.Parse("2006-01-02", *req.ExpectedDate)
		query += fmt.Sprintf(", expected_date = $%d", argNum)
		args = append(args, t)
		argNum++
	}
	if req.Notes != nil {
		query += fmt.Sprintf(", notes = $%d", argNum)
		args = append(args, *req.Notes)
		argNum++
	}
	if req.BuyerID != nil {
		query += fmt.Sprintf(", buyer_id = $%d", argNum)
		args = append(args, *req.BuyerID)
		argNum++
	}
	if req.Status != nil {
		query += fmt.Sprintf(", status = $%d", argNum)
		args = append(args, *req.Status)
		argNum++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argNum)
	args = append(args, id)

	result, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update purchase order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("purchase order not found")
	}

	return nil
}

func (s *purchaseOrderServiceImpl) Delete(ctx context.Context, id int) error {
	// Only allow deletion of DRAFT POs
	result, err := s.db.Exec(ctx, `DELETE FROM purchase_orders WHERE id = $1 AND status = 'DRAFT'`, id)
	if err != nil {
		return fmt.Errorf("failed to delete purchase order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("purchase order not found or cannot be deleted")
	}
	return nil
}

func (s *purchaseOrderServiceImpl) List(ctx context.Context, filters *models.PurchaseOrderListFilters) ([]models.PurchaseOrderWithDetails, int64, error) {
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
		whereClause += fmt.Sprintf(" AND po.vendor_id = $%d", argNum)
		args = append(args, *filters.VendorID)
		argNum++
	}
	if filters.WarehouseID != nil {
		whereClause += fmt.Sprintf(" AND po.warehouse_id = $%d", argNum)
		args = append(args, *filters.WarehouseID)
		argNum++
	}
	if filters.Status != nil {
		whereClause += fmt.Sprintf(" AND po.status = $%d", argNum)
		args = append(args, *filters.Status)
		argNum++
	}
	if filters.BuyerID != nil {
		whereClause += fmt.Sprintf(" AND po.buyer_id = $%d", argNum)
		args = append(args, *filters.BuyerID)
		argNum++
	}
	if filters.DateFrom != "" {
		whereClause += fmt.Sprintf(" AND po.order_date >= $%d", argNum)
		args = append(args, filters.DateFrom)
		argNum++
	}
	if filters.DateTo != "" {
		whereClause += fmt.Sprintf(" AND po.order_date <= $%d", argNum)
		args = append(args, filters.DateTo)
		argNum++
	}

	// Count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM purchase_orders po %s`, whereClause)
	var total int64
	s.db.QueryRow(ctx, countQuery, args...).Scan(&total)

	// Query
	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT po.id, po.po_number, po.vendor_id, po.warehouse_id, po.order_date,
			   po.expected_date, po.received_date, po.status, po.subtotal, po.tax_amount,
			   po.freight_amount, po.total_amount, po.notes, po.buyer_id, po.created_by,
			   po.created_at, po.updated_at,
			   v.name as vendor_name, v.code as vendor_code,
			   w.name as warehouse_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as buyer_name
		FROM purchase_orders po
		JOIN vendors v ON po.vendor_id = v.id
		JOIN warehouses w ON po.warehouse_id = w.id
		LEFT JOIN employees e ON po.buyer_id = e.id
		%s
		ORDER BY po.order_date DESC, po.id DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var orders []models.PurchaseOrderWithDetails
	for rows.Next() {
		var po models.PurchaseOrderWithDetails
		var expectedDate, receivedDate *time.Time
		var notes *string

		err := rows.Scan(
			&po.Order.ID, &po.Order.PONumber, &po.Order.VendorID, &po.Order.WarehouseID, &po.Order.OrderDate,
			&expectedDate, &receivedDate, &po.Order.Status, &po.Order.Subtotal, &po.Order.TaxAmount,
			&po.Order.FreightAmount, &po.Order.TotalAmount, &notes, &po.Order.BuyerID, &po.Order.CreatedBy,
			&po.Order.CreatedAt, &po.Order.UpdatedAt,
			&po.VendorName, &po.VendorCode, &po.WarehouseName, &po.BuyerName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan purchase order: %w", err)
		}
		if notes != nil {
			po.Order.Notes = *notes
		}
		orders = append(orders, po)
	}

	return orders, total, nil
}

func (s *purchaseOrderServiceImpl) Submit(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx,
		`UPDATE purchase_orders SET status = 'SUBMITTED', updated_at = NOW() WHERE id = $1 AND status = 'DRAFT'`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to submit purchase order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("purchase order not found or not in DRAFT status")
	}
	return nil
}

func (s *purchaseOrderServiceImpl) Cancel(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx,
		`UPDATE purchase_orders SET status = 'CANCELLED', updated_at = NOW() WHERE id = $1 AND status IN ('DRAFT', 'SUBMITTED')`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to cancel purchase order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("purchase order not found or cannot be cancelled")
	}
	return nil
}

// ============================================
// Purchase Order Lines
// ============================================

func (s *purchaseOrderServiceImpl) AddLine(ctx context.Context, poID int, req *models.CreatePurchaseOrderLineRequest) (int, error) {
	// Get next line number
	var maxLine int
	s.db.QueryRow(ctx, `SELECT COALESCE(MAX(line_number), 0) FROM purchase_order_lines WHERE po_id = $1`, poID).Scan(&maxLine)

	lineTotal := req.Quantity * req.UnitCost
	var expDate *time.Time
	if req.ExpectedDate != "" {
		t, _ := time.Parse("2006-01-02", req.ExpectedDate)
		expDate = &t
	}

	query := `
		INSERT INTO purchase_order_lines (
			po_id, line_number, product_id, description, quantity_ordered,
			unit_of_measure, unit_cost, line_total, expected_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		poID, maxLine+1, req.ProductID, req.Description, req.Quantity,
		req.UnitOfMeasure, req.UnitCost, lineTotal, expDate,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to add PO line: %w", err)
	}

	// Update PO totals
	s.recalculateTotals(ctx, poID)

	return id, nil
}

func (s *purchaseOrderServiceImpl) UpdateLine(ctx context.Context, lineID int, req *models.CreatePurchaseOrderLineRequest) error {
	lineTotal := req.Quantity * req.UnitCost
	var expDate *time.Time
	if req.ExpectedDate != "" {
		t, _ := time.Parse("2006-01-02", req.ExpectedDate)
		expDate = &t
	}

	query := `
		UPDATE purchase_order_lines SET
			product_id = $1, description = $2, quantity_ordered = $3,
			unit_of_measure = $4, unit_cost = $5, line_total = $6, expected_date = $7
		WHERE id = $8
		RETURNING po_id`

	var poID int
	err := s.db.QueryRow(ctx, query,
		req.ProductID, req.Description, req.Quantity,
		req.UnitOfMeasure, req.UnitCost, lineTotal, expDate, lineID,
	).Scan(&poID)

	if err != nil {
		return fmt.Errorf("failed to update PO line: %w", err)
	}

	s.recalculateTotals(ctx, poID)
	return nil
}

func (s *purchaseOrderServiceImpl) DeleteLine(ctx context.Context, lineID int) error {
	var poID int
	err := s.db.QueryRow(ctx, `DELETE FROM purchase_order_lines WHERE id = $1 RETURNING po_id`, lineID).Scan(&poID)
	if err != nil {
		return fmt.Errorf("failed to delete PO line: %w", err)
	}
	s.recalculateTotals(ctx, poID)
	return nil
}

// ============================================
// Receiving
// ============================================

func (s *purchaseOrderServiceImpl) CreateReceiving(ctx context.Context, req *models.CreateReceivingRequest, receivedBy int) (int, error) {
	// Generate receiving number
	recvNumber := s.generateReceivingNumber(ctx)

	// Insert receiving header
	query := `
		INSERT INTO receiving (
			receiving_number, po_id, warehouse_id, vendor_id, notes, received_by
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		recvNumber, req.POID, req.WarehouseID, req.VendorID, req.Notes, receivedBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create receiving: %w", err)
	}

	// Insert lines and update inventory
	for _, line := range req.Lines {
		var prodDate, expDate *time.Time
		if line.ProductionDate != "" {
			t, _ := time.Parse("2006-01-02", line.ProductionDate)
			prodDate = &t
		}
		if line.ExpiryDate != "" {
			t, _ := time.Parse("2006-01-02", line.ExpiryDate)
			expDate = &t
		}

		lineQuery := `
			INSERT INTO receiving_lines (
				receiving_id, po_line_id, product_id, quantity_received, unit_of_measure,
				lot_number, production_date, expiry_date, location_code, unit_cost, notes
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

		_, err := s.db.Exec(ctx, lineQuery,
			id, line.POLineID, line.ProductID, line.Quantity, line.UnitOfMeasure,
			line.LotNumber, prodDate, expDate, line.LocationCode, line.UnitCost, line.Notes,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to create receiving line: %w", err)
		}

		// Update PO line received quantity if linked
		if line.POLineID != nil {
			s.db.Exec(ctx,
				`UPDATE purchase_order_lines SET quantity_received = quantity_received + $1 WHERE id = $2`,
				line.Quantity, *line.POLineID,
			)
		}

		// Update inventory (upsert)
		invQuery := `
			INSERT INTO inventory (
				product_id, warehouse_id, location_code, lot_number, production_date, expiry_date,
				quantity_on_hand, last_cost, average_cost, last_movement_date
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8, NOW())
			ON CONFLICT (product_id, warehouse_id, location_code, lot_number)
			DO UPDATE SET
				quantity_on_hand = inventory.quantity_on_hand + $7,
				last_cost = $8,
				average_cost = (inventory.average_cost * inventory.quantity_on_hand + $8 * $7) / NULLIF(inventory.quantity_on_hand + $7, 0),
				last_movement_date = NOW(),
				updated_at = NOW()`

		s.db.Exec(ctx, invQuery,
			line.ProductID, req.WarehouseID, line.LocationCode, line.LotNumber,
			prodDate, expDate, line.Quantity, line.UnitCost,
		)
	}

	// Update PO status
	if req.POID != nil {
		s.updatePOStatus(ctx, *req.POID)
	}

	return id, nil
}

func (s *purchaseOrderServiceImpl) GetReceiving(ctx context.Context, id int) (*models.ReceivingWithDetails, error) {
	query := `
		SELECT r.id, r.receiving_number, r.po_id, r.warehouse_id, r.vendor_id,
			   r.received_date, r.notes, r.received_by, r.created_at,
			   v.name as vendor_name, w.name as warehouse_name,
			   COALESCE(po.po_number, '') as po_number,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as receiver_name
		FROM receiving r
		JOIN vendors v ON r.vendor_id = v.id
		JOIN warehouses w ON r.warehouse_id = w.id
		LEFT JOIN purchase_orders po ON r.po_id = po.id
		LEFT JOIN employees e ON r.received_by = e.id
		WHERE r.id = $1`

	var recv models.ReceivingWithDetails
	var notes *string

	err := s.db.QueryRow(ctx, query, id).Scan(
		&recv.Receiving.ID, &recv.Receiving.ReceivingNumber, &recv.Receiving.POID, &recv.Receiving.WarehouseID,
		&recv.Receiving.VendorID, &recv.Receiving.ReceivedDate, &notes, &recv.Receiving.ReceivedBy,
		&recv.Receiving.CreatedAt, &recv.VendorName, &recv.WarehouseName, &recv.PONumber, &recv.ReceiverName,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("receiving not found")
		}
		return nil, fmt.Errorf("failed to get receiving: %w", err)
	}
	if notes != nil {
		recv.Receiving.Notes = *notes
	}

	// Get lines
	linesQuery := `
		SELECT id, receiving_id, po_line_id, product_id, quantity_received, unit_of_measure,
			   lot_number, production_date, expiry_date, location_code, unit_cost, is_short_shipment, notes
		FROM receiving_lines
		WHERE receiving_id = $1`

	rows := s.db.Query(ctx, linesQuery, id)
	defer rows.Close()

	for rows.Next() {
		var line models.ReceivingLine
		var lotNum, locCode, lineNotes *string
		var prodDate, expDate *time.Time

		err := rows.Scan(
			&line.ID, &line.ReceivingID, &line.POLineID, &line.ProductID, &line.QuantityReceived,
			&line.UnitOfMeasure, &lotNum, &prodDate, &expDate, &locCode, &line.UnitCost,
			&line.IsShortShipment, &lineNotes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan receiving line: %w", err)
		}
		if lotNum != nil {
			line.LotNumber = *lotNum
		}
		if locCode != nil {
			line.LocationCode = *locCode
		}
		if lineNotes != nil {
			line.Notes = *lineNotes
		}
		recv.Lines = append(recv.Lines, line)
	}

	return &recv, nil
}

func (s *purchaseOrderServiceImpl) ListReceivings(ctx context.Context, poID *int, warehouseID *int, limit int) ([]models.ReceivingWithDetails, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if poID != nil {
		whereClause += fmt.Sprintf(" AND r.po_id = $%d", argNum)
		args = append(args, *poID)
		argNum++
	}
	if warehouseID != nil {
		whereClause += fmt.Sprintf(" AND r.warehouse_id = $%d", argNum)
		args = append(args, *warehouseID)
		argNum++
	}

	query := fmt.Sprintf(`
		SELECT r.id, r.receiving_number, r.po_id, r.warehouse_id, r.vendor_id,
			   r.received_date, r.notes, r.received_by, r.created_at,
			   v.name as vendor_name, w.name as warehouse_name,
			   COALESCE(po.po_number, '') as po_number,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as receiver_name
		FROM receiving r
		JOIN vendors v ON r.vendor_id = v.id
		JOIN warehouses w ON r.warehouse_id = w.id
		LEFT JOIN purchase_orders po ON r.po_id = po.id
		LEFT JOIN employees e ON r.received_by = e.id
		%s
		ORDER BY r.received_date DESC, r.id DESC
		LIMIT $%d`, whereClause, argNum)

	args = append(args, limit)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var receivings []models.ReceivingWithDetails
	for rows.Next() {
		var recv models.ReceivingWithDetails
		var notes *string

		err := rows.Scan(
			&recv.Receiving.ID, &recv.Receiving.ReceivingNumber, &recv.Receiving.POID, &recv.Receiving.WarehouseID,
			&recv.Receiving.VendorID, &recv.Receiving.ReceivedDate, &notes, &recv.Receiving.ReceivedBy,
			&recv.Receiving.CreatedAt, &recv.VendorName, &recv.WarehouseName, &recv.PONumber, &recv.ReceiverName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan receiving: %w", err)
		}
		if notes != nil {
			recv.Receiving.Notes = *notes
		}
		receivings = append(receivings, recv)
	}

	return receivings, nil
}

// ============================================
// Helper Functions
// ============================================

func (s *purchaseOrderServiceImpl) generatePONumber(ctx context.Context) string {
	var count int64
	s.db.QueryRow(ctx, `SELECT COUNT(*) FROM purchase_orders WHERE order_date = CURRENT_DATE`).Scan(&count)
	return fmt.Sprintf("PO%s%04d", time.Now().Format("20060102"), count+1)
}

func (s *purchaseOrderServiceImpl) generateReceivingNumber(ctx context.Context) string {
	var count int64
	s.db.QueryRow(ctx, `SELECT COUNT(*) FROM receiving WHERE received_date = CURRENT_DATE`).Scan(&count)
	return fmt.Sprintf("RCV%s%04d", time.Now().Format("20060102"), count+1)
}

func (s *purchaseOrderServiceImpl) recalculateTotals(ctx context.Context, poID int) {
	s.db.Exec(ctx, `
		UPDATE purchase_orders SET
			subtotal = (SELECT COALESCE(SUM(line_total), 0) FROM purchase_order_lines WHERE po_id = $1),
			total_amount = subtotal + tax_amount + freight_amount,
			updated_at = NOW()
		WHERE id = $1`, poID)
}

func (s *purchaseOrderServiceImpl) updatePOStatus(ctx context.Context, poID int) {
	// Check if fully or partially received
	var ordered, received float64
	s.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(quantity_ordered), 0), COALESCE(SUM(quantity_received), 0)
		FROM purchase_order_lines WHERE po_id = $1`, poID).Scan(&ordered, &received)

	var status models.POStatus
	if received >= ordered {
		status = models.POStatusReceived
	} else if received > 0 {
		status = models.POStatusPartial
	} else {
		return
	}

	s.db.Exec(ctx, `UPDATE purchase_orders SET status = $1, updated_at = NOW() WHERE id = $2`, status, poID)
}
