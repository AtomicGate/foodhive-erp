package pricing

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

type PricingService interface {
	// Price Lookup (5-level hierarchy)
	GetPrice(ctx context.Context, req *models.PriceLookupRequest) (*models.PriceLookupResult, error)
	GetPricesForProducts(ctx context.Context, productIDs []int, customerID *int) ([]models.PriceLookupResult, error)

	// Product Prices (Base & Level Prices)
	SetProductPrice(ctx context.Context, req *models.SetProductPriceRequest, createdBy int) (int, error)
	GetProductPrices(ctx context.Context, productID int) ([]models.ProductPrice, error)
	DeleteProductPrice(ctx context.Context, id int) error

	// Customer Prices
	SetCustomerPrice(ctx context.Context, req *models.SetCustomerPriceRequest, createdBy int) (int, error)
	GetCustomerPrices(ctx context.Context, customerID int) ([]models.CustomerPrice, error)
	DeleteCustomerPrice(ctx context.Context, id int) error

	// Contract Prices
	CreateContract(ctx context.Context, req *models.CreatePriceContractRequest, createdBy int) (int, error)
	GetContract(ctx context.Context, id int) (*models.ContractPrice, error)
	ListContracts(ctx context.Context, customerID *int, activeOnly bool) ([]models.ContractPrice, error)
	DeactivateContract(ctx context.Context, id int) error

	// Promotional Prices
	CreatePromotion(ctx context.Context, req *models.CreatePromotionRequest, createdBy int) (int, error)
	GetPromotion(ctx context.Context, id int) (*models.PromotionalPrice, error)
	ListPromotions(ctx context.Context, activeOnly bool) ([]models.PromotionalPrice, error)
	DeactivatePromotion(ctx context.Context, id int) error

	// Product Costs
	UpdateProductCost(ctx context.Context, req *models.UpdateProductCostRequest, updatedBy int) error
	GetProductCost(ctx context.Context, productID int) (*models.ProductCost, error)

	// Mass Price Maintenance
	MassPriceUpdate(ctx context.Context, req *models.MassPriceUpdateRequest, createdBy int) (int, error)

	// Price Lists
	GetPriceList(ctx context.Context, filters *models.PriceListFilters) ([]models.PriceListItem, int64, error)

	// Margin Check
	CheckBelowCost(ctx context.Context, productID int, price float64) (bool, float64, error)
}

// ============================================
// Service Implementation
// ============================================

type pricingServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) PricingService {
	return &pricingServiceImpl{db: db}
}

// ============================================
// Price Lookup (5-Level Hierarchy)
// ============================================

