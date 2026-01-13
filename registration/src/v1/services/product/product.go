package product

import (
	"context"
	"fmt"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

// ============================================
// Service Interface
// ============================================

type ProductService interface {
	// Products
	Create(ctx context.Context, req *models.CreateProductRequest) (int, error)
	GetByID(ctx context.Context, id int) (*models.ProductWithDetails, error)
	GetBySKU(ctx context.Context, sku string) (*models.Product, error)
	GetByBarcode(ctx context.Context, barcode string) (*models.Product, error)
	Update(ctx context.Context, id int, req *models.UpdateProductRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters *models.ProductListFilters) ([]models.Product, int64, error)

	// Categories
	CreateCategory(ctx context.Context, req *models.ProductCategory) (int, error)
	GetCategoryByID(ctx context.Context, id int) (*models.ProductCategory, error)
	UpdateCategory(ctx context.Context, id int, req *models.ProductCategory) error
	DeleteCategory(ctx context.Context, id int) error
	ListCategories(ctx context.Context) ([]models.ProductCategory, error)

	// Units
	AddUnit(ctx context.Context, req *models.ProductUnit) (int, error)
	GetUnits(ctx context.Context, productID int) ([]models.ProductUnit, error)
	UpdateUnit(ctx context.Context, id int, req *models.ProductUnit) error
	DeleteUnit(ctx context.Context, id int) error
}

// ============================================
// Service Implementation
// ============================================

type productServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) ProductService {
	return &productServiceImpl{db: db}
}

// ============================================
// Product CRUD
// ============================================

