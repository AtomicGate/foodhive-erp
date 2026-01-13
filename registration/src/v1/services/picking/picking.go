package picking

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

type PickingService interface {
	// Routes
	CreateRoute(ctx context.Context, req *models.CreateRouteRequest) (int, error)
	GetRoute(ctx context.Context, id int) (*models.RouteWithDetails, error)
	UpdateRoute(ctx context.Context, id int, req *models.UpdateRouteRequest) error
	DeleteRoute(ctx context.Context, id int) error
	ListRoutes(ctx context.Context, warehouseID *int, activeOnly bool) ([]models.RouteWithDetails, error)

	// Route Stops
	AddRouteStop(ctx context.Context, routeID int, req *models.AddRouteStopRequest) (int, error)
	UpdateRouteStop(ctx context.Context, stopID int, req *models.AddRouteStopRequest) error
	DeleteRouteStop(ctx context.Context, stopID int) error
	ReorderStops(ctx context.Context, routeID int, stopOrders []int) error // Stop IDs in new order

	// Pick Lists
	CreatePickList(ctx context.Context, req *models.CreatePickListRequest, createdBy int) (int, error)
	GeneratePickListForRoute(ctx context.Context, req *models.GeneratePickListRequest, createdBy int) (int, error)
	GetPickList(ctx context.Context, id int) (*models.PickListWithDetails, error)
	ListPickLists(ctx context.Context, filters *models.PickListFilters) ([]models.PickListWithDetails, int64, error)
	StartPicking(ctx context.Context, pickListID int, pickerID int) error
	CompletePicking(ctx context.Context, pickListID int) error
	CancelPickList(ctx context.Context, pickListID int) error

	// Pick Lines
	GetPickLines(ctx context.Context, pickListID int) ([]models.PickListLineWithProduct, error)
	ConfirmPickLine(ctx context.Context, lineID int, req *models.ConfirmPickLineRequest, pickerID int) error

	// Master Pick Report (consolidated picking across orders)
	GetMasterPickReport(ctx context.Context, warehouseID int, pickDate string, routeID *int) ([]models.MasterPickItem, error)

	// Suggested Picking (FEFO - First Expiry First Out)
	GetSuggestedPicking(ctx context.Context, productID, warehouseID int, quantity float64) ([]SuggestedPickLocation, error)
}

type SuggestedPickLocation struct {
	LocationCode string     `json:"location_code"`
	LotNumber    string     `json:"lot_number"`
	ExpiryDate   *time.Time `json:"expiry_date"`
	Available    float64    `json:"available"`
	Suggested    float64    `json:"suggested"`
}

// ============================================
// Service Implementation
// ============================================

type pickingServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) PickingService {
	return &pickingServiceImpl{db: db}
}

// ============================================
// Routes
// ============================================

