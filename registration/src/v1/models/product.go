package models

// ============================================
// Product Models
// ============================================

type Product struct {
	ID               int        `json:"id"`
	SKU              string     `json:"sku"`
	Barcode          string     `json:"barcode,omitempty"`
	UPC              string     `json:"upc,omitempty"`
	Name             string     `json:"name"`
	Description      string     `json:"description,omitempty"`
	CategoryID       *int       `json:"category_id,omitempty"`
	BaseUnit         string     `json:"base_unit"`
	IsCatchWeight    bool       `json:"is_catch_weight"`
	CatchWeightUnit  string     `json:"catch_weight_unit,omitempty"`
	CountryOfOrigin  string     `json:"country_of_origin,omitempty"`
	ShelfLifeDays    int        `json:"shelf_life_days,omitempty"`
	MinShelfLifeDays int        `json:"min_shelf_life_days,omitempty"`
	IsLotTracked     bool       `json:"is_lot_tracked"`
	IsSerialized     bool       `json:"is_serialized"`
	HACCPCategory    string     `json:"haccp_category,omitempty"`
	QCRequired       bool       `json:"qc_required"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        CustomDate `json:"created_at"`
	UpdatedAt        CustomDate `json:"updated_at"`
}

type ProductCategory struct {
	ID                   int    `json:"id"`
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	ParentID             *int   `json:"parent_id,omitempty"`
	GLSalesAccountID     *int   `json:"gl_sales_account_id,omitempty"`
	GLCOGSAccountID      *int   `json:"gl_cogs_account_id,omitempty"`
	GLInventoryAccountID *int   `json:"gl_inventory_account_id,omitempty"`
}

type ProductUnit struct {
	ID               int     `json:"id"`
	ProductID        int     `json:"product_id"`
	UnitName         string  `json:"unit_name"`
	Description      string  `json:"description,omitempty"`
	ConversionFactor float64 `json:"conversion_factor"`
	Barcode          string  `json:"barcode,omitempty"`
	Weight           float64 `json:"weight,omitempty"`
	IsPurchaseUnit   bool    `json:"is_purchase_unit"`
	IsSalesUnit      bool    `json:"is_sales_unit"`
}

type ProductWithDetails struct {
	Product      Product          `json:"product"`
	Category     *ProductCategory `json:"category,omitempty"`
	Units        []ProductUnit    `json:"units"`
	CurrentStock float64          `json:"current_stock"`
	AvgCost      float64          `json:"avg_cost"`
}

// ============================================
// Request/Response Types
// ============================================

type CreateProductRequest struct {
	SKU              string `json:"sku"`
	Barcode          string `json:"barcode,omitempty"`
	UPC              string `json:"upc,omitempty"`
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	CategoryID       *int   `json:"category_id,omitempty"`
	BaseUnit         string `json:"base_unit"`
	IsCatchWeight    bool   `json:"is_catch_weight"`
	CatchWeightUnit  string `json:"catch_weight_unit,omitempty"`
	CountryOfOrigin  string `json:"country_of_origin,omitempty"`
	ShelfLifeDays    int    `json:"shelf_life_days,omitempty"`
	MinShelfLifeDays int    `json:"min_shelf_life_days,omitempty"`
	IsLotTracked     bool   `json:"is_lot_tracked"`
	IsSerialized     bool   `json:"is_serialized"`
	HACCPCategory    string `json:"haccp_category,omitempty"`
	QCRequired       bool   `json:"qc_required"`
}

type UpdateProductRequest struct {
	Name             *string `json:"name,omitempty"`
	Description      *string `json:"description,omitempty"`
	CategoryID       *int    `json:"category_id,omitempty"`
	BaseUnit         *string `json:"base_unit,omitempty"`
	CountryOfOrigin  *string `json:"country_of_origin,omitempty"`
	ShelfLifeDays    *int    `json:"shelf_life_days,omitempty"`
	MinShelfLifeDays *int    `json:"min_shelf_life_days,omitempty"`
	HACCPCategory    *string `json:"haccp_category,omitempty"`
	QCRequired       *bool   `json:"qc_required,omitempty"`
	IsActive         *bool   `json:"is_active,omitempty"`
}

type ProductListFilters struct {
	Search     string `json:"search,omitempty"`
	CategoryID *int   `json:"category_id,omitempty"`
	IsActive   *bool  `json:"is_active,omitempty"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

// ============================================
// Validation
// ============================================

func ValidateProduct(v *Validator, p *CreateProductRequest) {
	v.Check(p.SKU != "", "sku", "SKU is required")
	v.Check(len(p.SKU) <= 50, "sku", "SKU must be 50 characters or less")
	v.Check(p.Name != "", "name", "Product name is required")
	v.Check(p.BaseUnit != "", "base_unit", "Base unit is required")
	if p.IsCatchWeight {
		v.Check(p.CatchWeightUnit != "", "catch_weight_unit", "Catch weight unit is required for catch weight items")
	}
	if p.CountryOfOrigin != "" {
		v.Check(len(p.CountryOfOrigin) == 3, "country_of_origin", "Country of origin must be a 3-letter code")
	}
}
