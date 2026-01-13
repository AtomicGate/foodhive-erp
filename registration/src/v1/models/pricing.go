package models

// ============================================
// Pricing Enums
// ============================================

type PriceLevel string

const (
	PriceLevelContract   PriceLevel = "CONTRACT"    // Level 1: Customer contract price
	PriceLevelCustomer   PriceLevel = "CUSTOMER"    // Level 2: Customer-specific price
	PriceLevelPromotion  PriceLevel = "PROMOTION"   // Level 3: Promotional price
	PriceLevelPriceLevel PriceLevel = "PRICE_LEVEL" // Level 4: Price level (A, B, C, etc.)
	PriceLevelBase       PriceLevel = "BASE"        // Level 5: Product base price
)

type CostingMethod string

const (
	CostMethodAverage  CostingMethod = "AVERAGE"  // Weighted average cost
	CostMethodLast     CostingMethod = "LAST"     // Last cost
	CostMethodLanded   CostingMethod = "LANDED"   // Landed cost (including freight, duties)
	CostMethodMarket   CostingMethod = "MARKET"   // Market cost
	CostMethodVendor   CostingMethod = "VENDOR"   // Default vendor cost
	CostMethodAdjusted CostingMethod = "ADJUSTED" // User adjusted cost
)

// ============================================
// Pricing Models
// ============================================

type ProductPrice struct {
	ID            int        `json:"id"`
	ProductID     int        `json:"product_id"`
	PriceLevel    PriceLevel `json:"price_level"`
	Price         float64    `json:"price"`
	EffectiveDate CustomDate `json:"effective_date"`
	ExpiryDate    CustomDate `json:"expiry_date,omitempty"`
	MinQuantity   float64    `json:"min_quantity,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     CustomDate `json:"created_at"`
	UpdatedAt     CustomDate `json:"updated_at"`
}

type CustomerPrice struct {
	ID            int        `json:"id"`
	CustomerID    int        `json:"customer_id"`
	ProductID     int        `json:"product_id"`
	Price         float64    `json:"price"`
	EffectiveDate CustomDate `json:"effective_date"`
	ExpiryDate    CustomDate `json:"expiry_date,omitempty"`
	Notes         string     `json:"notes,omitempty"`
	CreatedBy     int        `json:"created_by"`
	CreatedAt     CustomDate `json:"created_at"`
}

type ContractPrice struct {
	ID            int        `json:"id"`
	ContractCode  string     `json:"contract_code"`
	CustomerID    int        `json:"customer_id"`
	ProductID     int        `json:"product_id"`
	Price         float64    `json:"price"`
	EffectiveDate CustomDate `json:"effective_date"`
	ExpiryDate    CustomDate `json:"expiry_date"`
	MinQuantity   float64    `json:"min_quantity,omitempty"`
	MaxQuantity   float64    `json:"max_quantity,omitempty"`
	Notes         string     `json:"notes,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     CustomDate `json:"created_at"`
}