func (s *productServiceImpl) Create(ctx context.Context, req *models.CreateProductRequest) (int, error) {
	query := `
		INSERT INTO products (
			sku, barcode, upc, name, description, category_id,
			base_unit, is_catch_weight, catch_weight_unit, country_of_origin,
			shelf_life_days, min_shelf_life_days, is_lot_tracked, is_serialized,
			haccp_category, qc_required, is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, true
		) RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.SKU,
		req.Barcode,
		req.UPC,
		req.Name,
		req.Description,
		req.CategoryID,
		req.BaseUnit,
		req.IsCatchWeight,
		req.CatchWeightUnit,
		req.CountryOfOrigin,
		req.ShelfLifeDays,
		req.MinShelfLifeDays,
		req.IsLotTracked,
		req.IsSerialized,
		req.HACCPCategory,
		req.QCRequired,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create product: %w", err)
	}
	return id, nil
}

func (s *productServiceImpl) GetByID(ctx context.Context, id int) (*models.ProductWithDetails, error) {
	// Get product
	product, err := s.getProduct(ctx, "id = $1", id)
	if err != nil {
		return nil, err
	}

	result := &models.ProductWithDetails{
		Product: *product,
	}

	// Get category if exists
	if product.CategoryID != nil {
		category, err := s.GetCategoryByID(ctx, *product.CategoryID)
		if err == nil {
			result.Category = category
		}
	}

	// Get units
	units, err := s.GetUnits(ctx, id)
	if err == nil {
		result.Units = units
	}

	// Get current stock (sum from inventory)
	stockQuery := `SELECT COALESCE(SUM(quantity_on_hand), 0) FROM inventory WHERE product_id = $1`
	s.db.QueryRow(ctx, stockQuery, id).Scan(&result.CurrentStock)

	// Get average cost
	costQuery := `SELECT COALESCE(AVG(average_cost), 0) FROM inventory WHERE product_id = $1`
	s.db.QueryRow(ctx, costQuery, id).Scan(&result.AvgCost)

	return result, nil
}

func (s *productServiceImpl) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	return s.getProduct(ctx, "sku = $1", sku)
}

func (s *productServiceImpl) GetByBarcode(ctx context.Context, barcode string) (*models.Product, error) {
	return s.getProduct(ctx, "barcode = $1", barcode)
}

func (s *productServiceImpl) getProduct(ctx context.Context, whereClause string, arg interface{}) (*models.Product, error) {
	query := fmt.Sprintf(`
		SELECT id, sku, barcode, upc, name, description, category_id,
			   base_unit, is_catch_weight, catch_weight_unit, country_of_origin,
			   shelf_life_days, min_shelf_life_days, is_lot_tracked, is_serialized,
			   haccp_category, qc_required, is_active, created_at, updated_at
		FROM products
		WHERE %s`, whereClause)

	var p models.Product
	var barcode, upc, description, catchWeightUnit, countryOfOrigin, haccpCategory *string
	var shelfLifeDays, minShelfLifeDays *int

	err := s.db.QueryRow(ctx, query, arg).Scan(
		&p.ID, &p.SKU, &barcode, &upc, &p.Name, &description, &p.CategoryID,
		&p.BaseUnit, &p.IsCatchWeight, &catchWeightUnit, &countryOfOrigin,
		&shelfLifeDays, &minShelfLifeDays, &p.IsLotTracked, &p.IsSerialized,
		&haccpCategory, &p.QCRequired, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Handle nullable fields
	if barcode != nil {
		p.Barcode = *barcode
	}
	if upc != nil {
		p.UPC = *upc
	}
	if description != nil {
		p.Description = *description
	}
	if catchWeightUnit != nil {
		p.CatchWeightUnit = *catchWeightUnit
	}
	if countryOfOrigin != nil {
		p.CountryOfOrigin = *countryOfOrigin
	}
	if haccpCategory != nil {
		p.HACCPCategory = *haccpCategory
	}
	if shelfLifeDays != nil {
		p.ShelfLifeDays = *shelfLifeDays
	}
	if minShelfLifeDays != nil {
		p.MinShelfLifeDays = *minShelfLifeDays
	}

	return &p, nil
}

func (s *productServiceImpl) Update(ctx context.Context, id int, req *models.UpdateProductRequest) error {
	query := `
		UPDATE products SET
			name = COALESCE($1, name),
			description = COALESCE($2, description),
			category_id = COALESCE($3, category_id),
			base_unit = COALESCE($4, base_unit),
			country_of_origin = COALESCE($5, country_of_origin),
			shelf_life_days = COALESCE($6, shelf_life_days),
			min_shelf_life_days = COALESCE($7, min_shelf_life_days),
			haccp_category = COALESCE($8, haccp_category),
			qc_required = COALESCE($9, qc_required),
			is_active = COALESCE($10, is_active),
			updated_at = NOW()
		WHERE id = $11`

	result, err := s.db.Exec(ctx, query,
		req.Name,
		req.Description,
		req.CategoryID,
		req.BaseUnit,
		req.CountryOfOrigin,
		req.ShelfLifeDays,
		req.MinShelfLifeDays,
		req.HACCPCategory,
		req.QCRequired,
		req.IsActive,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (s *productServiceImpl) Delete(ctx context.Context, id int) error {
	// Soft delete by setting is_active = false
	query := `UPDATE products SET is_active = false, updated_at = NOW() WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (s *productServiceImpl) List(ctx context.Context, filters *models.ProductListFilters) ([]models.Product, int64, error) {
	// Set defaults
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		filters.PageSize = 20
	}

	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if filters.Search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d OR barcode ILIKE $%d)", argNum, argNum, argNum)
		args = append(args, "%"+filters.Search+"%")
		argNum++
	}

	if filters.CategoryID != nil {
		whereClause += fmt.Sprintf(" AND category_id = $%d", argNum)
		args = append(args, *filters.CategoryID)
		argNum++
	}

	if filters.IsActive != nil {
		whereClause += fmt.Sprintf(" AND is_active = $%d", argNum)
		args = append(args, *filters.IsActive)
		argNum++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Get paginated results
	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT id, sku, barcode, upc, name, description, category_id,
			   base_unit, is_catch_weight, catch_weight_unit, country_of_origin,
			   shelf_life_days, min_shelf_life_days, is_lot_tracked, is_serialized,
			   haccp_category, qc_required, is_active, created_at, updated_at
		FROM products
		%s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var barcode, upc, description, catchWeightUnit, countryOfOrigin, haccpCategory *string
		var shelfLifeDays, minShelfLifeDays *int

		err := rows.Scan(
			&p.ID, &p.SKU, &barcode, &upc, &p.Name, &description, &p.CategoryID,
			&p.BaseUnit, &p.IsCatchWeight, &catchWeightUnit, &countryOfOrigin,
			&shelfLifeDays, &minShelfLifeDays, &p.IsLotTracked, &p.IsSerialized,
			&haccpCategory, &p.QCRequired, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}

		// Handle nullable fields
		if barcode != nil {
			p.Barcode = *barcode
		}
		if upc != nil {
			p.UPC = *upc
		}
		if description != nil {
			p.Description = *description
		}
		if catchWeightUnit != nil {
			p.CatchWeightUnit = *catchWeightUnit
		}
		if countryOfOrigin != nil {
			p.CountryOfOrigin = *countryOfOrigin
		}
		if haccpCategory != nil {
			p.HACCPCategory = *haccpCategory
		}
		if shelfLifeDays != nil {
			p.ShelfLifeDays = *shelfLifeDays
		}
		if minShelfLifeDays != nil {
			p.MinShelfLifeDays = *minShelfLifeDays
		}

		products = append(products, p)
	}

	return products, total, nil
}

// ============================================
// Category CRUD
// ============================================

