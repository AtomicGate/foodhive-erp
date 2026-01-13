package models

// ============================================
// Picking & Routing Enums
// ============================================

type PickListStatus string

const (
	PickListStatusPending    PickListStatus = "PENDING"
	PickListStatusInProgress PickListStatus = "IN_PROGRESS"
	PickListStatusComplete   PickListStatus = "COMPLETE"
	PickListStatusCancelled  PickListStatus = "CANCELLED"
)

// ============================================
// Route Models
// ============================================

type Route struct {
	ID            int            `json:"id"`
	RouteCode     string         `json:"route_code"`
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	WarehouseID   *int           `json:"warehouse_id,omitempty"`
	DriverID      *int           `json:"driver_id,omitempty"`
	VehicleID     *int           `json:"vehicle_id,omitempty"`
	DepartureTime string         `json:"departure_time,omitempty"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     CustomDateTime `json:"created_at"`
}

type RouteWithDetails struct {
	Route         Route       `json:"route"`
	WarehouseName string      `json:"warehouse_name,omitempty"`
	DriverName    string      `json:"driver_name,omitempty"`
	StopCount     int         `json:"stop_count"`
	Stops         []RouteStop `json:"stops,omitempty"`
}

type RouteStop struct {
	ID               int    `json:"id"`
	RouteID          int    `json:"route_id"`
	CustomerID       int    `json:"customer_id"`
	CustomerName     string `json:"customer_name,omitempty"`
	ShipToID         *int   `json:"ship_to_id,omitempty"`
	ShipToName       string `json:"ship_to_name,omitempty"`
	ShipToAddress    string `json:"ship_to_address,omitempty"`
	StopSequence     int    `json:"stop_sequence"`
	EstimatedArrival string `json:"estimated_arrival,omitempty"`
	Notes            string `json:"notes,omitempty"`
}

// ============================================
// Pick List Models
// ============================================

type PickList struct {
	ID          int            `json:"id"`
	PickNumber  string         `json:"pick_number"`
	WarehouseID int            `json:"warehouse_id"`
	RouteID     *int           `json:"route_id,omitempty"`
	PickDate    CustomDate     `json:"pick_date"`
	Status      PickListStatus `json:"status"`
	PickerID    *int           `json:"picker_id,omitempty"`
	StartedAt   CustomDateTime `json:"started_at,omitempty"`
	CompletedAt CustomDateTime `json:"completed_at,omitempty"`
	CreatedBy   int            `json:"created_by"`
	CreatedAt   CustomDateTime `json:"created_at"`
}

type PickListLine struct {
	ID              int     `json:"id"`
	PickListID      int     `json:"pick_list_id"`
	OrderID         int     `json:"order_id"`
	OrderLineID     int     `json:"order_line_id"`
	ProductID       int     `json:"product_id"`
	LocationCode    string  `json:"location_code,omitempty"`
	LotNumber       string  `json:"lot_number,omitempty"`
	QuantityOrdered float64 `json:"quantity_ordered"`
	QuantityPicked  float64 `json:"quantity_picked"`
	CatchWeight     float64 `json:"catch_weight,omitempty"`
	PickedAt        string  `json:"picked_at,omitempty"`
	PickedBy        *int    `json:"picked_by,omitempty"`
	Notes           string  `json:"notes,omitempty"`
}

type PickListWithDetails struct {
	PickList      PickList       `json:"pick_list"`
	Lines         []PickListLine `json:"lines"`
	WarehouseName string         `json:"warehouse_name"`
	RouteName     string         `json:"route_name,omitempty"`
	PickerName    string         `json:"picker_name,omitempty"`
	TotalOrders   int            `json:"total_orders"`
	TotalLines    int            `json:"total_lines"`
	TotalPicked   int            `json:"total_picked"`
}

type PickListLineWithProduct struct {
	PickListLine  PickListLine `json:"line"`
	OrderNumber   string       `json:"order_number"`
	CustomerName  string       `json:"customer_name"`
	ProductSKU    string       `json:"product_sku"`
	ProductName   string       `json:"product_name"`
	UnitOfMeasure string       `json:"unit_of_measure"`
	ExpiryDate    CustomDate   `json:"expiry_date,omitempty"`
}

// ============================================
// Request/Response Types
// ============================================

type CreateRouteRequest struct {
	RouteCode     string `json:"route_code"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	WarehouseID   *int   `json:"warehouse_id,omitempty"`
	DriverID      *int   `json:"driver_id,omitempty"`
	DepartureTime string `json:"departure_time,omitempty"`
}