type PromotionalPrice struct {
	ID              int        `json:"id"`
	PromotionCode   string     `json:"promotion_code"`
	Name            string     `json:"name"`
	ProductID       *int       `json:"product_id,omitempty"`
	CategoryID      *int       `json:"category_id,omitempty"`
	DiscountPercent float64    `json:"discount_percent,omitempty"`
	DiscountAmount  float64    `json:"discount_amount,omitempty"`
	FixedPrice      float64    `json:"fixed_price,omitempty"`
	EffectiveDate   CustomDate `json:"effective_date"`
	ExpiryDate      CustomDate `json:"expiry_date"`
	CustomerGroupID *int       `json:"customer_group_id,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       CustomDate `json:"created_at"`
}

type ProductCost struct {
	ID             int           `json:"id"`
	ProductID      int           `json:"product_id"`
	CostingMethod  CostingMethod `json:"costing_method"`
	Cost           float64       `json:"cost"`
	EffectiveDate  CustomDate    `json:"effective_date"`
	FreightFactor  float64       `json:"freight_factor,omitempty"`
	DutyFactor     float64       `json:"duty_factor,omitempty"`
	HandlingFactor float64       `json:"handling_factor,omitempty"`
	LandedCost     float64       `json:"landed_cost,omitempty"`
	Notes          string        `json:"notes,omitempty"`
	UpdatedBy      int           `json:"updated_by"`
	UpdatedAt      CustomDate    `json:"updated_at"`
}

// ============================================
// Price Lookup Result
// ============================================

type PriceLookupResult struct {
	ProductID     int        `json:"product_id"`
	ProductSKU    string     `json:"product_sku"`
	ProductName   string     `json:"product_name"`
	PriceLevel    PriceLevel `json:"price_level"`
	Price         float64    `json:"price"`
	OriginalPrice float64    `json:"original_price,omitempty"`
	DiscountPct   float64    `json:"discount_percent,omitempty"`
	Cost          float64    `json:"cost"`
	Margin        float64    `json:"margin"`
	MarginPercent float64    `json:"margin_percent"`
	IsBelowCost   bool       `json:"is_below_cost"`
	PriceSource   string     `json:"price_source"`
}

type PriceListItem struct {
	ProductID     int     `json:"product_id"`
	ProductSKU    string  `json:"product_sku"`
	ProductName   string  `json:"product_name"`
	CategoryName  string  `json:"category_name"`
	UnitOfMeasure string  `json:"unit_of_measure"`
	BasePrice     float64 `json:"base_price"`
	LevelAPrice   float64 `json:"level_a_price,omitempty"`
	LevelBPrice   float64 `json:"level_b_price,omitempty"`
	LevelCPrice   float64 `json:"level_c_price,omitempty"`
	Cost          float64 `json:"cost"`
}

// ============================================
// Request/Response Types
// ============================================

type SetProductPriceRequest struct {
	ProductID     int        `json:"product_id"`
	PriceLevel    PriceLevel `json:"price_level"`
	Price         float64    `json:"price"`
	EffectiveDate string     `json:"effective_date"`
	ExpiryDate    string     `json:"expiry_date,omitempty"`
	MinQuantity   float64    `json:"min_quantity,omitempty"`
}

type SetCustomerPriceRequest struct {
	CustomerID    int     `json:"customer_id"`
	ProductID     int     `json:"product_id"`
	Price         float64 `json:"price"`
	EffectiveDate string  `json:"effective_date"`
	ExpiryDate    string  `json:"expiry_date,omitempty"`
	Notes         string  `json:"notes,omitempty"`
}

type CreatePriceContractRequest struct {
	ContractCode  string  `json:"contract_code"`
	CustomerID    int     `json:"customer_id"`
	ProductID     int     `json:"product_id"`
	Price         float64 `json:"price"`
	EffectiveDate string  `json:"effective_date"`
	ExpiryDate    string  `json:"expiry_date"`
	MinQuantity   float64 `json:"min_quantity,omitempty"`
	MaxQuantity   float64 `json:"max_quantity,omitempty"`
	Notes         string  `json:"notes,omitempty"`
}

type CreatePromotionRequest struct {
	PromotionCode   string  `json:"promotion_code"`
	Name            string  `json:"name"`
	ProductID       *int    `json:"product_id,omitempty"`
	CategoryID      *int    `json:"category_id,omitempty"`
	DiscountPercent float64 `json:"discount_percent,omitempty"`
	DiscountAmount  float64 `json:"discount_amount,omitempty"`
	FixedPrice      float64 `json:"fixed_price,omitempty"`
	EffectiveDate   string  `json:"effective_date"`
	ExpiryDate      string  `json:"expiry_date"`
	CustomerGroupID *int    `json:"customer_group_id,omitempty"`
}

type UpdateProductCostRequest struct {
	ProductID      int           `json:"product_id"`
	CostingMethod  CostingMethod `json:"costing_method"`
	Cost           float64       `json:"cost"`
	FreightFactor  float64       `json:"freight_factor,omitempty"`
	DutyFactor     float64       `json:"duty_factor,omitempty"`
	HandlingFactor float64       `json:"handling_factor,omitempty"`
	Notes          string        `json:"notes,omitempty"`
}

type MassPriceUpdateRequest struct {
	ProductIDs      []int      `json:"product_ids,omitempty"`
	CategoryID      *int       `json:"category_id,omitempty"`
	PriceLevel      PriceLevel `json:"price_level"`
	AdjustmentType  string     `json:"adjustment_type"` // "PERCENT" or "AMOUNT"
	AdjustmentValue float64    `json:"adjustment_value"`
	EffectiveDate   string     `json:"effective_date"`
}

type PriceLookupRequest struct {
	ProductID  int     `json:"product_id"`
	CustomerID *int    `json:"customer_id,omitempty"`
	Quantity   float64 `json:"quantity"`
	AsOfDate   string  `json:"as_of_date,omitempty"`
}

type PriceListFilters struct {
	CategoryID    *int   `json:"category_id,omitempty"`
	EffectiveDate string `json:"effective_date,omitempty"`
	IncludeCost   bool   `json:"include_cost"`
	Page          int    `json:"page"`
	PageSize      int    `json:"page_size"`
}

// ============================================
// Validation
// ============================================

func ValidateProductPrice(v *Validator, req *SetProductPriceRequest) {
	v.Check(req.ProductID > 0, "product_id", "Product is required")
	v.Check(req.PriceLevel != "", "price_level", "Price level is required")
	v.Check(req.Price >= 0, "price", "Price must be non-negative")
	v.Check(req.EffectiveDate != "", "effective_date", "Effective date is required")
}

func ValidateCustomerPrice(v *Validator, req *SetCustomerPriceRequest) {
	v.Check(req.CustomerID > 0, "customer_id", "Customer is required")
	v.Check(req.ProductID > 0, "product_id", "Product is required")
	v.Check(req.Price >= 0, "price", "Price must be non-negative")
	v.Check(req.EffectiveDate != "", "effective_date", "Effective date is required")
}

func ValidatePriceContract(v *Validator, req *CreatePriceContractRequest) {
	v.Check(req.ContractCode != "", "contract_code", "Contract code is required")
	v.Check(req.CustomerID > 0, "customer_id", "Customer is required")
	v.Check(req.ProductID > 0, "product_id", "Product is required")
	v.Check(req.Price >= 0, "price", "Price must be non-negative")
	v.Check(req.EffectiveDate != "", "effective_date", "Effective date is required")
	v.Check(req.ExpiryDate != "", "expiry_date", "Expiry date is required")
}

func ValidatePromotion(v *Validator, req *CreatePromotionRequest) {
	v.Check(req.PromotionCode != "", "promotion_code", "Promotion code is required")
	v.Check(req.Name != "", "name", "Name is required")
	v.Check(req.ProductID != nil || req.CategoryID != nil, "product_id", "Product or category is required")
	v.Check(req.DiscountPercent > 0 || req.DiscountAmount > 0 || req.FixedPrice > 0,
		"discount", "Discount percent, amount, or fixed price is required")
	v.Check(req.EffectiveDate != "", "effective_date", "Effective date is required")
	v.Check(req.ExpiryDate != "", "expiry_date", "Expiry date is required")
}

func ValidateMassPriceUpdate(v *Validator, req *MassPriceUpdateRequest) {
	v.Check(len(req.ProductIDs) > 0 || req.CategoryID != nil, "products", "Products or category is required")
	v.Check(req.PriceLevel != "", "price_level", "Price level is required")
	v.Check(req.AdjustmentType == "PERCENT" || req.AdjustmentType == "AMOUNT",
		"adjustment_type", "Must be PERCENT or AMOUNT")
	v.Check(req.EffectiveDate != "", "effective_date", "Effective date is required")
}