func (s *pickingServiceImpl) CreateRoute(ctx context.Context, req *models.CreateRouteRequest) (int, error) {
	query := `
		INSERT INTO routes (route_code, name, description, warehouse_id, driver_id, departure_time)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	var depTime *string
	if req.DepartureTime != "" {
		depTime = &req.DepartureTime
	}

	var id int
	err := s.db.QueryRow(ctx, query,
		req.RouteCode, req.Name, req.Description, req.WarehouseID, req.DriverID, depTime,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create route: %w", err)
	}
	return id, nil
}

func (s *pickingServiceImpl) GetRoute(ctx context.Context, id int) (*models.RouteWithDetails, error) {
	query := `
		SELECT r.id, r.route_code, r.name, r.description, r.warehouse_id, r.driver_id,
			   r.vehicle_id, r.departure_time, r.is_active, r.created_at,
			   COALESCE(w.name, '') as warehouse_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as driver_name,
			   (SELECT COUNT(*) FROM route_stops WHERE route_id = r.id) as stop_count
		FROM routes r
		LEFT JOIN warehouses w ON r.warehouse_id = w.id
		LEFT JOIN employees e ON r.driver_id = e.id
		WHERE r.id = $1`

	var route models.RouteWithDetails
	var desc, depTime *string

	err := s.db.QueryRow(ctx, query, id).Scan(
		&route.Route.ID, &route.Route.RouteCode, &route.Route.Name, &desc,
		&route.Route.WarehouseID, &route.Route.DriverID, &route.Route.VehicleID,
		&depTime, &route.Route.IsActive, &route.Route.CreatedAt,
		&route.WarehouseName, &route.DriverName, &route.StopCount,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("route not found")
		}
		return nil, fmt.Errorf("failed to get route: %w", err)
	}

	if desc != nil {
		route.Route.Description = *desc
	}
	if depTime != nil {
		route.Route.DepartureTime = *depTime
	}

	// Get stops
	stopsQuery := `
		SELECT rs.id, rs.route_id, rs.customer_id, c.name as customer_name,
			   rs.ship_to_id, COALESCE(cst.name, '') as ship_to_name,
			   COALESCE(cst.address_line1 || ', ' || cst.city, '') as ship_to_address,
			   rs.stop_sequence, rs.estimated_arrival, rs.notes
		FROM route_stops rs
		JOIN customers c ON rs.customer_id = c.id
		LEFT JOIN customer_ship_to cst ON rs.ship_to_id = cst.id
		WHERE rs.route_id = $1
		ORDER BY rs.stop_sequence`

	rows := s.db.Query(ctx, stopsQuery, id)
	defer rows.Close()

	for rows.Next() {
		var stop models.RouteStop
		var estArr, notes *string

		err := rows.Scan(
			&stop.ID, &stop.RouteID, &stop.CustomerID, &stop.CustomerName,
			&stop.ShipToID, &stop.ShipToName, &stop.ShipToAddress,
			&stop.StopSequence, &estArr, &notes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan route stop: %w", err)
		}
		if estArr != nil {
			stop.EstimatedArrival = *estArr
		}
		if notes != nil {
			stop.Notes = *notes
		}
		route.Stops = append(route.Stops, stop)
	}

	return &route, nil
}

func (s *pickingServiceImpl) UpdateRoute(ctx context.Context, id int, req *models.UpdateRouteRequest) error {
	query := `UPDATE routes SET id = id`
	args := []interface{}{}
	argNum := 1

	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argNum)
		args = append(args, *req.Name)
		argNum++
	}
	if req.Description != nil {
		query += fmt.Sprintf(", description = $%d", argNum)
		args = append(args, *req.Description)
		argNum++
	}
	if req.WarehouseID != nil {
		query += fmt.Sprintf(", warehouse_id = $%d", argNum)
		args = append(args, *req.WarehouseID)
		argNum++
	}
	if req.DriverID != nil {
		query += fmt.Sprintf(", driver_id = $%d", argNum)
		args = append(args, *req.DriverID)
		argNum++
	}
	if req.DepartureTime != nil {
		query += fmt.Sprintf(", departure_time = $%d", argNum)
		args = append(args, *req.DepartureTime)
		argNum++
	}
	if req.IsActive != nil {
		query += fmt.Sprintf(", is_active = $%d", argNum)
		args = append(args, *req.IsActive)
		argNum++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argNum)
	args = append(args, id)

	result, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("route not found")
	}
	return nil
}

func (s *pickingServiceImpl) DeleteRoute(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx, `DELETE FROM routes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("route not found")
	}
	return nil
}

