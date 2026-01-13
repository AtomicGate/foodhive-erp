package inventory

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

type InventoryService interface {
	// Inventory Queries
	GetByID(ctx context.Context, id int) (*models.InventoryWithDetails, error)
	GetByProduct(ctx context.Context, productID int) ([]models.Inventory, error)
	GetByWarehouse(ctx context.Context, warehouseID int) ([]models.Inventory, error)
	GetByLocation(ctx context.Context, warehouseID int, locationCode string) ([]models.Inventory, error)
	GetByLot(ctx context.Context, lotNumber string) ([]models.Inventory, error)
	List(ctx context.Context, filters *models.InventoryInquiryRequest) ([]models.InventoryWithDetails, int64, error)

	// Inventory Summary
	GetProductSummary(ctx context.Context, productID int) (*models.InventorySummary, error)
	GetExpiringInventory(ctx context.Context, daysToExpiry int, warehouseID *int) ([]models.InventoryWithDetails, error)

	// Inventory Operations
	Receive(ctx context.Context, req *ReceiveRequest, createdBy int) (int, error)
	Adjust(ctx context.Context, req *models.AdjustInventoryRequest, createdBy int) error
	Transfer(ctx context.Context, req *models.TransferInventoryRequest, createdBy int) error

	// Transaction History
	GetTransactions(ctx context.Context, productID, warehouseID *int, limit int) ([]models.InventoryTransaction, error)
}

// ============================================
// Additional Request Types
// ============================================

