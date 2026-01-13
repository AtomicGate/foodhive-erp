package sales_order

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

type SalesOrderService interface {
	// Sales Orders
	Create(ctx context.Context, req *models.CreateSalesOrderRequest, createdBy int) (int, error)
	GetByID(ctx context.Context, id int) (*models.SalesOrderWithDetails, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*models.SalesOrderWithDetails, error)
	Update(ctx context.Context, id int, req *models.UpdateSalesOrderRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters *models.SalesOrderListFilters) ([]models.SalesOrderWithDetails, int64, error)
	Confirm(ctx context.Context, id int) error
	Cancel(ctx context.Context, id int) error
	Ship(ctx context.Context, id int) error

	// Sales Order Lines
	AddLine(ctx context.Context, orderID int, req *models.CreateSalesOrderLineRequest) (int, error)
	UpdateLine(ctx context.Context, lineID int, req *models.CreateSalesOrderLineRequest) error
	DeleteLine(ctx context.Context, lineID int) error

	// Order Guide
	GetOrderGuide(ctx context.Context, customerID int, warehouseID int) ([]models.OrderGuideEntry, error)

	// Lost Sales
	RecordLostSale(ctx context.Context, orderID, productID int, qtyRequested, qtyAvailable float64, reason string) error
	GetLostSales(ctx context.Context, orderID *int, limit int) ([]models.LostSale, error)
}

// ============================================
// Service Implementation
// ============================================

type salesOrderServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) SalesOrderService {
	return &salesOrderServiceImpl{db: db}
}

// ============================================
// Sales Order CRUD
// ============================================

func (s *salesOrderServiceImpl) Create(ctx context.Context, req *models.CreateSalesOrderRequest, createdBy int) (int, error) {
	// Generate order number
	orderNumber := s.generateOrderNumber(ctx)

	// Parse requested ship date
	var reqShipDate *time.Time
	if req.RequestedShipDate != "" {
		t, err := time.Parse("2006-01-02", req.RequestedShipDate)
		if err == nil {
			reqShipDate = &t
		}
	}

	// Set default order type if not provided
	orderType := req.OrderType
	if orderType == "" {
		orderType = models.OrderTypeStandard
	}

	// Insert header
	query := `
		INSERT INTO sales_orders (
			order_number, customer_id, ship_to_id, order_type, requested_ship_date,
			warehouse_id, route_id, status, notes, po_number, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 'DRAFT', $8, $9, $10)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		orderNumber, req.CustomerID, req.ShipToID, orderType, reqShipDate,
		req.WarehouseID, req.RouteID, req.Notes, req.PONumber, createdBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create sales order: %w", err)
	}

	// Insert lines
	for i, line := range req.Lines {
		// Get product price if not provided
		unitPrice := line.UnitPrice
		if unitPrice == 0 {
			unitPrice = s.getProductPrice(ctx, req.CustomerID, line.ProductID)
		}

		lineTotal := line.Quantity * unitPrice * (1 - line.DiscountPercent/100)

		lineQuery := `
			INSERT INTO sales_order_lines (
				order_id, line_number, product_id, description, quantity_ordered,
				unit_of_measure, unit_price, discount_percent, line_total, lot_number
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		_, err := s.db.Exec(ctx, lineQuery,
			id, i+1, line.ProductID, line.Notes, line.Quantity,
			line.UnitOfMeasure, unitPrice, line.DiscountPercent, lineTotal, line.LotNumber,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to create order line: %w", err)
		}
	}

	// Calculate totals
	s.recalculateTotals(ctx, id)

	return id, nil
}

func (s *salesOrderServiceImpl) GetByID(ctx context.Context, id int) (*models.SalesOrderWithDetails, error) {
	return s.getSalesOrder(ctx, "so.id = $1", id)
}

func (s *salesOrderServiceImpl) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.SalesOrderWithDetails, error) {
	return s.getSalesOrder(ctx, "so.order_number = $1", orderNumber)
}

