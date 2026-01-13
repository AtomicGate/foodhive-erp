package models

// ============================================
// Purchase Order Enums
// ============================================

type POStatus string

const (
	POStatusDraft     POStatus = "DRAFT"
	POStatusSubmitted POStatus = "SUBMITTED"
	POStatusConfirmed POStatus = "CONFIRMED"
	POStatusPartial   POStatus = "PARTIAL"
	POStatusReceived  POStatus = "RECEIVED"
	POStatusClosed    POStatus = "CLOSED"
	POStatusCancelled POStatus = "CANCELLED"
)

// ============================================
// Purchase Order Models
// ============================================

type PurchaseOrder struct {
	ID            int        `json:"id"`
	PONumber      string     `json:"po_number"`
	VendorID      int        `json:"vendor_id"`
	WarehouseID   int        `json:"warehouse_id"`
	OrderDate     CustomDate `json:"order_date"`
	ExpectedDate  CustomDate `json:"expected_date,omitempty"`
	ReceivedDate  CustomDate `json:"received_date,omitempty"`
	Status        POStatus   `json:"status"`
	Subtotal      float64    `json:"subtotal"`
	TaxAmount     float64    `json:"tax_amount"`
	FreightAmount float64    `json:"freight_amount"`
	TotalAmount   float64    `json:"total_amount"`
	Notes         string     `json:"notes,omitempty"`
	BuyerID       *int       `json:"buyer_id,omitempty"`
	CreatedBy     int        `json:"created_by"`
	CreatedAt     CustomDate `json:"created_at"`
	UpdatedAt     CustomDate `json:"updated_at"`
}

type PurchaseOrderLine struct {
	ID               int        `json:"id"`
	POID             int        `json:"po_id"`
	LineNumber       int        `json:"line_number"`
	ProductID        int        `json:"product_id"`
	Description      string     `json:"description,omitempty"`
	QuantityOrdered  float64    `json:"quantity_ordered"`
	QuantityReceived float64    `json:"quantity_received"`
	UnitOfMeasure    string     `json:"unit_of_measure"`
	UnitCost         float64    `json:"unit_cost"`
	LineTotal        float64    `json:"line_total"`
	ExpectedDate     CustomDate `json:"expected_date,omitempty"`
}

type PurchaseOrderWithDetails struct {
	Order         PurchaseOrder       `json:"order"`
	Lines         []PurchaseOrderLine `json:"lines"`
	VendorName    string              `json:"vendor_name"`
	VendorCode    string              `json:"vendor_code"`
	WarehouseName string              `json:"warehouse_name"`
	BuyerName     string              `json:"buyer_name,omitempty"`
}

// ============================================
// Receiving Models
// ============================================

type Receiving struct {
	ID              int        `json:"id"`
	ReceivingNumber string     `json:"receiving_number"`
	POID            *int       `json:"po_id,omitempty"`
	WarehouseID     int        `json:"warehouse_id"`
	VendorID        int        `json:"vendor_id"`
	ReceivedDate    CustomDate `json:"received_date"`
	Notes           string     `json:"notes,omitempty"`
	ReceivedBy      int        `json:"received_by"`
	CreatedAt       CustomDate `json:"created_at"`
}

type ReceivingLine struct {
	ID               int        `json:"id"`
	ReceivingID      int        `json:"receiving_id"`
	POLineID         *int       `json:"po_line_id,omitempty"`
	ProductID        int        `json:"product_id"`
	QuantityReceived float64    `json:"quantity_received"`
	UnitOfMeasure    string     `json:"unit_of_measure"`
	LotNumber        string     `json:"lot_number,omitempty"`
	ProductionDate   CustomDate `json:"production_date,omitempty"`
	ExpiryDate       CustomDate `json:"expiry_date,omitempty"`
	LocationCode     string     `json:"location_code,omitempty"`
	UnitCost         float64    `json:"unit_cost"`
	IsShortShipment  bool       `json:"is_short_shipment"`
	Notes            string     `json:"notes,omitempty"`
}

type ReceivingWithDetails struct {
	Receiving     Receiving       `json:"receiving"`
	Lines         []ReceivingLine `json:"lines"`
	VendorName    string          `json:"vendor_name"`
	WarehouseName string          `json:"warehouse_name"`
	PONumber      string          `json:"po_number,omitempty"`
	ReceiverName  string          `json:"receiver_name"`
}