func (s *pickingServiceImpl) ListRoutes(ctx context.Context, warehouseID *int, activeOnly bool) ([]models.RouteWithDetails, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if warehouseID != nil {
		whereClause += fmt.Sprintf(" AND r.warehouse_id = $%d", argNum)
		args = append(args, *warehouseID)
		argNum++
	}
	if activeOnly {
		whereClause += " AND r.is_active = true"
	}

	query := fmt.Sprintf(`
		SELECT r.id, r.route_code, r.name, r.description, r.warehouse_id, r.driver_id,
			   r.vehicle_id, r.departure_time, r.is_active, r.created_at,
			   COALESCE(w.name, '') as warehouse_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as driver_name,
			   (SELECT COUNT(*) FROM route_stops WHERE route_id = r.id) as stop_count
		FROM routes r
		LEFT JOIN warehouses w ON r.warehouse_id = w.id
		LEFT JOIN employees e ON r.driver_id = e.id
		%s
		ORDER BY r.route_code`, whereClause)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var routes []models.RouteWithDetails
	for rows.Next() {
		var route models.RouteWithDetails
		var desc, depTime *string

		err := rows.Scan(
			&route.Route.ID, &route.Route.RouteCode, &route.Route.Name, &desc,
			&route.Route.WarehouseID, &route.Route.DriverID, &route.Route.VehicleID,
			&depTime, &route.Route.IsActive, &route.Route.CreatedAt,
			&route.WarehouseName, &route.DriverName, &route.StopCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan route: %w", err)
		}
		if desc != nil {
			route.Route.Description = *desc
		}
		if depTime != nil {
			route.Route.DepartureTime = *depTime
		}
		routes = append(routes, route)
	}

	return routes, nil
}

// ============================================
// Route Stops
// ============================================

