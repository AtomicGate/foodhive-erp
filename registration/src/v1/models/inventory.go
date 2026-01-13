package models

import "time"

// ============================================
// Inventory Models
// ============================================

type Inventory struct {
	ID                int        `json:"id"`
	ProductID         int        `json:"product_id"`
	WarehouseID       int        `json:"warehouse_id"`
	LocationCode      string     `json:"location_code,omitempty"`
	LotNumber         string     `json:"lot_number,omitempty"`
	ProductionDate    CustomDate `json:"production_date,omitempty"`
	ExpiryDate        CustomDate `json:"expiry_date,omitempty"`
	QuantityOnHand    float64    `json:"quantity_on_hand"`
	QuantityAllocated float64    `json:"quantity_allocated"`
	QuantityOnOrder   float64    `json:"quantity_on_order"`
	QuantityAvailable float64    `json:"quantity_available"`
	LastCost          float64    `json:"last_cost"`
	AverageCost       float64    `json:"average_cost"`
	LastCountedDate   CustomDate `json:"last_counted_date,omitempty"`
	LastMovementDate  CustomDate `json:"last_movement_date,omitempty"`
	CreatedAt         CustomDate `json:"created_at"`
	UpdatedAt         CustomDate `json:"updated_at"`
}

type InventoryWithDetails struct {
	Inventory     Inventory `json:"inventory"`
	ProductName   string    `json:"product_name"`
	ProductSKU    string    `json:"product_sku"`
	WarehouseName string    `json:"warehouse_name"`
	DaysToExpiry  int       `json:"days_to_expiry,omitempty"`
	AgeInDays     int       `json:"age_in_days"`
}

// InventoryTransactionType represents types of inventory transactions
type InventoryTransactionType string

const (
	TxReceive     InventoryTransactionType = "RECEIVE"
	TxShip        InventoryTransactionType = "SHIP"
	TxAdjustIn    InventoryTransactionType = "ADJUST_IN"
	TxAdjustOut   InventoryTransactionType = "ADJUST_OUT"
	TxTransferIn  InventoryTransactionType = "TRANSFER_IN"
	TxTransferOut InventoryTransactionType = "TRANSFER_OUT"
	TxReturn      InventoryTransactionType = "RETURN"
	TxDispose     InventoryTransactionType = "DISPOSE"
	TxCycleCount  InventoryTransactionType = "CYCLE_COUNT"
)

type InventoryTransaction struct {
	ID              int                      `json:"id"`
	ProductID       int                      `json:"product_id"`
	WarehouseID     int                      `json:"warehouse_id"`
	LocationCode    string                   `json:"location_code,omitempty"`
	TransactionType InventoryTransactionType `json:"transaction_type"`
	Quantity        float64                  `json:"quantity"`
	LotNumber       string                   `json:"lot_number,omitempty"`
	UnitCost        float64                  `json:"unit_cost"`
	ReferenceType   string                   `json:"reference_type,omitempty"`
	ReferenceID     int                      `json:"reference_id,omitempty"`
	ReferenceNumber string                   `json:"reference_number,omitempty"`
	Notes           string                   `json:"notes,omitempty"`
	CreatedBy       int                      `json:"created_by"`
	CreatedAt       time.Time                `json:"created_at"`
}

// ============================================
// Request/Response Types
// ============================================

type AdjustInventoryRequest struct {
	ProductID    int     `json:"product_id"`
	WarehouseID  int     `json:"warehouse_id"`
	LocationCode string  `json:"location_code,omitempty"`
	LotNumber    string  `json:"lot_number,omitempty"`
	Quantity     float64 `json:"quantity"`
	Reason       string  `json:"reason"`
	Notes        string  `json:"notes,omitempty"`
}

type TransferInventoryRequest struct {
	ProductID        int     `json:"product_id"`
	FromWarehouseID  int     `json:"from_warehouse_id"`
	ToWarehouseID    int     `json:"to_warehouse_id"`
	FromLocationCode string  `json:"from_location_code,omitempty"`
	ToLocationCode   string  `json:"to_location_code,omitempty"`
	LotNumber        string  `json:"lot_number,omitempty"`
	Quantity         float64 `json:"quantity"`
	Notes            string  `json:"notes,omitempty"`
}

type InventoryInquiryRequest struct {
	ProductID    *int   `json:"product_id,omitempty"`
	WarehouseID  *int   `json:"warehouse_id,omitempty"`
	CategoryID   *int   `json:"category_id,omitempty"`
	LotNumber    string `json:"lot_number,omitempty"`
	ShowExpiring bool   `json:"show_expiring"`
	DaysToExpiry int    `json:"days_to_expiry,omitempty"`
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
}

type InventorySummary struct {
	ProductID          int                 `json:"product_id"`
	ProductSKU         string              `json:"product_sku"`
	ProductName        string              `json:"product_name"`
	TotalOnHand        float64             `json:"total_on_hand"`
	TotalAllocated     float64             `json:"total_allocated"`
	TotalOnOrder       float64             `json:"total_on_order"`
	TotalAvailable     float64             `json:"total_available"`
	AverageCost        float64             `json:"average_cost"`
	InventoryValue     float64             `json:"inventory_value"`
	WarehouseBreakdown []WarehouseQuantity `json:"warehouse_breakdown"`
}

type WarehouseQuantity struct {
	WarehouseID   int     `json:"warehouse_id"`
	WarehouseName string  `json:"warehouse_name"`
	OnHand        float64 `json:"on_hand"`
	Allocated     float64 `json:"allocated"`
	Available     float64 `json:"available"`
}

// ============================================
// Lot Tracking
// ============================================

type LotInfo struct {
	LotNumber      string     `json:"lot_number"`
	ProductID      int        `json:"product_id"`
	ProductName    string     `json:"product_name"`
	ProductionDate CustomDate `json:"production_date"`
	ExpiryDate     CustomDate `json:"expiry_date"`
	Quantity       float64    `json:"quantity"`
	WarehouseID    int        `json:"warehouse_id"`
	WarehouseName  string     `json:"warehouse_name"`
	LocationCode   string     `json:"location_code"`
	Status         string     `json:"status"`
}

// ============================================
// Validation
// ============================================

func ValidateAdjustInventory(v *Validator, req *AdjustInventoryRequest) {
	v.Check(req.ProductID > 0, "product_id", "Product ID is required")
	v.Check(req.WarehouseID > 0, "warehouse_id", "Warehouse ID is required")
	v.Check(req.Quantity != 0, "quantity", "Quantity cannot be zero")
	v.Check(req.Reason != "", "reason", "Reason is required")
}

func ValidateTransferInventory(v *Validator, req *TransferInventoryRequest) {
	v.Check(req.ProductID > 0, "product_id", "Product ID is required")
	v.Check(req.FromWarehouseID > 0, "from_warehouse_id", "Source warehouse is required")
	v.Check(req.ToWarehouseID > 0, "to_warehouse_id", "Destination warehouse is required")
	v.Check(req.FromWarehouseID != req.ToWarehouseID, "to_warehouse_id", "Destination must be different from source")
	v.Check(req.Quantity > 0, "quantity", "Quantity must be positive")
}