func (s *pricingServiceImpl) GetPrice(ctx context.Context, req *models.PriceLookupRequest) (*models.PriceLookupResult, error) {
	asOfDate := time.Now()
	if req.AsOfDate != "" {
		t, _ := time.Parse("2006-01-02", req.AsOfDate)
		if !t.IsZero() {
			asOfDate = t
		}
	}

	result := &models.PriceLookupResult{
		ProductID: req.ProductID,
	}

	// Get product info
	s.db.QueryRow(ctx, `SELECT sku, name FROM products WHERE id = $1`, req.ProductID).Scan(
		&result.ProductSKU, &result.ProductName,
	)

	// Get cost
	s.db.QueryRow(ctx, `SELECT COALESCE(average_cost, 0) FROM inventory WHERE product_id = $1 LIMIT 1`,
		req.ProductID).Scan(&result.Cost)

	var price float64
	var priceFound bool

	// Level 1: Contract Price (if customer provided)
	if req.CustomerID != nil && !priceFound {
		err := s.db.QueryRow(ctx, `
			SELECT price FROM contract_prices
			WHERE customer_id = $1 AND product_id = $2 AND is_active = true
			AND effective_date <= $3 AND expiry_date >= $3
			AND (min_quantity IS NULL OR min_quantity <= $4)
			ORDER BY min_quantity DESC NULLS LAST LIMIT 1`,
			*req.CustomerID, req.ProductID, asOfDate, req.Quantity,
		).Scan(&price)
		if err == nil {
			result.Price = price
			result.PriceLevel = models.PriceLevelContract
			result.PriceSource = "Contract"
			priceFound = true
		}
	}

	// Level 2: Customer-Specific Price
	if req.CustomerID != nil && !priceFound {
		err := s.db.QueryRow(ctx, `
			SELECT price FROM customer_pricing
			WHERE customer_id = $1 AND product_id = $2
			AND effective_date <= $3 AND (expiry_date IS NULL OR expiry_date >= $3)
			ORDER BY effective_date DESC LIMIT 1`,
			*req.CustomerID, req.ProductID, asOfDate,
		).Scan(&price)
		if err == nil {
			result.Price = price
			result.PriceLevel = models.PriceLevelCustomer
			result.PriceSource = "Customer Price"
			priceFound = true
		}
	}

	// Level 3: Promotional Price
	if !priceFound {
		var discountPct, discountAmt, fixedPrice float64
		err := s.db.QueryRow(ctx, `
			SELECT discount_percent, discount_amount, fixed_price
			FROM promotional_prices
			WHERE (product_id = $1 OR category_id = (SELECT category_id FROM products WHERE id = $1))
			AND is_active = true AND effective_date <= $2 AND expiry_date >= $2
			ORDER BY fixed_price DESC NULLS LAST, discount_percent DESC LIMIT 1`,
			req.ProductID, asOfDate,
		).Scan(&discountPct, &discountAmt, &fixedPrice)
		if err == nil {
			// Get base price first
			var basePrice float64
			s.db.QueryRow(ctx, `SELECT COALESCE(base_price, 0) FROM products WHERE id = $1`, req.ProductID).Scan(&basePrice)

			if fixedPrice > 0 {
				result.Price = fixedPrice
			} else if discountPct > 0 {
				result.Price = basePrice * (1 - discountPct/100)
				result.DiscountPct = discountPct
			} else if discountAmt > 0 {
				result.Price = basePrice - discountAmt
			}
			result.OriginalPrice = basePrice
			result.PriceLevel = models.PriceLevelPromotion
			result.PriceSource = "Promotion"
			priceFound = true
		}
	}

	// Level 4: Price Level (if customer has price level assigned)
	if req.CustomerID != nil && !priceFound {
		var priceLevel string
		s.db.QueryRow(ctx, `SELECT COALESCE(price_level, '') FROM customers WHERE id = $1`, *req.CustomerID).Scan(&priceLevel)
		if priceLevel != "" {
			err := s.db.QueryRow(ctx, `
				SELECT price FROM product_prices
				WHERE product_id = $1 AND price_level = $2 AND is_active = true
				AND effective_date <= $3 AND (expiry_date IS NULL OR expiry_date >= $3)
				ORDER BY effective_date DESC LIMIT 1`,
				req.ProductID, priceLevel, asOfDate,
			).Scan(&price)
			if err == nil {
				result.Price = price
				result.PriceLevel = models.PriceLevelPriceLevel
				result.PriceSource = "Price Level " + priceLevel
				priceFound = true
			}
		}
	}

	// Level 5: Base Price
	if !priceFound {
		err := s.db.QueryRow(ctx, `SELECT COALESCE(base_price, 0) FROM products WHERE id = $1`, req.ProductID).Scan(&price)
		if err != nil {
			return nil, fmt.Errorf("product not found")
		}
		result.Price = price
		result.PriceLevel = models.PriceLevelBase
		result.PriceSource = "Base Price"
	}

	// Calculate margin
	if result.Cost > 0 {
		result.Margin = result.Price - result.Cost
		result.MarginPercent = (result.Margin / result.Price) * 100
		result.IsBelowCost = result.Price < result.Cost
	}

	return result, nil
}

func (s *pricingServiceImpl) GetPricesForProducts(ctx context.Context, productIDs []int, customerID *int) ([]models.PriceLookupResult, error) {
	var results []models.PriceLookupResult
	for _, pid := range productIDs {
		req := &models.PriceLookupRequest{
			ProductID:  pid,
			CustomerID: customerID,
			Quantity:   1,
		}
		result, err := s.GetPrice(ctx, req)
		if err == nil {
			results = append(results, *result)
		}
	}
	return results, nil
}

