package models

// ============================================
// Customer Models
// ============================================

type Customer struct {
	ID                 int        `json:"id"`
	CustomerCode       string     `json:"customer_code"`
	Name               string     `json:"name"`
	BillingAddressID   *int       `json:"billing_address_id,omitempty"`
	CreditLimit        float64    `json:"credit_limit"`
	CurrentBalance     float64    `json:"current_balance"`
	PaymentTermsDays   int        `json:"payment_terms_days"`
	Currency           string     `json:"currency"`
	SalesRepID         *int       `json:"sales_rep_id,omitempty"`
	DefaultRouteID     *int       `json:"default_route_id,omitempty"`
	DefaultWarehouseID *int       `json:"default_warehouse_id,omitempty"`
	TaxExempt          bool       `json:"tax_exempt"`
	IsActive           bool       `json:"is_active"`
	CreatedBy          int        `json:"created_by"`
	CreatedAt          CustomDate `json:"created_at"`
	UpdatedAt          CustomDate `json:"updated_at"`
}

type CustomerShipTo struct {
	ID           int    `json:"id"`
	CustomerID   int    `json:"customer_id"`
	ShipToCode   string `json:"ship_to_code"`
	Name         string `json:"name"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone,omitempty"`
	IsDefault    bool   `json:"is_default"`
	WarehouseID  *int   `json:"warehouse_id,omitempty"`
	RouteID      *int   `json:"route_id,omitempty"`
}

type CustomerOrderGuide struct {
	ID                  int        `json:"id"`
	CustomerID          int        `json:"customer_id"`
	ProductID           int        `json:"product_id"`
	ProductName         string     `json:"product_name,omitempty"`
	ProductSKU          string     `json:"product_sku,omitempty"`
	DefaultQuantity     float64    `json:"default_quantity"`
	LastOrderedDate     CustomDate `json:"last_ordered_date,omitempty"`
	LastOrderedQuantity float64    `json:"last_ordered_quantity"`
	AvgWeeklyQuantity   float64    `json:"avg_weekly_quantity"`
	TimesOrdered        int        `json:"times_ordered"`
	IsPushItem          bool       `json:"is_push_item"`
	CustomPrice         *float64   `json:"custom_price,omitempty"`
}

type CustomerWithDetails struct {
	Customer      Customer         `json:"customer"`
	ShipToList    []CustomerShipTo `json:"ship_to_addresses"`
	SalesRepName  string           `json:"sales_rep_name,omitempty"`
	WarehouseName string           `json:"warehouse_name,omitempty"`
}

// ============================================
// Request/Response Types
// ============================================

type CreateCustomerRequest struct {
	CustomerCode       string  `json:"customer_code"`
	Name               string  `json:"name"`
	CreditLimit        float64 `json:"credit_limit"`
	PaymentTermsDays   int     `json:"payment_terms_days"`
	Currency           string  `json:"currency"`
	SalesRepID         *int    `json:"sales_rep_id,omitempty"`
	DefaultWarehouseID *int    `json:"default_warehouse_id,omitempty"`
	TaxExempt          bool    `json:"tax_exempt"`
}

type UpdateCustomerRequest struct {
	Name               *string  `json:"name,omitempty"`
	CreditLimit        *float64 `json:"credit_limit,omitempty"`
	PaymentTermsDays   *int     `json:"payment_terms_days,omitempty"`
	Currency           *string  `json:"currency,omitempty"`
	SalesRepID         *int     `json:"sales_rep_id,omitempty"`
	DefaultWarehouseID *int     `json:"default_warehouse_id,omitempty"`
	TaxExempt          *bool    `json:"tax_exempt,omitempty"`
	IsActive           *bool    `json:"is_active,omitempty"`
}

type CustomerListFilters struct {
	Search      string `json:"search,omitempty"`
	SalesRepID  *int   `json:"sales_rep_id,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
	WarehouseID *int   `json:"warehouse_id,omitempty"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
}

// ============================================
// Validation
// ============================================

func ValidateCustomer(v *Validator, c *CreateCustomerRequest) {
	v.Check(c.CustomerCode != "", "customer_code", "Customer code is required")
	v.Check(len(c.CustomerCode) <= 20, "customer_code", "Customer code must be 20 characters or less")
	v.Check(c.Name != "", "name", "Customer name is required")
	v.Check(c.CreditLimit >= 0, "credit_limit", "Credit limit cannot be negative")
	v.Check(c.PaymentTermsDays >= 0, "payment_terms_days", "Payment terms cannot be negative")
	if c.Currency != "" {
		v.Check(len(c.Currency) == 3, "currency", "Currency must be a 3-letter code")
	}
}