type UpdateRouteRequest struct {
	Name          *string `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	WarehouseID   *int    `json:"warehouse_id,omitempty"`
	DriverID      *int    `json:"driver_id,omitempty"`
	DepartureTime *string `json:"departure_time,omitempty"`
	IsActive      *bool   `json:"is_active,omitempty"`
}

type AddRouteStopRequest struct {
	CustomerID       int    `json:"customer_id"`
	ShipToID         *int   `json:"ship_to_id,omitempty"`
	StopSequence     int    `json:"stop_sequence"`
	EstimatedArrival string `json:"estimated_arrival,omitempty"`
	Notes            string `json:"notes,omitempty"`
}

type CreatePickListRequest struct {
	WarehouseID int    `json:"warehouse_id"`
	RouteID     *int   `json:"route_id,omitempty"`
	PickDate    string `json:"pick_date"`
	OrderIDs    []int  `json:"order_ids"`
}

type GeneratePickListRequest struct {
	WarehouseID int    `json:"warehouse_id"`
	RouteID     *int   `json:"route_id,omitempty"`
	PickDate    string `json:"pick_date"`
}

type ConfirmPickLineRequest struct {
	QuantityPicked float64 `json:"quantity_picked"`
	CatchWeight    float64 `json:"catch_weight,omitempty"`
	LotNumber      string  `json:"lot_number,omitempty"`
	LocationCode   string  `json:"location_code,omitempty"`
	Notes          string  `json:"notes,omitempty"`
}

type PickListFilters struct {
	WarehouseID *int            `json:"warehouse_id,omitempty"`
	RouteID     *int            `json:"route_id,omitempty"`
	Status      *PickListStatus `json:"status,omitempty"`
	PickerID    *int            `json:"picker_id,omitempty"`
	DateFrom    string          `json:"date_from,omitempty"`
	DateTo      string          `json:"date_to,omitempty"`
	Page        int             `json:"page"`
	PageSize    int             `json:"page_size"`
}

// ============================================
// Master Pick Report
// ============================================

type MasterPickItem struct {
	ProductID     int     `json:"product_id"`
	ProductSKU    string  `json:"product_sku"`
	ProductName   string  `json:"product_name"`
	CategoryName  string  `json:"category_name"`
	LocationCode  string  `json:"location_code"`
	TotalQuantity float64 `json:"total_quantity"`
	UnitOfMeasure string  `json:"unit_of_measure"`
	OrderCount    int     `json:"order_count"`
	LotNumber     string  `json:"lot_number,omitempty"`
	ExpiryDate    string  `json:"expiry_date,omitempty"`
}

// ============================================
// Validation
// ============================================

func ValidateRoute(v *Validator, req *CreateRouteRequest) {
	v.Check(req.RouteCode != "", "route_code", "Route code is required")
	v.Check(len(req.RouteCode) <= 10, "route_code", "Route code must be 10 characters or less")
	v.Check(req.Name != "", "name", "Name is required")
}

func ValidatePickList(v *Validator, req *CreatePickListRequest) {
	v.Check(req.WarehouseID > 0, "warehouse_id", "Warehouse is required")
	v.Check(req.PickDate != "", "pick_date", "Pick date is required")
	v.Check(len(req.OrderIDs) > 0, "order_ids", "At least one order is required")
}

func ValidateRouteStop(v *Validator, req *AddRouteStopRequest) {
	v.Check(req.CustomerID > 0, "customer_id", "Customer is required")
	v.Check(req.StopSequence > 0, "stop_sequence", "Stop sequence must be positive")
}