// ============================================
// Product Prices
// ============================================

func (s *pricingServiceImpl) SetProductPrice(ctx context.Context, req *models.SetProductPriceRequest, createdBy int) (int, error) {
	effDate, _ := time.Parse("2006-01-02", req.EffectiveDate)
	var expDate *time.Time
	if req.ExpiryDate != "" {
		t, _ := time.Parse("2006-01-02", req.ExpiryDate)
		expDate = &t
	}

	query := `
		INSERT INTO product_prices (product_id, price_level, price, effective_date, expiry_date, min_quantity, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)
		ON CONFLICT (product_id, price_level, effective_date)
		DO UPDATE SET price = $3, expiry_date = $5, min_quantity = $6, is_active = true, updated_at = NOW()
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.ProductID, req.PriceLevel, req.Price, effDate, expDate, req.MinQuantity,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to set product price: %w", err)
	}
	return id, nil
}

func (s *pricingServiceImpl) GetProductPrices(ctx context.Context, productID int) ([]models.ProductPrice, error) {
	query := `
		SELECT id, product_id, price_level, price, effective_date, expiry_date, 
			   min_quantity, is_active, created_at, updated_at
		FROM product_prices
		WHERE product_id = $1
		ORDER BY price_level, effective_date DESC`

	rows := s.db.Query(ctx, query, productID)
	defer rows.Close()

	var prices []models.ProductPrice
	for rows.Next() {
		var p models.ProductPrice
		var expDate *time.Time

		err := rows.Scan(&p.ID, &p.ProductID, &p.PriceLevel, &p.Price, &p.EffectiveDate,
			&expDate, &p.MinQuantity, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product price: %w", err)
		}
		prices = append(prices, p)
	}

	return prices, nil
}

func (s *pricingServiceImpl) DeleteProductPrice(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx, `DELETE FROM product_prices WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete product price: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("price not found")
	}
	return nil
}

// ============================================
// Customer Prices
// ============================================

func (s *pricingServiceImpl) SetCustomerPrice(ctx context.Context, req *models.SetCustomerPriceRequest, createdBy int) (int, error) {
	effDate, _ := time.Parse("2006-01-02", req.EffectiveDate)
	var expDate *time.Time
	if req.ExpiryDate != "" {
		t, _ := time.Parse("2006-01-02", req.ExpiryDate)
		expDate = &t
	}

	query := `
		INSERT INTO customer_pricing (customer_id, product_id, price, effective_date, expiry_date, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (customer_id, product_id, effective_date)
		DO UPDATE SET price = $3, expiry_date = $5, notes = $6
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.CustomerID, req.ProductID, req.Price, effDate, expDate, req.Notes, createdBy,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to set customer price: %w", err)
	}
	return id, nil
}

func (s *pricingServiceImpl) GetCustomerPrices(ctx context.Context, customerID int) ([]models.CustomerPrice, error) {
	query := `
		SELECT cp.id, cp.customer_id, cp.product_id, cp.price, cp.effective_date, 
			   cp.expiry_date, cp.notes, cp.created_by, cp.created_at
		FROM customer_pricing cp
		WHERE cp.customer_id = $1
		ORDER BY cp.effective_date DESC`

	rows := s.db.Query(ctx, query, customerID)
	defer rows.Close()

	var prices []models.CustomerPrice
	for rows.Next() {
		var p models.CustomerPrice
		var expDate *time.Time
		var notes *string

		err := rows.Scan(&p.ID, &p.CustomerID, &p.ProductID, &p.Price, &p.EffectiveDate,
			&expDate, &notes, &p.CreatedBy, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer price: %w", err)
		}
		if notes != nil {
			p.Notes = *notes
		}
		prices = append(prices, p)
	}

	return prices, nil
}

func (s *pricingServiceImpl) DeleteCustomerPrice(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx, `DELETE FROM customer_pricing WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer price: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("price not found")
	}
	return nil
}