func (s *salesOrderServiceImpl) getSalesOrder(ctx context.Context, whereClause string, arg interface{}) (*models.SalesOrderWithDetails, error) {
	query := fmt.Sprintf(`
		SELECT so.id, so.order_number, so.customer_id, so.ship_to_id, so.order_type,
			   so.order_date, so.requested_ship_date, so.actual_ship_date, so.warehouse_id,
			   so.route_id, so.status, so.subtotal, so.tax_amount, so.freight_amount,
			   so.discount_amount, so.total_amount, so.notes, so.po_number, so.sales_rep_id,
			   so.created_by, so.created_at, so.updated_at,
			   c.name as customer_name, c.customer_code,
			   COALESCE(cst.name, '') as ship_to_name,
			   COALESCE(cst.address_line1 || ', ' || cst.city, '') as ship_to_address,
			   w.name as warehouse_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as sales_rep_name,
			   COALESCE(rt.name, '') as route_name
		FROM sales_orders so
		JOIN customers c ON so.customer_id = c.id
		JOIN warehouses w ON so.warehouse_id = w.id
		LEFT JOIN customer_ship_to cst ON so.ship_to_id = cst.id
		LEFT JOIN employees e ON so.sales_rep_id = e.id
		LEFT JOIN routes rt ON so.route_id = rt.id
		WHERE %s`, whereClause)

	var order models.SalesOrderWithDetails
	var reqShipDate, actShipDate *time.Time
	var notes, poNumber *string

	err := s.db.QueryRow(ctx, query, arg).Scan(
		&order.Order.ID, &order.Order.OrderNumber, &order.Order.CustomerID, &order.Order.ShipToID,
		&order.Order.OrderType, &order.Order.OrderDate, &reqShipDate, &actShipDate,
		&order.Order.WarehouseID, &order.Order.RouteID, &order.Order.Status,
		&order.Order.Subtotal, &order.Order.TaxAmount, &order.Order.FreightAmount,
		&order.Order.DiscountAmount, &order.Order.TotalAmount, &notes, &poNumber,
		&order.Order.SalesRepID, &order.Order.CreatedBy, &order.Order.CreatedAt, &order.Order.UpdatedAt,
		&order.CustomerName, &order.CustomerCode, &order.ShipToName, &order.ShipToAddress,
		&order.WarehouseName, &order.SalesRepName, &order.RouteName,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("sales order not found")
		}
		return nil, fmt.Errorf("failed to get sales order: %w", err)
	}

	if notes != nil {
		order.Order.Notes = *notes
	}
	if poNumber != nil {
		order.Order.PONumber = *poNumber
	}

	// Get lines
	linesQuery := `
		SELECT id, order_id, line_number, product_id, description, quantity_ordered,
			   quantity_shipped, unit_of_measure, unit_price, discount_percent, line_total,
			   lot_number, expiry_date, catch_weight, cost
		FROM sales_order_lines
		WHERE order_id = $1
		ORDER BY line_number`

	rows := s.db.Query(ctx, linesQuery, order.Order.ID)
	defer rows.Close()

	for rows.Next() {
		var line models.SalesOrderLine
		var desc, lotNum *string
		var expDate *time.Time

		err := rows.Scan(
			&line.ID, &line.OrderID, &line.LineNumber, &line.ProductID, &desc,
			&line.QuantityOrdered, &line.QuantityShipped, &line.UnitOfMeasure,
			&line.UnitPrice, &line.DiscountPercent, &line.LineTotal,
			&lotNum, &expDate, &line.CatchWeight, &line.Cost,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order line: %w", err)
		}
		if desc != nil {
			line.Description = *desc
		}
		if lotNum != nil {
			line.LotNumber = *lotNum
		}
		order.Lines = append(order.Lines, line)
	}

	return &order, nil
}

