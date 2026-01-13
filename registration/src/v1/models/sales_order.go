package models

// ============================================
// Sales Order Enums
// ============================================

type OrderType string

const (
	OrderTypeStandard   OrderType = "STANDARD"
	OrderTypeAdvance    OrderType = "ADVANCE"
	OrderTypePrePaid    OrderType = "PRE_PAID"
	OrderTypeOnHold     OrderType = "ON_HOLD"
	OrderTypeQuote      OrderType = "QUOTE"
	OrderTypeCreditMemo OrderType = "CREDIT_MEMO"
	OrderTypePickUp     OrderType = "PICK_UP"
)

type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "DRAFT"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusPicking   OrderStatus = "PICKING"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusInvoiced  OrderStatus = "INVOICED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// ============================================
// Sales Order Models
// ============================================

type SalesOrder struct {
	ID                int         `json:"id"`
	OrderNumber       string      `json:"order_number"`
	CustomerID        int         `json:"customer_id"`
	ShipToID          *int        `json:"ship_to_id,omitempty"`
	OrderType         OrderType   `json:"order_type"`
	OrderDate         CustomDate  `json:"order_date"`
	RequestedShipDate CustomDate  `json:"requested_ship_date,omitempty"`
	ActualShipDate    CustomDate  `json:"actual_ship_date,omitempty"`
	WarehouseID       int         `json:"warehouse_id"`
	RouteID           *int        `json:"route_id,omitempty"`
	Status            OrderStatus `json:"status"`
	Subtotal          float64     `json:"subtotal"`
	TaxAmount         float64     `json:"tax_amount"`
	FreightAmount     float64     `json:"freight_amount"`
	DiscountAmount    float64     `json:"discount_amount"`
	TotalAmount       float64     `json:"total_amount"`
	Notes             string      `json:"notes,omitempty"`
	PONumber          string      `json:"po_number,omitempty"`
	SalesRepID        *int        `json:"sales_rep_id,omitempty"`
	CreatedBy         int         `json:"created_by"`
	CreatedAt         CustomDate  `json:"created_at"`
	UpdatedAt         CustomDate  `json:"updated_at"`
}

type SalesOrderLine struct {
	ID              int        `json:"id"`
	OrderID         int        `json:"order_id"`
	LineNumber      int        `json:"line_number"`
	ProductID       int        `json:"product_id"`
	Description     string     `json:"description,omitempty"`
	QuantityOrdered float64    `json:"quantity_ordered"`
	QuantityShipped float64    `json:"quantity_shipped"`
	UnitOfMeasure   string     `json:"unit_of_measure"`
	UnitPrice       float64    `json:"unit_price"`
	DiscountPercent float64    `json:"discount_percent"`
	LineTotal       float64    `json:"line_total"`
	LotNumber       string     `json:"lot_number,omitempty"`
	ExpiryDate      CustomDate `json:"expiry_date,omitempty"`
	CatchWeight     float64    `json:"catch_weight,omitempty"`
	Cost            float64    `json:"cost"`
}

type SalesOrderWithDetails struct {
	Order         SalesOrder       `json:"order"`
	Lines         []SalesOrderLine `json:"lines"`
	CustomerName  string           `json:"customer_name"`
	CustomerCode  string           `json:"customer_code"`
	ShipToName    string           `json:"ship_to_name,omitempty"`
	ShipToAddress string           `json:"ship_to_address,omitempty"`
	WarehouseName string           `json:"warehouse_name"`
	SalesRepName  string           `json:"sales_rep_name,omitempty"`
	RouteName     string           `json:"route_name,omitempty"`
}

// ============================================
// Request/Response Types
// ============================================