type ReceiveRequest struct {
	ProductID       int        `json:"product_id"`
	WarehouseID     int        `json:"warehouse_id"`
	LocationCode    string     `json:"location_code,omitempty"`
	LotNumber       string     `json:"lot_number,omitempty"`
	ProductionDate  *time.Time `json:"production_date,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	Quantity        float64    `json:"quantity"`
	UnitCost        float64    `json:"unit_cost"`
	ReferenceType   string     `json:"reference_type,omitempty"`
	ReferenceID     int        `json:"reference_id,omitempty"`
	ReferenceNumber string     `json:"reference_number,omitempty"`
	Notes           string     `json:"notes,omitempty"`
}

// ============================================
// Service Implementation
// ============================================

type inventoryServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) InventoryService {
	return &inventoryServiceImpl{db: db}
}

// ============================================
// Inventory Queries
// ============================================

func (s *inventoryServiceImpl) GetByID(ctx context.Context, id int) (*models.InventoryWithDetails, error) {
	query := `
		SELECT i.id, i.product_id, i.warehouse_id, i.location_code, i.lot_number,
			   i.production_date, i.expiry_date, i.quantity_on_hand, i.quantity_allocated,
			   i.quantity_on_order, i.quantity_available, i.last_cost, i.average_cost,
			   i.last_counted_date, i.last_movement_date, i.created_at, i.updated_at,
			   p.name as product_name, p.sku as product_sku, w.name as warehouse_name,
			   COALESCE(i.expiry_date - CURRENT_DATE, 0) as days_to_expiry,
			   COALESCE(CURRENT_DATE - i.production_date::date, 0) as age_in_days
		FROM inventory i
		JOIN products p ON i.product_id = p.id
		JOIN warehouses w ON i.warehouse_id = w.id
		WHERE i.id = $1`

	var inv models.InventoryWithDetails
	var locCode, lotNum *string
	var prodDate, expDate, countDate, moveDate *time.Time

	err := s.db.QueryRow(ctx, query, id).Scan(
		&inv.Inventory.ID, &inv.Inventory.ProductID, &inv.Inventory.WarehouseID,
		&locCode, &lotNum, &prodDate, &expDate,
		&inv.Inventory.QuantityOnHand, &inv.Inventory.QuantityAllocated,
		&inv.Inventory.QuantityOnOrder, &inv.Inventory.QuantityAvailable,
		&inv.Inventory.LastCost, &inv.Inventory.AverageCost,
		&countDate, &moveDate, &inv.Inventory.CreatedAt, &inv.Inventory.UpdatedAt,
		&inv.ProductName, &inv.ProductSKU, &inv.WarehouseName,
		&inv.DaysToExpiry, &inv.AgeInDays,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("inventory not found")
		}
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	if locCode != nil {
		inv.Inventory.LocationCode = *locCode
	}
	if lotNum != nil {
		inv.Inventory.LotNumber = *lotNum
	}

	return &inv, nil
}

func (s *inventoryServiceImpl) GetByProduct(ctx context.Context, productID int) ([]models.Inventory, error) {
	return s.getInventoryList(ctx, "product_id = $1", productID)
}

func (s *inventoryServiceImpl) GetByWarehouse(ctx context.Context, warehouseID int) ([]models.Inventory, error) {
	return s.getInventoryList(ctx, "warehouse_id = $1", warehouseID)
}

func (s *inventoryServiceImpl) GetByLocation(ctx context.Context, warehouseID int, locationCode string) ([]models.Inventory, error) {
	query := `
		SELECT id, product_id, warehouse_id, location_code, lot_number,
			   production_date, expiry_date, quantity_on_hand, quantity_allocated,
			   quantity_on_order, quantity_available, last_cost, average_cost,
			   last_counted_date, last_movement_date, created_at, updated_at
		FROM inventory
		WHERE warehouse_id = $1 AND location_code = $2
		ORDER BY product_id, lot_number`

	rows := s.db.Query(ctx, query, warehouseID, locationCode)
	defer rows.Close()

	return s.scanInventoryRows(rows)
}

func (s *inventoryServiceImpl) GetByLot(ctx context.Context, lotNumber string) ([]models.Inventory, error) {
	return s.getInventoryList(ctx, "lot_number = $1", lotNumber)
}

func (s *inventoryServiceImpl) getInventoryList(ctx context.Context, whereClause string, arg interface{}) ([]models.Inventory, error) {
	query := fmt.Sprintf(`
		SELECT id, product_id, warehouse_id, location_code, lot_number,
			   production_date, expiry_date, quantity_on_hand, quantity_allocated,
			   quantity_on_order, quantity_available, last_cost, average_cost,
			   last_counted_date, last_movement_date, created_at, updated_at
		FROM inventory
		WHERE %s
		ORDER BY warehouse_id, location_code, lot_number`, whereClause)

	rows := s.db.Query(ctx, query, arg)
	defer rows.Close()

	return s.scanInventoryRows(rows)
}

func (s *inventoryServiceImpl) scanInventoryRows(rows postgres.Rows) ([]models.Inventory, error) {
	var inventories []models.Inventory

	for rows.Next() {
		var inv models.Inventory
		var locCode, lotNum *string
		var prodDate, expDate, countDate, moveDate *time.Time

		err := rows.Scan(
			&inv.ID, &inv.ProductID, &inv.WarehouseID,
			&locCode, &lotNum, &prodDate, &expDate,
			&inv.QuantityOnHand, &inv.QuantityAllocated,
			&inv.QuantityOnOrder, &inv.QuantityAvailable,
			&inv.LastCost, &inv.AverageCost,
			&countDate, &moveDate, &inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}

		if locCode != nil {
			inv.LocationCode = *locCode
		}
		if lotNum != nil {
			inv.LotNumber = *lotNum
		}

		inventories = append(inventories, inv)
	}

	return inventories, nil
}

func (s *inventoryServiceImpl) List(ctx context.Context, filters *models.InventoryInquiryRequest) ([]models.InventoryWithDetails, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 20
	}

	whereClause := "WHERE i.quantity_on_hand > 0"
	args := []interface{}{}
	argNum := 1

	if filters.ProductID != nil {
		whereClause += fmt.Sprintf(" AND i.product_id = $%d", argNum)
		args = append(args, *filters.ProductID)
		argNum++
	}

	if filters.WarehouseID != nil {
		whereClause += fmt.Sprintf(" AND i.warehouse_id = $%d", argNum)
		args = append(args, *filters.WarehouseID)
		argNum++
	}

	if filters.LotNumber != "" {
		whereClause += fmt.Sprintf(" AND i.lot_number = $%d", argNum)
		args = append(args, filters.LotNumber)
		argNum++
	}

	if filters.ShowExpiring && filters.DaysToExpiry > 0 {
		whereClause += fmt.Sprintf(" AND i.expiry_date <= CURRENT_DATE + $%d", argNum)
		args = append(args, filters.DaysToExpiry)
		argNum++
	}

	// Count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM inventory i
		JOIN products p ON i.product_id = p.id
		%s`, whereClause)
	var total int64
	s.db.QueryRow(ctx, countQuery, args...).Scan(&total)

	// Query
	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT i.id, i.product_id, i.warehouse_id, i.location_code, i.lot_number,
			   i.production_date, i.expiry_date, i.quantity_on_hand, i.quantity_allocated,
			   i.quantity_on_order, i.quantity_available, i.last_cost, i.average_cost,
			   i.last_counted_date, i.last_movement_date, i.created_at, i.updated_at,
			   p.name as product_name, p.sku as product_sku, w.name as warehouse_name,
			   COALESCE(i.expiry_date - CURRENT_DATE, 0) as days_to_expiry,
			   COALESCE(CURRENT_DATE - i.production_date::date, 0) as age_in_days
		FROM inventory i
		JOIN products p ON i.product_id = p.id
		JOIN warehouses w ON i.warehouse_id = w.id
		%s
		ORDER BY p.name, i.expiry_date NULLS LAST
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var inventories []models.InventoryWithDetails
	for rows.Next() {
		var inv models.InventoryWithDetails
		var locCode, lotNum *string
		var prodDate, expDate, countDate, moveDate *time.Time

		err := rows.Scan(
			&inv.Inventory.ID, &inv.Inventory.ProductID, &inv.Inventory.WarehouseID,
			&locCode, &lotNum, &prodDate, &expDate,
			&inv.Inventory.QuantityOnHand, &inv.Inventory.QuantityAllocated,
			&inv.Inventory.QuantityOnOrder, &inv.Inventory.QuantityAvailable,
			&inv.Inventory.LastCost, &inv.Inventory.AverageCost,
			&countDate, &moveDate, &inv.Inventory.CreatedAt, &inv.Inventory.UpdatedAt,
			&inv.ProductName, &inv.ProductSKU, &inv.WarehouseName,
			&inv.DaysToExpiry, &inv.AgeInDays,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory: %w", err)
		}

		if locCode != nil {
			inv.Inventory.LocationCode = *locCode
		}
		if lotNum != nil {
			inv.Inventory.LotNumber = *lotNum
		}

		inventories = append(inventories, inv)
	}

	return inventories, total, nil
}