func (s *salesOrderServiceImpl) Update(ctx context.Context, id int, req *models.UpdateSalesOrderRequest) error {
	query := `UPDATE sales_orders SET updated_at = NOW()`
	args := []interface{}{}
	argNum := 1

	if req.ShipToID != nil {
		query += fmt.Sprintf(", ship_to_id = $%d", argNum)
		args = append(args, *req.ShipToID)
		argNum++
	}
	if req.RequestedShipDate != nil {
		t, _ := time.Parse("2006-01-02", *req.RequestedShipDate)
		query += fmt.Sprintf(", requested_ship_date = $%d", argNum)
		args = append(args, t)
		argNum++
	}
	if req.RouteID != nil {
		query += fmt.Sprintf(", route_id = $%d", argNum)
		args = append(args, *req.RouteID)
		argNum++
	}
	if req.Notes != nil {
		query += fmt.Sprintf(", notes = $%d", argNum)
		args = append(args, *req.Notes)
		argNum++
	}
	if req.PONumber != nil {
		query += fmt.Sprintf(", po_number = $%d", argNum)
		args = append(args, *req.PONumber)
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
		return fmt.Errorf("failed to update sales order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("sales order not found")
	}

	return nil
}

func (s *salesOrderServiceImpl) Delete(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx, `DELETE FROM sales_orders WHERE id = $1 AND status = 'DRAFT'`, id)
	if err != nil {
		return fmt.Errorf("failed to delete sales order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("sales order not found or cannot be deleted")
	}
	return nil
}

func (s *salesOrderServiceImpl) List(ctx context.Context, filters *models.SalesOrderListFilters) ([]models.SalesOrderWithDetails, int64, error) {
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
		whereClause += fmt.Sprintf(" AND so.customer_id = $%d", argNum)
		args = append(args, *filters.CustomerID)
		argNum++
	}
	if filters.Status != nil {
		whereClause += fmt.Sprintf(" AND so.status = $%d", argNum)
		args = append(args, *filters.Status)
		argNum++
	}
	if filters.OrderType != nil {
		whereClause += fmt.Sprintf(" AND so.order_type = $%d", argNum)
		args = append(args, *filters.OrderType)
		argNum++
	}
	if filters.WarehouseID != nil {
		whereClause += fmt.Sprintf(" AND so.warehouse_id = $%d", argNum)
		args = append(args, *filters.WarehouseID)
		argNum++
	}
	if filters.RouteID != nil {
		whereClause += fmt.Sprintf(" AND so.route_id = $%d", argNum)
		args = append(args, *filters.RouteID)
		argNum++
	}
	if filters.SalesRepID != nil {
		whereClause += fmt.Sprintf(" AND so.sales_rep_id = $%d", argNum)
		args = append(args, *filters.SalesRepID)
		argNum++
	}
	if filters.DateFrom != "" {
		whereClause += fmt.Sprintf(" AND so.order_date >= $%d", argNum)
		args = append(args, filters.DateFrom)
		argNum++
	}
	if filters.DateTo != "" {
		whereClause += fmt.Sprintf(" AND so.order_date <= $%d", argNum)
		args = append(args, filters.DateTo)
		argNum++
	}

	// Count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM sales_orders so %s`, whereClause)
	var total int64
	s.db.QueryRow(ctx, countQuery, args...).Scan(&total)

	// Query
	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT so.id, so.order_number, so.customer_id, so.ship_to_id, so.order_type,
			   so.order_date, so.requested_ship_date, so.actual_ship_date, so.warehouse_id,
			   so.route_id, so.status, so.subtotal, so.tax_amount, so.freight_amount,
			   so.discount_amount, so.total_amount, so.notes, so.po_number, so.sales_rep_id,
			   so.created_by, so.created_at, so.updated_at,
			   c.name as customer_name, c.customer_code,
			   COALESCE(cst.name, '') as ship_to_name,
			   COALESCE(cst.address_line1 || ', ' || cst.city, '') as ship_to_address,
			   w.name as warehouse_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as sales_rep_name,
			   COALESCE(rt.name, '') as route_name
		FROM sales_orders so
		JOIN customers c ON so.customer_id = c.id
		JOIN warehouses w ON so.warehouse_id = w.id
		LEFT JOIN customer_ship_to cst ON so.ship_to_id = cst.id
		LEFT JOIN employees e ON so.sales_rep_id = e.id
		LEFT JOIN routes rt ON so.route_id = rt.id
		%s
		ORDER BY so.order_date DESC, so.id DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var orders []models.SalesOrderWithDetails
	for rows.Next() {
		var order models.SalesOrderWithDetails
		var reqShipDate, actShipDate *time.Time
		var notes, poNumber *string

		err := rows.Scan(
			&order.Order.ID, &order.Order.OrderNumber, &order.Order.CustomerID, &order.Order.ShipToID,
			&order.Order.OrderType, &order.Order.OrderDate, &reqShipDate, &actShipDate,
			&order.Order.WarehouseID, &order.Order.RouteID, &order.Order.Status,
			&order.Order.Subtotal, &order.Order.TaxAmount, &order.Order.FreightAmount,
			&order.Order.DiscountAmount, &order.Order.TotalAmount, &notes, &poNumber,
			&order.Order.SalesRepID, &order.Order.CreatedBy, &order.Order.CreatedAt, &order.Order.UpdatedAt,
			&order.CustomerName, &order.CustomerCode, &order.ShipToName, &order.ShipToAddress,
			&order.WarehouseName, &order.SalesRepName, &order.RouteName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sales order: %w", err)
		}
		if notes != nil {
			order.Order.Notes = *notes
		}
		if poNumber != nil {
			order.Order.PONumber = *poNumber
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

func (s *salesOrderServiceImpl) Confirm(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx,
		`UPDATE sales_orders SET status = 'CONFIRMED', updated_at = NOW() WHERE id = $1 AND status = 'DRAFT'`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to confirm sales order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("sales order not found or not in DRAFT status")
	}
	return nil
}

func (s *salesOrderServiceImpl) Cancel(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx,
		`UPDATE sales_orders SET status = 'CANCELLED', updated_at = NOW() WHERE id = $1 AND status IN ('DRAFT', 'CONFIRMED')`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to cancel sales order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("sales order not found or cannot be cancelled")
	}
	return nil
}

func (s *salesOrderServiceImpl) Ship(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx,
		`UPDATE sales_orders SET status = 'SHIPPED', actual_ship_date = CURRENT_DATE, updated_at = NOW() 
		 WHERE id = $1 AND status = 'CONFIRMED'`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to ship sales order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("sales order not found or not in CONFIRMED status")
	}
	return nil
}