func (s *productServiceImpl) CreateCategory(ctx context.Context, req *models.ProductCategory) (int, error) {
	query := `
		INSERT INTO product_categories (code, name, parent_id, gl_sales_account_id, gl_cogs_account_id, gl_inventory_account_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.Code, req.Name, req.ParentID,
		req.GLSalesAccountID, req.GLCOGSAccountID, req.GLInventoryAccountID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create category: %w", err)
	}
	return id, nil
}

func (s *productServiceImpl) GetCategoryByID(ctx context.Context, id int) (*models.ProductCategory, error) {
	query := `
		SELECT id, code, name, parent_id, gl_sales_account_id, gl_cogs_account_id, gl_inventory_account_id
		FROM product_categories
		WHERE id = $1`

	var c models.ProductCategory
	var code *string
	err := s.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &code, &c.Name, &c.ParentID,
		&c.GLSalesAccountID, &c.GLCOGSAccountID, &c.GLInventoryAccountID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if code != nil {
		c.Code = *code
	}

	return &c, nil
}

func (s *productServiceImpl) UpdateCategory(ctx context.Context, id int, req *models.ProductCategory) error {
	query := `
		UPDATE product_categories SET
			code = $1, name = $2, parent_id = $3,
			gl_sales_account_id = $4, gl_cogs_account_id = $5, gl_inventory_account_id = $6
		WHERE id = $7`

	result, err := s.db.Exec(ctx, query,
		req.Code, req.Name, req.ParentID,
		req.GLSalesAccountID, req.GLCOGSAccountID, req.GLInventoryAccountID, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (s *productServiceImpl) DeleteCategory(ctx context.Context, id int) error {
	// Check if category has products
	var count int
	s.db.QueryRow(ctx, "SELECT COUNT(*) FROM products WHERE category_id = $1", id).Scan(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete category with %d products", count)
	}

	// Check if category has children
	s.db.QueryRow(ctx, "SELECT COUNT(*) FROM product_categories WHERE parent_id = $1", id).Scan(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete category with %d subcategories", count)
	}

	query := `DELETE FROM product_categories WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (s *productServiceImpl) ListCategories(ctx context.Context) ([]models.ProductCategory, error) {
	query := `
		SELECT id, code, name, parent_id, gl_sales_account_id, gl_cogs_account_id, gl_inventory_account_id
		FROM product_categories
		ORDER BY name ASC`

	rows := s.db.Query(ctx, query)
	defer rows.Close()

	var categories []models.ProductCategory
	for rows.Next() {
		var c models.ProductCategory
		var code *string
		err := rows.Scan(
			&c.ID, &code, &c.Name, &c.ParentID,
			&c.GLSalesAccountID, &c.GLCOGSAccountID, &c.GLInventoryAccountID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		if code != nil {
			c.Code = *code
		}
		categories = append(categories, c)
	}

	return categories, nil
}

// ============================================
// Unit CRUD
// ============================================

func (s *productServiceImpl) AddUnit(ctx context.Context, req *models.ProductUnit) (int, error) {
	query := `
		INSERT INTO product_units (product_id, unit_name, description, conversion_factor, barcode, weight, is_purchase_unit, is_sales_unit)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.ProductID, req.UnitName, req.Description, req.ConversionFactor,
		req.Barcode, req.Weight, req.IsPurchaseUnit, req.IsSalesUnit,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to add unit: %w", err)
	}
	return id, nil
}

func (s *productServiceImpl) GetUnits(ctx context.Context, productID int) ([]models.ProductUnit, error) {
	query := `
		SELECT id, product_id, unit_name, description, conversion_factor, barcode, weight, is_purchase_unit, is_sales_unit
		FROM product_units
		WHERE product_id = $1
		ORDER BY conversion_factor ASC`

	rows := s.db.Query(ctx, query, productID)
	defer rows.Close()

	var units []models.ProductUnit
	for rows.Next() {
		var u models.ProductUnit
		var description, barcode *string
		var weight *float64

		err := rows.Scan(
			&u.ID, &u.ProductID, &u.UnitName, &description, &u.ConversionFactor,
			&barcode, &weight, &u.IsPurchaseUnit, &u.IsSalesUnit,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan unit: %w", err)
		}

		if description != nil {
			u.Description = *description
		}
		if barcode != nil {
			u.Barcode = *barcode
		}
		if weight != nil {
			u.Weight = *weight
		}

		units = append(units, u)
	}

	return units, nil
}

func (s *productServiceImpl) UpdateUnit(ctx context.Context, id int, req *models.ProductUnit) error {
	query := `
		UPDATE product_units SET
			unit_name = $1, description = $2, conversion_factor = $3,
			barcode = $4, weight = $5, is_purchase_unit = $6, is_sales_unit = $7
		WHERE id = $8`

	result, err := s.db.Exec(ctx, query,
		req.UnitName, req.Description, req.ConversionFactor,
		req.Barcode, req.Weight, req.IsPurchaseUnit, req.IsSalesUnit, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update unit: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("unit not found")
	}

	return nil
}

func (s *productServiceImpl) DeleteUnit(ctx context.Context, id int) error {
	query := `DELETE FROM product_units WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("unit not found")
	}

	return nil
}