// ============================================
// Inventory Summary
// ============================================

func (s *inventoryServiceImpl) GetProductSummary(ctx context.Context, productID int) (*models.InventorySummary, error) {
	// Get product info and totals
	query := `
		SELECT p.id, p.sku, p.name,
			   COALESCE(SUM(i.quantity_on_hand), 0) as total_on_hand,
			   COALESCE(SUM(i.quantity_allocated), 0) as total_allocated,
			   COALESCE(SUM(i.quantity_on_order), 0) as total_on_order,
			   COALESCE(SUM(i.quantity_available), 0) as total_available,
			   COALESCE(AVG(i.average_cost), 0) as avg_cost
		FROM products p
		LEFT JOIN inventory i ON p.id = i.product_id
		WHERE p.id = $1
		GROUP BY p.id, p.sku, p.name`

	var summary models.InventorySummary
	err := s.db.QueryRow(ctx, query, productID).Scan(
		&summary.ProductID, &summary.ProductSKU, &summary.ProductName,
		&summary.TotalOnHand, &summary.TotalAllocated, &summary.TotalOnOrder,
		&summary.TotalAvailable, &summary.AverageCost,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product summary: %w", err)
	}

	summary.InventoryValue = summary.TotalOnHand * summary.AverageCost

	// Get warehouse breakdown
	whQuery := `
		SELECT i.warehouse_id, w.name,
			   SUM(i.quantity_on_hand) as on_hand,
			   SUM(i.quantity_allocated) as allocated,
			   SUM(i.quantity_available) as available
		FROM inventory i
		JOIN warehouses w ON i.warehouse_id = w.id
		WHERE i.product_id = $1
		GROUP BY i.warehouse_id, w.name
		ORDER BY w.name`

	rows := s.db.Query(ctx, whQuery, productID)
	defer rows.Close()

	for rows.Next() {
		var wq models.WarehouseQuantity
		err := rows.Scan(&wq.WarehouseID, &wq.WarehouseName, &wq.OnHand, &wq.Allocated, &wq.Available)
		if err != nil {
			return nil, fmt.Errorf("failed to scan warehouse quantity: %w", err)
		}
		summary.WarehouseBreakdown = append(summary.WarehouseBreakdown, wq)
	}

	return &summary, nil
}