// ============================================
// Contract Prices
// ============================================

func (s *pricingServiceImpl) CreateContract(ctx context.Context, req *models.CreatePriceContractRequest, createdBy int) (int, error) {
	effDate, _ := time.Parse("2006-01-02", req.EffectiveDate)
	expDate, _ := time.Parse("2006-01-02", req.ExpiryDate)

	query := `
		INSERT INTO contract_prices (
			contract_code, customer_id, product_id, price, effective_date, expiry_date,
			min_quantity, max_quantity, notes, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.ContractCode, req.CustomerID, req.ProductID, req.Price,
		effDate, expDate, req.MinQuantity, req.MaxQuantity, req.Notes,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create contract: %w", err)
	}
	return id, nil
}

func (s *pricingServiceImpl) GetContract(ctx context.Context, id int) (*models.ContractPrice, error) {
	query := `
		SELECT id, contract_code, customer_id, product_id, price, effective_date, expiry_date,
			   min_quantity, max_quantity, notes, is_active, created_at
		FROM contract_prices
		WHERE id = $1`

	var c models.ContractPrice
	var notes *string

	err := s.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.ContractCode, &c.CustomerID, &c.ProductID, &c.Price,
		&c.EffectiveDate, &c.ExpiryDate, &c.MinQuantity, &c.MaxQuantity,
		&notes, &c.IsActive, &c.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("contract not found")
		}
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}
	if notes != nil {
		c.Notes = *notes
	}

	return &c, nil
}

func (s *pricingServiceImpl) ListContracts(ctx context.Context, customerID *int, activeOnly bool) ([]models.ContractPrice, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if customerID != nil {
		whereClause += fmt.Sprintf(" AND customer_id = $%d", argNum)
		args = append(args, *customerID)
		argNum++
	}
	if activeOnly {
		whereClause += " AND is_active = true AND expiry_date >= CURRENT_DATE"
	}

	query := fmt.Sprintf(`
		SELECT id, contract_code, customer_id, product_id, price, effective_date, expiry_date,
			   min_quantity, max_quantity, notes, is_active, created_at
		FROM contract_prices
		%s
		ORDER BY expiry_date DESC`, whereClause)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var contracts []models.ContractPrice
	for rows.Next() {
		var c models.ContractPrice
		var notes *string

		err := rows.Scan(&c.ID, &c.ContractCode, &c.CustomerID, &c.ProductID, &c.Price,
			&c.EffectiveDate, &c.ExpiryDate, &c.MinQuantity, &c.MaxQuantity,
			&notes, &c.IsActive, &c.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contract: %w", err)
		}
		if notes != nil {
			c.Notes = *notes
		}
		contracts = append(contracts, c)
	}

	return contracts, nil
}

func (s *pricingServiceImpl) DeactivateContract(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx, `UPDATE contract_prices SET is_active = false WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate contract: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("contract not found")
	}
	return nil
}

// ============================================
// Promotional Prices
// ============================================