// ============================================
// Sales Order Lines
// ============================================

func (s *salesOrderServiceImpl) AddLine(ctx context.Context, orderID int, req *models.CreateSalesOrderLineRequest) (int, error) {
	// Get next line number
	var maxLine int
	s.db.QueryRow(ctx, `SELECT COALESCE(MAX(line_number), 0) FROM sales_order_lines WHERE order_id = $1`, orderID).Scan(&maxLine)

	// Get customer ID for pricing
	var customerID int
	s.db.QueryRow(ctx, `SELECT customer_id FROM sales_orders WHERE id = $1`, orderID).Scan(&customerID)

	// Get product price if not provided
	unitPrice := req.UnitPrice
	if unitPrice == 0 {
		unitPrice = s.getProductPrice(ctx, customerID, req.ProductID)
	}

	lineTotal := req.Quantity * unitPrice * (1 - req.DiscountPercent/100)

	query := `
		INSERT INTO sales_order_lines (
			order_id, line_number, product_id, description, quantity_ordered,
			unit_of_measure, unit_price, discount_percent, line_total, lot_number
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		orderID, maxLine+1, req.ProductID, req.Notes, req.Quantity,
		req.UnitOfMeasure, unitPrice, req.DiscountPercent, lineTotal, req.LotNumber,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to add order line: %w", err)
	}

	s.recalculateTotals(ctx, orderID)
	return id, nil
}

func (s *salesOrderServiceImpl) UpdateLine(ctx context.Context, lineID int, req *models.CreateSalesOrderLineRequest) error {
	lineTotal := req.Quantity * req.UnitPrice * (1 - req.DiscountPercent/100)

	query := `
		UPDATE sales_order_lines SET
			product_id = $1, description = $2, quantity_ordered = $3,
			unit_of_measure = $4, unit_price = $5, discount_percent = $6, 
			line_total = $7, lot_number = $8
		WHERE id = $9
		RETURNING order_id`

	var orderID int
	err := s.db.QueryRow(ctx, query,
		req.ProductID, req.Notes, req.Quantity,
		req.UnitOfMeasure, req.UnitPrice, req.DiscountPercent,
		lineTotal, req.LotNumber, lineID,
	).Scan(&orderID)

	if err != nil {
		return fmt.Errorf("failed to update order line: %w", err)
	}

	s.recalculateTotals(ctx, orderID)
	return nil
}

func (s *salesOrderServiceImpl) DeleteLine(ctx context.Context, lineID int) error {
	var orderID int
	err := s.db.QueryRow(ctx, `DELETE FROM sales_order_lines WHERE id = $1 RETURNING order_id`, lineID).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order line: %w", err)
	}
	s.recalculateTotals(ctx, orderID)
	return nil
}

// ============================================
// Order Guide
// ============================================

func (s *salesOrderServiceImpl) GetOrderGuide(ctx context.Context, customerID int, warehouseID int) ([]models.OrderGuideEntry, error) {
	query := `
		SELECT p.id, p.sku, p.name,
			   COALESCE(cog.default_quantity, 0) as default_qty,
			   COALESCE(cog.last_ordered_qty, 0) as last_ordered_qty,
			   COALESCE(cog.avg_weekly_qty, 0) as avg_weekly_qty,
			   COALESCE(i.quantity_on_hand, 0) as on_hand,
			   COALESCE(i.quantity_allocated, 0) as allocated,
			   COALESCE(i.quantity_available, 0) as available,
			   COALESCE(cp.price, p.base_price, 0) as unit_price,
			   COALESCE(pu.abbreviation, 'EA') as unit_of_measure,
			   COALESCE(cog.is_push_item, false) as is_push_item
		FROM customer_order_guides cog
		JOIN products p ON cog.product_id = p.id
		LEFT JOIN product_units pu ON p.default_unit_id = pu.id
		LEFT JOIN customer_pricing cp ON cp.customer_id = $1 AND cp.product_id = p.id
			AND cp.effective_date <= CURRENT_DATE AND (cp.expiry_date IS NULL OR cp.expiry_date >= CURRENT_DATE)
		LEFT JOIN inventory i ON i.product_id = p.id AND i.warehouse_id = $2
		WHERE cog.customer_id = $1 AND p.is_active = true
		ORDER BY p.name`

	rows := s.db.Query(ctx, query, customerID, warehouseID)
	defer rows.Close()

	var entries []models.OrderGuideEntry
	for rows.Next() {
		var entry models.OrderGuideEntry
		err := rows.Scan(
			&entry.ProductID, &entry.ProductSKU, &entry.ProductName,
			&entry.DefaultQuantity, &entry.LastOrderedQty, &entry.AvgWeeklyQty,
			&entry.OnHand, &entry.Allocated, &entry.Available,
			&entry.UnitPrice, &entry.UnitOfMeasure, &entry.IsPushItem,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order guide entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// ============================================
// Lost Sales
// ============================================

func (s *salesOrderServiceImpl) RecordLostSale(ctx context.Context, orderID, productID int, qtyRequested, qtyAvailable float64, reason string) error {
	query := `
		INSERT INTO lost_sales (order_id, product_id, quantity_requested, quantity_available, reason)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(ctx, query, orderID, productID, qtyRequested, qtyAvailable, reason)
	if err != nil {
		return fmt.Errorf("failed to record lost sale: %w", err)
	}
	return nil
}