func (s *inventoryServiceImpl) GetExpiringInventory(ctx context.Context, daysToExpiry int, warehouseID *int) ([]models.InventoryWithDetails, error) {
	whereClause := "WHERE i.expiry_date IS NOT NULL AND i.expiry_date <= CURRENT_DATE + $1 AND i.quantity_on_hand > 0"
	args := []interface{}{daysToExpiry}
	argNum := 2

	if warehouseID != nil {
		whereClause += fmt.Sprintf(" AND i.warehouse_id = $%d", argNum)
		args = append(args, *warehouseID)
	}

	query := fmt.Sprintf(`
		SELECT i.id, i.product_id, i.warehouse_id, i.location_code, i.lot_number,
			   i.production_date, i.expiry_date, i.quantity_on_hand, i.quantity_allocated,
			   i.quantity_on_order, i.quantity_available, i.last_cost, i.average_cost,
			   i.last_counted_date, i.last_movement_date, i.created_at, i.updated_at,
			   p.name as product_name, p.sku as product_sku, w.name as warehouse_name,
			   i.expiry_date - CURRENT_DATE as days_to_expiry,
			   COALESCE(CURRENT_DATE - i.production_date::date, 0) as age_in_days
		FROM inventory i
		JOIN products p ON i.product_id = p.id
		JOIN warehouses w ON i.warehouse_id = w.id
		%s
		ORDER BY i.expiry_date ASC`, whereClause)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var inventories []models.InventoryWithDetails
	for rows.Next() {
		var inv models.InventoryWithDetails
		var locCode, lotNum *string
		var prodDate, expDate, countDate, moveDate *time.Time

		err := rows.Scan(
			&inv.Inventory.ID, &inv.Inventory.ProductID, &inv.Inventory.WarehouseID,
			&locCode, &lotNum, &prodDate, &expDate,
			&inv.Inventory.QuantityOnHand, &inv.Inventory.QuantityAllocated,
			&inv.Inventory.QuantityOnOrder, &inv.Inventory.QuantityAvailable,
			&inv.Inventory.LastCost, &inv.Inventory.AverageCost,
			&countDate, &moveDate, &inv.Inventory.CreatedAt, &inv.Inventory.UpdatedAt,
			&inv.ProductName, &inv.ProductSKU, &inv.WarehouseName,
			&inv.DaysToExpiry, &inv.AgeInDays,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}

		if locCode != nil {
			inv.Inventory.LocationCode = *locCode
		}
		if lotNum != nil {
			inv.Inventory.LotNumber = *lotNum
		}

		inventories = append(inventories, inv)
	}

	return inventories, nil
}

// ============================================
// Inventory Operations
// ============================================