func (s *pickingServiceImpl) AddRouteStop(ctx context.Context, routeID int, req *models.AddRouteStopRequest) (int, error) {
	query := `
		INSERT INTO route_stops (route_id, customer_id, ship_to_id, stop_sequence, estimated_arrival, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	var estArr *string
	if req.EstimatedArrival != "" {
		estArr = &req.EstimatedArrival
	}

	var id int
	err := s.db.QueryRow(ctx, query,
		routeID, req.CustomerID, req.ShipToID, req.StopSequence, estArr, req.Notes,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to add route stop: %w", err)
	}
	return id, nil
}

func (s *pickingServiceImpl) UpdateRouteStop(ctx context.Context, stopID int, req *models.AddRouteStopRequest) error {
	query := `
		UPDATE route_stops SET
			customer_id = $1, ship_to_id = $2, stop_sequence = $3,
			estimated_arrival = $4, notes = $5
		WHERE id = $6`

	var estArr *string
	if req.EstimatedArrival != "" {
		estArr = &req.EstimatedArrival
	}

	result, err := s.db.Exec(ctx, query,
		req.CustomerID, req.ShipToID, req.StopSequence, estArr, req.Notes, stopID,
	)
	if err != nil {
		return fmt.Errorf("failed to update route stop: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("route stop not found")
	}
	return nil
}

func (s *pickingServiceImpl) DeleteRouteStop(ctx context.Context, stopID int) error {
	result, err := s.db.Exec(ctx, `DELETE FROM route_stops WHERE id = $1`, stopID)
	if err != nil {
		return fmt.Errorf("failed to delete route stop: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("route stop not found")
	}
	return nil
}

func (s *pickingServiceImpl) ReorderStops(ctx context.Context, routeID int, stopOrders []int) error {
	for i, stopID := range stopOrders {
		_, err := s.db.Exec(ctx,
			`UPDATE route_stops SET stop_sequence = $1 WHERE id = $2 AND route_id = $3`,
			i+1, stopID, routeID,
		)
		if err != nil {
			return fmt.Errorf("failed to reorder stops: %w", err)
		}
	}
	return nil
}

// ============================================
// Pick Lists
// ============================================

func (s *pickingServiceImpl) CreatePickList(ctx context.Context, req *models.CreatePickListRequest, createdBy int) (int, error) {
	pickNumber := s.generatePickNumber(ctx)

	// Parse pick date
	pickDate, _ := time.Parse("2006-01-02", req.PickDate)

	query := `
		INSERT INTO pick_lists (pick_number, warehouse_id, route_id, pick_date, status, created_by)
		VALUES ($1, $2, $3, $4, 'PENDING', $5)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		pickNumber, req.WarehouseID, req.RouteID, pickDate, createdBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create pick list: %w", err)
	}

	// Add lines from orders
	for _, orderID := range req.OrderIDs {
		err := s.addOrderToPickList(ctx, id, orderID)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (s *pickingServiceImpl) GeneratePickListForRoute(ctx context.Context, req *models.GeneratePickListRequest, createdBy int) (int, error) {
	// Find all confirmed orders for this route and date
	orderQuery := `
		SELECT id FROM sales_orders
		WHERE warehouse_id = $1 AND status = 'CONFIRMED'
		AND requested_ship_date = $2`

	args := []interface{}{req.WarehouseID, req.PickDate}

	if req.RouteID != nil {
		orderQuery += " AND route_id = $3"
		args = append(args, *req.RouteID)
	}

	rows := s.db.Query(ctx, orderQuery, args...)
	defer rows.Close()

	var orderIDs []int
	for rows.Next() {
		var orderID int
		rows.Scan(&orderID)
		orderIDs = append(orderIDs, orderID)
	}

	if len(orderIDs) == 0 {
		return 0, fmt.Errorf("no orders found for the specified criteria")
	}

	createReq := &models.CreatePickListRequest{
		WarehouseID: req.WarehouseID,
		RouteID:     req.RouteID,
		PickDate:    req.PickDate,
		OrderIDs:    orderIDs,
	}

	return s.CreatePickList(ctx, createReq, createdBy)
}

func (s *pickingServiceImpl) addOrderToPickList(ctx context.Context, pickListID, orderID int) error {
	query := `
		INSERT INTO pick_list_lines (pick_list_id, order_id, order_line_id, product_id, quantity_ordered)
		SELECT $1, sol.order_id, sol.id, sol.product_id, sol.quantity_ordered - sol.quantity_shipped
		FROM sales_order_lines sol
		WHERE sol.order_id = $2 AND sol.quantity_ordered > sol.quantity_shipped`

	_, err := s.db.Exec(ctx, query, pickListID, orderID)
	if err != nil {
		return fmt.Errorf("failed to add order to pick list: %w", err)
	}
	return nil
}

func (s *pickingServiceImpl) GetPickList(ctx context.Context, id int) (*models.PickListWithDetails, error) {
	query := `
		SELECT pl.id, pl.pick_number, pl.warehouse_id, pl.route_id, pl.pick_date,
			   pl.status, pl.picker_id, pl.started_at, pl.completed_at, pl.created_by, pl.created_at,
			   w.name as warehouse_name,
			   COALESCE(r.name, '') as route_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as picker_name,
			   (SELECT COUNT(DISTINCT order_id) FROM pick_list_lines WHERE pick_list_id = pl.id) as total_orders,
			   (SELECT COUNT(*) FROM pick_list_lines WHERE pick_list_id = pl.id) as total_lines,
			   (SELECT COUNT(*) FROM pick_list_lines WHERE pick_list_id = pl.id AND quantity_picked > 0) as total_picked
		FROM pick_lists pl
		JOIN warehouses w ON pl.warehouse_id = w.id
		LEFT JOIN routes r ON pl.route_id = r.id
		LEFT JOIN employees e ON pl.picker_id = e.id
		WHERE pl.id = $1`

	var pickList models.PickListWithDetails
	var startedAt, completedAt *time.Time

	err := s.db.QueryRow(ctx, query, id).Scan(
		&pickList.PickList.ID, &pickList.PickList.PickNumber, &pickList.PickList.WarehouseID,
		&pickList.PickList.RouteID, &pickList.PickList.PickDate, &pickList.PickList.Status,
		&pickList.PickList.PickerID, &startedAt, &completedAt,
		&pickList.PickList.CreatedBy, &pickList.PickList.CreatedAt,
		&pickList.WarehouseName, &pickList.RouteName, &pickList.PickerName,
		&pickList.TotalOrders, &pickList.TotalLines, &pickList.TotalPicked,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("pick list not found")
		}
		return nil, fmt.Errorf("failed to get pick list: %w", err)
	}

	return &pickList, nil
}