func (s *pricingServiceImpl) CreatePromotion(ctx context.Context, req *models.CreatePromotionRequest, createdBy int) (int, error) {
	effDate, _ := time.Parse("2006-01-02", req.EffectiveDate)
	expDate, _ := time.Parse("2006-01-02", req.ExpiryDate)

	query := `
		INSERT INTO promotional_prices (
			promotion_code, name, product_id, category_id, discount_percent, discount_amount,
			fixed_price, effective_date, expiry_date, customer_group_id, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, true)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.PromotionCode, req.Name, req.ProductID, req.CategoryID,
		req.DiscountPercent, req.DiscountAmount, req.FixedPrice,
		effDate, expDate, req.CustomerGroupID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create promotion: %w", err)
	}
	return id, nil
}

func (s *pricingServiceImpl) GetPromotion(ctx context.Context, id int) (*models.PromotionalPrice, error) {
	query := `
		SELECT id, promotion_code, name, product_id, category_id, discount_percent,
			   discount_amount, fixed_price, effective_date, expiry_date, customer_group_id,
			   is_active, created_at
		FROM promotional_prices
		WHERE id = $1`

	var p models.PromotionalPrice
	err := s.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.PromotionCode, &p.Name, &p.ProductID, &p.CategoryID,
		&p.DiscountPercent, &p.DiscountAmount, &p.FixedPrice,
		&p.EffectiveDate, &p.ExpiryDate, &p.CustomerGroupID, &p.IsActive, &p.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("promotion not found")
		}
		return nil, fmt.Errorf("failed to get promotion: %w", err)
	}

	return &p, nil
}

func (s *pricingServiceImpl) ListPromotions(ctx context.Context, activeOnly bool) ([]models.PromotionalPrice, error) {
	whereClause := ""
	if activeOnly {
		whereClause = "WHERE is_active = true AND expiry_date >= CURRENT_DATE"
	}

	query := fmt.Sprintf(`
		SELECT id, promotion_code, name, product_id, category_id, discount_percent,
			   discount_amount, fixed_price, effective_date, expiry_date, customer_group_id,
			   is_active, created_at
		FROM promotional_prices
		%s
		ORDER BY effective_date DESC`, whereClause)

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var promotions []models.PromotionalPrice
	for rows.Next() {
		var p models.PromotionalPrice
		err := rows.Scan(&p.ID, &p.PromotionCode, &p.Name, &p.ProductID, &p.CategoryID,
			&p.DiscountPercent, &p.DiscountAmount, &p.FixedPrice,
			&p.EffectiveDate, &p.ExpiryDate, &p.CustomerGroupID, &p.IsActive, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan promotion: %w", err)
		}
		promotions = append(promotions, p)
	}

	return promotions, nil
}

func (s *pricingServiceImpl) DeactivatePromotion(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx, `UPDATE promotional_prices SET is_active = false WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate promotion: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("promotion not found")
	}
	return nil
}

// ============================================
// Product Costs
// ============================================

func (s *pricingServiceImpl) UpdateProductCost(ctx context.Context, req *models.UpdateProductCostRequest, updatedBy int) error {
	landedCost := req.Cost * (1 + req.FreightFactor/100 + req.DutyFactor/100 + req.HandlingFactor/100)

	query := `
		INSERT INTO product_costs (
			product_id, costing_method, cost, effective_date, freight_factor,
			duty_factor, handling_factor, landed_cost, notes, updated_by
		) VALUES ($1, $2, $3, CURRENT_DATE, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (product_id, costing_method)
		DO UPDATE SET cost = $3, freight_factor = $4, duty_factor = $5,
			handling_factor = $6, landed_cost = $7, notes = $8, updated_by = $9, updated_at = NOW()`

	_, err := s.db.Exec(ctx, query,
		req.ProductID, req.CostingMethod, req.Cost, req.FreightFactor,
		req.DutyFactor, req.HandlingFactor, landedCost, req.Notes, updatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to update product cost: %w", err)
	}
	return nil
}

func (s *pricingServiceImpl) GetProductCost(ctx context.Context, productID int) (*models.ProductCost, error) {
	query := `
		SELECT id, product_id, costing_method, cost, effective_date, freight_factor,
			   duty_factor, handling_factor, landed_cost, notes, updated_by, updated_at
		FROM product_costs
		WHERE product_id = $1
		ORDER BY updated_at DESC LIMIT 1`

	var c models.ProductCost
	var notes *string

	err := s.db.QueryRow(ctx, query, productID).Scan(
		&c.ID, &c.ProductID, &c.CostingMethod, &c.Cost, &c.EffectiveDate,
		&c.FreightFactor, &c.DutyFactor, &c.HandlingFactor, &c.LandedCost,
		&notes, &c.UpdatedBy, &c.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("cost not found")
		}
		return nil, fmt.Errorf("failed to get product cost: %w", err)
	}
	if notes != nil {
		c.Notes = *notes
	}

	return &c, nil
}

// ============================================
// Mass Price Maintenance
// ============================================