func (s *salesOrderServiceImpl) GetLostSales(ctx context.Context, orderID *int, limit int) ([]models.LostSale, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if orderID != nil {
		whereClause += fmt.Sprintf(" AND ls.order_id = $%d", argNum)
		args = append(args, *orderID)
		argNum++
	}

	query := fmt.Sprintf(`
		SELECT ls.id, ls.order_id, ls.product_id, p.name, ls.quantity_requested,
			   ls.quantity_available, ls.reason, ls.created_at
		FROM lost_sales ls
		JOIN products p ON ls.product_id = p.id
		%s
		ORDER BY ls.created_at DESC
		LIMIT $%d`, whereClause, argNum)

	args = append(args, limit)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var lostSales []models.LostSale
	for rows.Next() {
		var ls models.LostSale
		var reason *string
		err := rows.Scan(
			&ls.ID, &ls.OrderID, &ls.ProductID, &ls.ProductName,
			&ls.QuantityRequested, &ls.QuantityAvailable, &reason, &ls.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lost sale: %w", err)
		}
		if reason != nil {
			ls.Reason = *reason
		}
		lostSales = append(lostSales, ls)
	}

	return lostSales, nil
}

// ============================================
// Helper Functions
// ============================================

func (s *salesOrderServiceImpl) generateOrderNumber(ctx context.Context) string {
	var count int64
	s.db.QueryRow(ctx, `SELECT COUNT(*) FROM sales_orders WHERE order_date = CURRENT_DATE`).Scan(&count)
	return fmt.Sprintf("SO%s%04d", time.Now().Format("20060102"), count+1)
}

func (s *salesOrderServiceImpl) getProductPrice(ctx context.Context, customerID, productID int) float64 {
	var price float64

	// Try customer-specific price first
	s.db.QueryRow(ctx, `
		SELECT price FROM customer_pricing 
		WHERE customer_id = $1 AND product_id = $2
		AND effective_date <= CURRENT_DATE 
		AND (expiry_date IS NULL OR expiry_date >= CURRENT_DATE)
		ORDER BY effective_date DESC LIMIT 1`, customerID, productID).Scan(&price)

	if price > 0 {
		return price
	}

	// Fall back to product base price
	s.db.QueryRow(ctx, `SELECT COALESCE(base_price, 0) FROM products WHERE id = $1`, productID).Scan(&price)
	return price
}

func (s *salesOrderServiceImpl) recalculateTotals(ctx context.Context, orderID int) {
	s.db.Exec(ctx, `
		UPDATE sales_orders SET
			subtotal = (SELECT COALESCE(SUM(line_total), 0) FROM sales_order_lines WHERE order_id = $1),
			total_amount = subtotal + tax_amount + freight_amount - discount_amount,
			updated_at = NOW()
		WHERE id = $1`, orderID)
}