func (s *pickingServiceImpl) ListPickLists(ctx context.Context, filters *models.PickListFilters) ([]models.PickListWithDetails, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 20
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if filters.WarehouseID != nil {
		whereClause += fmt.Sprintf(" AND pl.warehouse_id = $%d", argNum)
		args = append(args, *filters.WarehouseID)
		argNum++
	}
	if filters.RouteID != nil {
		whereClause += fmt.Sprintf(" AND pl.route_id = $%d", argNum)
		args = append(args, *filters.RouteID)
		argNum++
	}
	if filters.Status != nil {
		whereClause += fmt.Sprintf(" AND pl.status = $%d", argNum)
		args = append(args, *filters.Status)
		argNum++
	}
	if filters.PickerID != nil {
		whereClause += fmt.Sprintf(" AND pl.picker_id = $%d", argNum)
		args = append(args, *filters.PickerID)
		argNum++
	}
	if filters.DateFrom != "" {
		whereClause += fmt.Sprintf(" AND pl.pick_date >= $%d", argNum)
		args = append(args, filters.DateFrom)
		argNum++
	}
	if filters.DateTo != "" {
		whereClause += fmt.Sprintf(" AND pl.pick_date <= $%d", argNum)
		args = append(args, filters.DateTo)
		argNum++
	}

	// Count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM pick_lists pl %s`, whereClause)
	var total int64
	s.db.QueryRow(ctx, countQuery, args...).Scan(&total)

	// Query
	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT pl.id, pl.pick_number, pl.warehouse_id, pl.route_id, pl.pick_date,
			   pl.status, pl.picker_id, pl.started_at, pl.completed_at, pl.created_by, pl.created_at,
			   w.name as warehouse_name,
			   COALESCE(r.name, '') as route_name,
			   COALESCE(e.first_name || ' ' || e.last_name, '') as picker_name,
			   (SELECT COUNT(DISTINCT order_id) FROM pick_list_lines WHERE pick_list_id = pl.id) as total_orders,
			   (SELECT COUNT(*) FROM pick_list_lines WHERE pick_list_id = pl.id) as total_lines,
			   (SELECT COUNT(*) FROM pick_list_lines WHERE pick_list_id = pl.id AND quantity_picked > 0) as total_picked
		FROM pick_lists pl
		JOIN warehouses w ON pl.warehouse_id = w.id
		LEFT JOIN routes r ON pl.route_id = r.id
		LEFT JOIN employees e ON pl.picker_id = e.id
		%s
		ORDER BY pl.pick_date DESC, pl.id DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var pickLists []models.PickListWithDetails
	for rows.Next() {
		var pickList models.PickListWithDetails
		var startedAt, completedAt *time.Time

		err := rows.Scan(
			&pickList.PickList.ID, &pickList.PickList.PickNumber, &pickList.PickList.WarehouseID,
			&pickList.PickList.RouteID, &pickList.PickList.PickDate, &pickList.PickList.Status,
			&pickList.PickList.PickerID, &startedAt, &completedAt,
			&pickList.PickList.CreatedBy, &pickList.PickList.CreatedAt,
			&pickList.WarehouseName, &pickList.RouteName, &pickList.PickerName,
			&pickList.TotalOrders, &pickList.TotalLines, &pickList.TotalPicked,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pick list: %w", err)
		}
		pickLists = append(pickLists, pickList)
	}

	return pickLists, total, nil
}

