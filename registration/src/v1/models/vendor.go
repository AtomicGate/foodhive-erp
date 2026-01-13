package models

// ============================================
// Vendor Models
// ============================================

type Vendor struct {
	ID               int        `json:"id"`
	VendorCode       string     `json:"vendor_code"`
	Name             string     `json:"name"`
	AddressLine1     string     `json:"address_line1,omitempty"`
	AddressLine2     string     `json:"address_line2,omitempty"`
	City             string     `json:"city,omitempty"`
	State            string     `json:"state,omitempty"`
	PostalCode       string     `json:"postal_code,omitempty"`
	Country          string     `json:"country,omitempty"`
	Phone            string     `json:"phone,omitempty"`
	Email            string     `json:"email,omitempty"`
	PaymentTermsDays int        `json:"payment_terms_days"`
	Currency         string     `json:"currency"`
	LeadTimeDays     int        `json:"lead_time_days"`
	MinimumOrder     float64    `json:"minimum_order,omitempty"`
	BuyerID          *int       `json:"buyer_id,omitempty"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        CustomDate `json:"created_at"`
	UpdatedAt        CustomDate `json:"updated_at"`
}

type VendorProduct struct {
	ID                int     `json:"id"`
	VendorID          int     `json:"vendor_id"`
	ProductID         int     `json:"product_id"`
	VendorSKU         string  `json:"vendor_sku,omitempty"`
	VendorDescription string  `json:"vendor_description,omitempty"`
	UnitOfMeasure     string  `json:"unit_of_measure,omitempty"`
	UnitCost          float64 `json:"unit_cost"`
	MinimumOrderQty   float64 `json:"minimum_order_qty,omitempty"`
	LeadTimeDays      *int    `json:"lead_time_days,omitempty"`
	IsPreferred       bool    `json:"is_preferred"`
}

type VendorDiscount struct {
	ID              int     `json:"id"`
	VendorID        int     `json:"vendor_id"`
	DiscountDays    int     `json:"discount_days"`
	DiscountPercent float64 `json:"discount_percent"`
}

type VendorWithDetails struct {
	Vendor    Vendor           `json:"vendor"`
	BuyerName string           `json:"buyer_name,omitempty"`
	Products  []VendorProduct  `json:"products,omitempty"`
	Discounts []VendorDiscount `json:"discounts,omitempty"`
}

// ============================================
// Request/Response Types
// ============================================

type CreateVendorRequest struct {
	VendorCode       string  `json:"vendor_code"`
	Name             string  `json:"name"`
	AddressLine1     string  `json:"address_line1,omitempty"`
	AddressLine2     string  `json:"address_line2,omitempty"`
	City             string  `json:"city,omitempty"`
	State            string  `json:"state,omitempty"`
	PostalCode       string  `json:"postal_code,omitempty"`
	Country          string  `json:"country,omitempty"`
	Phone            string  `json:"phone,omitempty"`
	Email            string  `json:"email,omitempty"`
	PaymentTermsDays int     `json:"payment_terms_days"`
	Currency         string  `json:"currency"`
	LeadTimeDays     int     `json:"lead_time_days"`
	MinimumOrder     float64 `json:"minimum_order,omitempty"`
	BuyerID          *int    `json:"buyer_id,omitempty"`
}

type UpdateVendorRequest struct {
	Name             *string  `json:"name,omitempty"`
	AddressLine1     *string  `json:"address_line1,omitempty"`
	AddressLine2     *string  `json:"address_line2,omitempty"`
	City             *string  `json:"city,omitempty"`
	State            *string  `json:"state,omitempty"`
	PostalCode       *string  `json:"postal_code,omitempty"`
	Country          *string  `json:"country,omitempty"`
	Phone            *string  `json:"phone,omitempty"`
	Email            *string  `json:"email,omitempty"`
	PaymentTermsDays *int     `json:"payment_terms_days,omitempty"`
	Currency         *string  `json:"currency,omitempty"`
	LeadTimeDays     *int     `json:"lead_time_days,omitempty"`
	MinimumOrder     *float64 `json:"minimum_order,omitempty"`
	BuyerID          *int     `json:"buyer_id,omitempty"`
	IsActive         *bool    `json:"is_active,omitempty"`
}

type VendorListFilters struct {
	Search   string `json:"search,omitempty"`
	BuyerID  *int   `json:"buyer_id,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// ============================================
// Validation
// ============================================

func ValidateVendor(v *Validator, req *CreateVendorRequest) {
	v.Check(req.VendorCode != "", "vendor_code", "Vendor code is required")
	v.Check(len(req.VendorCode) <= 20, "vendor_code", "Vendor code must be 20 characters or less")
	v.Check(req.Name != "", "name", "Vendor name is required")
	if req.Currency != "" {
		v.Check(len(req.Currency) == 3, "currency", "Currency must be a 3-letter code")
	}
	if req.Email != "" {
		v.Check(len(req.Email) <= 255, "email", "Email must be 255 characters or less")
	}
	v.Check(req.PaymentTermsDays >= 0, "payment_terms_days", "Payment terms must be 0 or greater")
	v.Check(req.LeadTimeDays >= 0, "lead_time_days", "Lead time must be 0 or greater")
}

func ValidateVendorProduct(v *Validator, req *VendorProduct) {
	v.Check(req.VendorID > 0, "vendor_id", "Vendor ID is required")
	v.Check(req.ProductID > 0, "product_id", "Product ID is required")
	v.Check(req.UnitCost >= 0, "unit_cost", "Unit cost must be 0 or greater")
}

func ValidateVendorDiscount(v *Validator, req *VendorDiscount) {
	v.Check(req.VendorID > 0, "vendor_id", "Vendor ID is required")
	v.Check(req.DiscountDays > 0, "discount_days", "Discount days must be greater than 0")
	v.Check(req.DiscountPercent > 0 && req.DiscountPercent <= 100, "discount_percent", "Discount percent must be between 0 and 100")
}