type CreateSalesOrderRequest struct {
	CustomerID        int                           `json:"customer_id"`
	ShipToID          *int                          `json:"ship_to_id,omitempty"`
	OrderType         OrderType                     `json:"order_type"`
	RequestedShipDate string                        `json:"requested_ship_date,omitempty"`
	WarehouseID       int                           `json:"warehouse_id"`
	RouteID           *int                          `json:"route_id,omitempty"`
	Notes             string                        `json:"notes,omitempty"`
	PONumber          string                        `json:"po_number,omitempty"`
	Lines             []CreateSalesOrderLineRequest `json:"lines"`
}

type CreateSalesOrderLineRequest struct {
	ProductID       int     `json:"product_id"`
	Quantity        float64 `json:"quantity"`
	UnitOfMeasure   string  `json:"unit_of_measure"`
	UnitPrice       float64 `json:"unit_price,omitempty"` // Optional, use customer price if not provided
	DiscountPercent float64 `json:"discount_percent,omitempty"`
	LotNumber       string  `json:"lot_number,omitempty"`
	Notes           string  `json:"notes,omitempty"`
}

type UpdateSalesOrderRequest struct {
	ShipToID          *int         `json:"ship_to_id,omitempty"`
	RequestedShipDate *string      `json:"requested_ship_date,omitempty"`
	RouteID           *int         `json:"route_id,omitempty"`
	Notes             *string      `json:"notes,omitempty"`
	PONumber          *string      `json:"po_number,omitempty"`
	Status            *OrderStatus `json:"status,omitempty"`
}

type SalesOrderListFilters struct {
	CustomerID  *int         `json:"customer_id,omitempty"`
	Status      *OrderStatus `json:"status,omitempty"`
	OrderType   *OrderType   `json:"order_type,omitempty"`
	WarehouseID *int         `json:"warehouse_id,omitempty"`
	RouteID     *int         `json:"route_id,omitempty"`
	SalesRepID  *int         `json:"sales_rep_id,omitempty"`
	DateFrom    string       `json:"date_from,omitempty"`
	DateTo      string       `json:"date_to,omitempty"`
	Page        int          `json:"page"`
	PageSize    int          `json:"page_size"`
}

// ============================================
// Order Guide Entry
// ============================================

type OrderGuideEntry struct {
	ProductID       int     `json:"product_id"`
	ProductSKU      string  `json:"product_sku"`
	ProductName     string  `json:"product_name"`
	DefaultQuantity float64 `json:"default_quantity"`
	LastOrderedQty  float64 `json:"last_ordered_qty"`
	AvgWeeklyQty    float64 `json:"avg_weekly_qty"`
	OnHand          float64 `json:"on_hand"`
	Allocated       float64 `json:"allocated"`
	Available       float64 `json:"available"`
	UnitPrice       float64 `json:"unit_price"`
	UnitOfMeasure   string  `json:"unit_of_measure"`
	IsPushItem      bool    `json:"is_push_item"`
}

// ============================================
// Lost Sales
// ============================================

type LostSale struct {
	ID                int        `json:"id"`
	OrderID           int        `json:"order_id"`
	ProductID         int        `json:"product_id"`
	ProductName       string     `json:"product_name"`
	QuantityRequested float64    `json:"quantity_requested"`
	QuantityAvailable float64    `json:"quantity_available"`
	Reason            string     `json:"reason"`
	CreatedAt         CustomDate `json:"created_at"`
}

// ============================================
// Validation
// ============================================

func ValidateSalesOrder(v *Validator, req *CreateSalesOrderRequest) {
	v.Check(req.CustomerID > 0, "customer_id", "Customer is required")
	v.Check(req.WarehouseID > 0, "warehouse_id", "Warehouse is required")
	v.Check(len(req.Lines) > 0, "lines", "At least one line is required")

	for i, line := range req.Lines {
		v.Check(line.ProductID > 0, "lines", "Product ID is required for all lines")
		v.Check(line.Quantity > 0, "lines", "Quantity must be positive for all lines")
		v.Check(line.UnitOfMeasure != "", "lines", "Unit of measure is required for all lines")
		_ = i // Avoid unused variable
	}
}