func (s *pickingServiceImpl) StartPicking(ctx context.Context, pickListID int, pickerID int) error {
	result, err := s.db.Exec(ctx,
		`UPDATE pick_lists SET status = 'IN_PROGRESS', picker_id = $1, started_at = NOW()
		 WHERE id = $2 AND status = 'PENDING'`,
		pickerID, pickListID,
	)
	if err != nil {
		return fmt.Errorf("failed to start picking: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("pick list not found or already started")
	}
	return nil
}

func (s *pickingServiceImpl) CompletePicking(ctx context.Context, pickListID int) error {
	result, err := s.db.Exec(ctx,
		`UPDATE pick_lists SET status = 'COMPLETE', completed_at = NOW()
		 WHERE id = $1 AND status = 'IN_PROGRESS'`,
		pickListID,
	)
	if err != nil {
		return fmt.Errorf("failed to complete picking: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("pick list not found or not in progress")
	}

	// Update sales order lines with picked quantities
	s.db.Exec(ctx, `
		UPDATE sales_order_lines sol
		SET quantity_shipped = quantity_shipped + pll.quantity_picked
		FROM pick_list_lines pll
		WHERE pll.order_line_id = sol.id AND pll.pick_list_id = $1`, pickListID)

	return nil
}

func (s *pickingServiceImpl) CancelPickList(ctx context.Context, pickListID int) error {
	result, err := s.db.Exec(ctx,
		`UPDATE pick_lists SET status = 'CANCELLED' WHERE id = $1 AND status IN ('PENDING', 'IN_PROGRESS')`,
		pickListID,
	)
	if err != nil {
		return fmt.Errorf("failed to cancel pick list: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("pick list not found or cannot be cancelled")
	}
	return nil
}

// ============================================
// Pick Lines
// ============================================

func (s *pickingServiceImpl) GetPickLines(ctx context.Context, pickListID int) ([]models.PickListLineWithProduct, error) {
	query := `
		SELECT pll.id, pll.pick_list_id, pll.order_id, pll.order_line_id, pll.product_id,
			   pll.location_code, pll.lot_number, pll.quantity_ordered, pll.quantity_picked,
			   pll.catch_weight, pll.picked_at, pll.picked_by, pll.notes,
			   so.order_number, c.name as customer_name,
			   p.sku as product_sku, p.name as product_name,
			   COALESCE(pu.abbreviation, 'EA') as unit_of_measure,
			   sol.expiry_date
		FROM pick_list_lines pll
		JOIN sales_orders so ON pll.order_id = so.id
		JOIN customers c ON so.customer_id = c.id
		JOIN products p ON pll.product_id = p.id
		JOIN sales_order_lines sol ON pll.order_line_id = sol.id
		LEFT JOIN product_units pu ON p.default_unit_id = pu.id
		WHERE pll.pick_list_id = $1
		ORDER BY p.name, pll.id`

	rows := s.db.Query(ctx, query, pickListID)
	defer rows.Close()

	var lines []models.PickListLineWithProduct
	for rows.Next() {
		var line models.PickListLineWithProduct
		var locCode, lotNum, notes, pickedAt *string
		var expDate *time.Time

		err := rows.Scan(
			&line.PickListLine.ID, &line.PickListLine.PickListID, &line.PickListLine.OrderID,
			&line.PickListLine.OrderLineID, &line.PickListLine.ProductID,
			&locCode, &lotNum, &line.PickListLine.QuantityOrdered, &line.PickListLine.QuantityPicked,
			&line.PickListLine.CatchWeight, &pickedAt, &line.PickListLine.PickedBy, &notes,
			&line.OrderNumber, &line.CustomerName, &line.ProductSKU, &line.ProductName,
			&line.UnitOfMeasure, &expDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pick line: %w", err)
		}
		if locCode != nil {
			line.PickListLine.LocationCode = *locCode
		}
		if lotNum != nil {
			line.PickListLine.LotNumber = *lotNum
		}
		if notes != nil {
			line.PickListLine.Notes = *notes
		}
		if pickedAt != nil {
			line.PickListLine.PickedAt = *pickedAt
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func (s *pickingServiceImpl) ConfirmPickLine(ctx context.Context, lineID int, req *models.ConfirmPickLineRequest, pickerID int) error {
	query := `
		UPDATE pick_list_lines SET
			quantity_picked = $1, catch_weight = $2, lot_number = $3,
			location_code = $4, notes = $5, picked_by = $6, picked_at = NOW()
		WHERE id = $7`

	result, err := s.db.Exec(ctx, query,
		req.QuantityPicked, req.CatchWeight, req.LotNumber,
		req.LocationCode, req.Notes, pickerID, lineID,
	)
	if err != nil {
		return fmt.Errorf("failed to confirm pick line: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("pick line not found")
	}
	return nil
}

// ============================================
// Master Pick Report
// ============================================

func (s *pickingServiceImpl) GetMasterPickReport(ctx context.Context, warehouseID int, pickDate string, routeID *int) ([]models.MasterPickItem, error) {
	whereClause := "WHERE pl.warehouse_id = $1 AND pl.pick_date = $2 AND pl.status != 'CANCELLED'"
	args := []interface{}{warehouseID, pickDate}
	argNum := 3

	if routeID != nil {
		whereClause += fmt.Sprintf(" AND pl.route_id = $%d", argNum)
		args = append(args, *routeID)
	}

	query := fmt.Sprintf(`
		SELECT pll.product_id, p.sku, p.name, COALESCE(pc.name, '') as category_name,
			   COALESCE(pll.location_code, '') as location_code,
			   SUM(pll.quantity_ordered) as total_quantity,
			   COALESCE(pu.abbreviation, 'EA') as unit_of_measure,
			   COUNT(DISTINCT pll.order_id) as order_count,
			   COALESCE(pll.lot_number, '') as lot_number,
			   '' as expiry_date
		FROM pick_list_lines pll
		JOIN pick_lists pl ON pll.pick_list_id = pl.id
		JOIN products p ON pll.product_id = p.id
		LEFT JOIN product_categories pc ON p.category_id = pc.id
		LEFT JOIN product_units pu ON p.default_unit_id = pu.id
		%s
		GROUP BY pll.product_id, p.sku, p.name, pc.name, pll.location_code, pu.abbreviation, pll.lot_number
		ORDER BY pc.name, p.name, pll.location_code`, whereClause)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var items []models.MasterPickItem
	for rows.Next() {
		var item models.MasterPickItem
		err := rows.Scan(
			&item.ProductID, &item.ProductSKU, &item.ProductName, &item.CategoryName,
			&item.LocationCode, &item.TotalQuantity, &item.UnitOfMeasure, &item.OrderCount,
			&item.LotNumber, &item.ExpiryDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan master pick item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// ============================================
// Suggested Picking (FEFO)
// ============================================

func (s *pickingServiceImpl) GetSuggestedPicking(ctx context.Context, productID, warehouseID int, quantity float64) ([]SuggestedPickLocation, error) {
	query := `
		SELECT location_code, lot_number, expiry_date, quantity_available
		FROM inventory
		WHERE product_id = $1 AND warehouse_id = $2 AND quantity_available > 0
		ORDER BY expiry_date NULLS LAST, production_date, lot_number`

	rows := s.db.Query(ctx, query, productID, warehouseID)
	defer rows.Close()

	var suggestions []SuggestedPickLocation
	remaining := quantity

	for rows.Next() && remaining > 0 {
		var loc SuggestedPickLocation
		var locCode, lotNum *string

		err := rows.Scan(&locCode, &lotNum, &loc.ExpiryDate, &loc.Available)
		if err != nil {
			return nil, fmt.Errorf("failed to scan suggestion: %w", err)
		}

		if locCode != nil {
			loc.LocationCode = *locCode
		}
		if lotNum != nil {
			loc.LotNumber = *lotNum
		}

		if loc.Available >= remaining {
			loc.Suggested = remaining
			remaining = 0
		} else {
			loc.Suggested = loc.Available
			remaining -= loc.Available
		}

		suggestions = append(suggestions, loc)
	}

	return suggestions, nil
}

// ============================================
// Helpers
// ============================================

func (s *pickingServiceImpl) generatePickNumber(ctx context.Context) string {
	var count int64
	s.db.QueryRow(ctx, `SELECT COUNT(*) FROM pick_lists WHERE pick_date = CURRENT_DATE`).Scan(&count)
	return fmt.Sprintf("PK%s%04d", time.Now().Format("20060102"), count+1)
}