func (s *pricingServiceImpl) MassPriceUpdate(ctx context.Context, req *models.MassPriceUpdateRequest, createdBy int) (int, error) {
	effDate, _ := time.Parse("2006-01-02", req.EffectiveDate)

	// Get products to update
	var productIDs []int
	if len(req.ProductIDs) > 0 {
		productIDs = req.ProductIDs
	} else if req.CategoryID != nil {
		rows := s.db.Query(ctx, `SELECT id FROM products WHERE category_id = $1 AND is_active = true`, *req.CategoryID)
		defer rows.Close()
		for rows.Next() {
			var id int
			rows.Scan(&id)
			productIDs = append(productIDs, id)
		}
	}

	count := 0
	for _, productID := range productIDs {
		// Get current price
		var currentPrice float64
		s.db.QueryRow(ctx, `
			SELECT price FROM product_prices 
			WHERE product_id = $1 AND price_level = $2 AND is_active = true
			ORDER BY effective_date DESC LIMIT 1`, productID, req.PriceLevel).Scan(&currentPrice)

		if currentPrice == 0 {
			s.db.QueryRow(ctx, `SELECT COALESCE(base_price, 0) FROM products WHERE id = $1`, productID).Scan(&currentPrice)
		}

		// Calculate new price
		var newPrice float64
		if req.AdjustmentType == "PERCENT" {
			newPrice = currentPrice * (1 + req.AdjustmentValue/100)
		} else {
			newPrice = currentPrice + req.AdjustmentValue
		}

		// Insert new price
		_, err := s.db.Exec(ctx, `
			INSERT INTO product_prices (product_id, price_level, price, effective_date, is_active)
			VALUES ($1, $2, $3, $4, true)
			ON CONFLICT (product_id, price_level, effective_date)
			DO UPDATE SET price = $3`,
			productID, req.PriceLevel, newPrice, effDate,
		)
		if err == nil {
			count++
		}
	}

	return count, nil
}

// ============================================
// Price Lists
// ============================================

func (s *pricingServiceImpl) GetPriceList(ctx context.Context, filters *models.PriceListFilters) ([]models.PriceListItem, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 50
	}

	whereClause := "WHERE p.is_active = true"
	args := []interface{}{}
	argNum := 1

	if filters.CategoryID != nil {
		whereClause += fmt.Sprintf(" AND p.category_id = $%d", argNum)
		args = append(args, *filters.CategoryID)
		argNum++
	}

	// Count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM products p %s`, whereClause)
	var total int64
	s.db.QueryRow(ctx, countQuery, args...).Scan(&total)

	// Query
	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT p.id, p.sku, p.name, COALESCE(pc.name, '') as category_name,
			   COALESCE(pu.abbreviation, 'EA') as unit_of_measure,
			   COALESCE(p.base_price, 0) as base_price,
			   COALESCE(i.average_cost, 0) as cost
		FROM products p
		LEFT JOIN product_categories pc ON p.category_id = pc.id
		LEFT JOIN product_units pu ON p.default_unit_id = pu.id
		LEFT JOIN inventory i ON i.product_id = p.id
		%s
		GROUP BY p.id, p.sku, p.name, pc.name, pu.abbreviation, p.base_price, i.average_cost
		ORDER BY p.name
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var items []models.PriceListItem
	for rows.Next() {
		var item models.PriceListItem
		err := rows.Scan(&item.ProductID, &item.ProductSKU, &item.ProductName,
			&item.CategoryName, &item.UnitOfMeasure, &item.BasePrice, &item.Cost)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan price list item: %w", err)
		}
		items = append(items, item)
	}

	return items, total, nil
}

// ============================================
// Margin Check
// ============================================

func (s *pricingServiceImpl) CheckBelowCost(ctx context.Context, productID int, price float64) (bool, float64, error) {
	var cost float64
	err := s.db.QueryRow(ctx, `SELECT COALESCE(average_cost, 0) FROM inventory WHERE product_id = $1 LIMIT 1`, productID).Scan(&cost)
	if err != nil {
		return false, 0, nil
	}
	return price < cost, cost, nil
}