// ============================================
// Request/Response Types
// ============================================

type CreatePurchaseOrderRequest struct {
	VendorID     int                              `json:"vendor_id"`
	WarehouseID  int                              `json:"warehouse_id"`
	ExpectedDate string                           `json:"expected_date,omitempty"`
	Notes        string                           `json:"notes,omitempty"`
	BuyerID      *int                             `json:"buyer_id,omitempty"`
	Lines        []CreatePurchaseOrderLineRequest `json:"lines"`
}

type CreatePurchaseOrderLineRequest struct {
	ProductID     int     `json:"product_id"`
	Quantity      float64 `json:"quantity"`
	UnitOfMeasure string  `json:"unit_of_measure"`
	UnitCost      float64 `json:"unit_cost"`
	Description   string  `json:"description,omitempty"`
	ExpectedDate  string  `json:"expected_date,omitempty"`
}

type UpdatePurchaseOrderRequest struct {
	ExpectedDate *string   `json:"expected_date,omitempty"`
	Notes        *string   `json:"notes,omitempty"`
	BuyerID      *int      `json:"buyer_id,omitempty"`
	Status       *POStatus `json:"status,omitempty"`
}

type PurchaseOrderListFilters struct {
	VendorID    *int      `json:"vendor_id,omitempty"`
	WarehouseID *int      `json:"warehouse_id,omitempty"`
	Status      *POStatus `json:"status,omitempty"`
	BuyerID     *int      `json:"buyer_id,omitempty"`
	DateFrom    string    `json:"date_from,omitempty"`
	DateTo      string    `json:"date_to,omitempty"`
	Page        int       `json:"page"`
	PageSize    int       `json:"page_size"`
}

// ============================================
// Receiving Request Types
// ============================================

type CreateReceivingRequest struct {
	POID        *int                         `json:"po_id,omitempty"`
	VendorID    int                          `json:"vendor_id"`
	WarehouseID int                          `json:"warehouse_id"`
	Notes       string                       `json:"notes,omitempty"`
	Lines       []CreateReceivingLineRequest `json:"lines"`
}

type CreateReceivingLineRequest struct {
	POLineID       *int    `json:"po_line_id,omitempty"`
	ProductID      int     `json:"product_id"`
	Quantity       float64 `json:"quantity"`
	UnitOfMeasure  string  `json:"unit_of_measure"`
	LotNumber      string  `json:"lot_number,omitempty"`
	ProductionDate string  `json:"production_date,omitempty"`
	ExpiryDate     string  `json:"expiry_date,omitempty"`
	LocationCode   string  `json:"location_code,omitempty"`
	UnitCost       float64 `json:"unit_cost"`
	Notes          string  `json:"notes,omitempty"`
}

// ============================================
// Validation
// ============================================

func ValidatePurchaseOrder(v *Validator, req *CreatePurchaseOrderRequest) {
	v.Check(req.VendorID > 0, "vendor_id", "Vendor is required")
	v.Check(req.WarehouseID > 0, "warehouse_id", "Warehouse is required")
	v.Check(len(req.Lines) > 0, "lines", "At least one line is required")

	for i, line := range req.Lines {
		v.Check(line.ProductID > 0, "lines", "Product ID is required for all lines")
		v.Check(line.Quantity > 0, "lines", "Quantity must be positive for all lines")
		v.Check(line.UnitOfMeasure != "", "lines", "Unit of measure is required for all lines")
		v.Check(line.UnitCost >= 0, "lines", "Unit cost must be non-negative")
		_ = i
	}
}

func ValidateReceiving(v *Validator, req *CreateReceivingRequest) {
	v.Check(req.VendorID > 0, "vendor_id", "Vendor is required")
	v.Check(req.WarehouseID > 0, "warehouse_id", "Warehouse is required")
	v.Check(len(req.Lines) > 0, "lines", "At least one line is required")

	for i, line := range req.Lines {
		v.Check(line.ProductID > 0, "lines", "Product ID is required for all lines")
		v.Check(line.Quantity > 0, "lines", "Quantity must be positive for all lines")
		v.Check(line.UnitOfMeasure != "", "lines", "Unit of measure is required for all lines")
		_ = i
	}
}