func (s *inventoryServiceImpl) Receive(ctx context.Context, req *ReceiveRequest, createdBy int) (int, error) {
	// Upsert inventory record
	query := `
		INSERT INTO inventory (
			product_id, warehouse_id, location_code, lot_number,
			production_date, expiry_date, quantity_on_hand, last_cost, average_cost,
			last_movement_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8, NOW())
		ON CONFLICT (product_id, warehouse_id, location_code, lot_number)
		DO UPDATE SET
			quantity_on_hand = inventory.quantity_on_hand + $7,
			last_cost = $8,
			average_cost = (inventory.average_cost * inventory.quantity_on_hand + $8 * $7) / (inventory.quantity_on_hand + $7),
			last_movement_date = NOW(),
			updated_at = NOW()
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.ProductID, req.WarehouseID, req.LocationCode, req.LotNumber,
		req.ProductionDate, req.ExpiryDate, req.Quantity, req.UnitCost,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to receive inventory: %w", err)
	}

	// Log transaction
	s.logTransaction(ctx, req.ProductID, req.WarehouseID, req.LocationCode,
		models.TxReceive, req.Quantity, req.LotNumber, req.UnitCost,
		req.ReferenceType, req.ReferenceID, req.ReferenceNumber, req.Notes, createdBy)

	return id, nil
}

func (s *inventoryServiceImpl) Adjust(ctx context.Context, req *models.AdjustInventoryRequest, createdBy int) error {
	txType := models.TxAdjustIn
	if req.Quantity < 0 {
		txType = models.TxAdjustOut
	}

	query := `
		UPDATE inventory SET
			quantity_on_hand = quantity_on_hand + $1,
			last_movement_date = NOW(),
			updated_at = NOW()
		WHERE product_id = $2 AND warehouse_id = $3
			AND COALESCE(location_code, '') = COALESCE($4, '')
			AND COALESCE(lot_number, '') = COALESCE($5, '')`

	result, err := s.db.Exec(ctx, query,
		req.Quantity, req.ProductID, req.WarehouseID, req.LocationCode, req.LotNumber,
	)

	if err != nil {
		return fmt.Errorf("failed to adjust inventory: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("inventory record not found")
	}

	// Log transaction
	s.logTransaction(ctx, req.ProductID, req.WarehouseID, req.LocationCode,
		txType, req.Quantity, req.LotNumber, 0, "ADJUSTMENT", 0, "", req.Reason+": "+req.Notes, createdBy)

	return nil
}

func (s *inventoryServiceImpl) Transfer(ctx context.Context, req *models.TransferInventoryRequest, createdBy int) error {
	// Deduct from source
	deductQuery := `
		UPDATE inventory SET
			quantity_on_hand = quantity_on_hand - $1,
			last_movement_date = NOW(),
			updated_at = NOW()
		WHERE product_id = $2 AND warehouse_id = $3
			AND COALESCE(location_code, '') = COALESCE($4, '')
			AND COALESCE(lot_number, '') = COALESCE($5, '')
			AND quantity_on_hand >= $1`

	result, err := s.db.Exec(ctx, deductQuery,
		req.Quantity, req.ProductID, req.FromWarehouseID, req.FromLocationCode, req.LotNumber,
	)

	if err != nil {
		return fmt.Errorf("failed to deduct from source: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("insufficient quantity or inventory not found")
	}

	// Get cost for the transferred item
	var cost float64
	s.db.QueryRow(ctx, `SELECT COALESCE(average_cost, 0) FROM inventory WHERE product_id = $1 AND warehouse_id = $2 LIMIT 1`,
		req.ProductID, req.FromWarehouseID).Scan(&cost)

	// Add to destination
	addQuery := `
		INSERT INTO inventory (
			product_id, warehouse_id, location_code, lot_number, quantity_on_hand, average_cost, last_movement_date
		) VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (product_id, warehouse_id, location_code, lot_number)
		DO UPDATE SET
			quantity_on_hand = inventory.quantity_on_hand + $5,
			last_movement_date = NOW(),
			updated_at = NOW()`

	_, err = s.db.Exec(ctx, addQuery,
		req.ProductID, req.ToWarehouseID, req.ToLocationCode, req.LotNumber, req.Quantity, cost,
	)

	if err != nil {
		return fmt.Errorf("failed to add to destination: %w", err)
	}

	// Log transactions
	s.logTransaction(ctx, req.ProductID, req.FromWarehouseID, req.FromLocationCode,
		models.TxTransferOut, -req.Quantity, req.LotNumber, cost, "TRANSFER", 0, "", req.Notes, createdBy)

	s.logTransaction(ctx, req.ProductID, req.ToWarehouseID, req.ToLocationCode,
		models.TxTransferIn, req.Quantity, req.LotNumber, cost, "TRANSFER", 0, "", req.Notes, createdBy)

	return nil
}

// ============================================
// Transaction History
// ============================================

func (s *inventoryServiceImpl) GetTransactions(ctx context.Context, productID, warehouseID *int, limit int) ([]models.InventoryTransaction, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if productID != nil {
		whereClause += fmt.Sprintf(" AND product_id = $%d", argNum)
		args = append(args, *productID)
		argNum++
	}

	if warehouseID != nil {
		whereClause += fmt.Sprintf(" AND warehouse_id = $%d", argNum)
		args = append(args, *warehouseID)
		argNum++
	}

	query := fmt.Sprintf(`
		SELECT id, product_id, warehouse_id, location_code, transaction_type,
			   quantity, lot_number, unit_cost, reference_type, reference_id,
			   reference_number, notes, created_by, created_at
		FROM inventory_transactions
		%s
		ORDER BY created_at DESC
		LIMIT $%d`, whereClause, argNum)

	args = append(args, limit)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var transactions []models.InventoryTransaction
	for rows.Next() {
		var tx models.InventoryTransaction
		var locCode, lotNum, refType, refNum, notes *string
		var refID *int

		err := rows.Scan(
			&tx.ID, &tx.ProductID, &tx.WarehouseID, &locCode, &tx.TransactionType,
			&tx.Quantity, &lotNum, &tx.UnitCost, &refType, &refID,
			&refNum, &notes, &tx.CreatedBy, &tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		if locCode != nil {
			tx.LocationCode = *locCode
		}
		if lotNum != nil {
			tx.LotNumber = *lotNum
		}
		if refType != nil {
			tx.ReferenceType = *refType
		}
		if refID != nil {
			tx.ReferenceID = *refID
		}
		if refNum != nil {
			tx.ReferenceNumber = *refNum
		}
		if notes != nil {
			tx.Notes = *notes
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// ============================================
// Helper Functions
// ============================================

func (s *inventoryServiceImpl) logTransaction(
	ctx context.Context,
	productID, warehouseID int,
	locationCode string,
	txType models.InventoryTransactionType,
	quantity float64,
	lotNumber string,
	unitCost float64,
	refType string,
	refID int,
	refNumber, notes string,
	createdBy int,
) {
	query := `
		INSERT INTO inventory_transactions (
			product_id, warehouse_id, location_code, transaction_type,
			quantity, lot_number, unit_cost, reference_type, reference_id,
			reference_number, notes, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	s.db.Exec(ctx, query,
		productID, warehouseID, locationCode, txType,
		quantity, lotNumber, unitCost, refType, refID,
		refNumber, notes, createdBy,
	)
}
