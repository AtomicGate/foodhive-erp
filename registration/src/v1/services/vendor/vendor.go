package vendor

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

type VendorService interface {
	// Vendors
	Create(ctx context.Context, req *models.CreateVendorRequest) (int, error)
	GetByID(ctx context.Context, id int) (*models.VendorWithDetails, error)
	GetByCode(ctx context.Context, code string) (*models.Vendor, error)
	Update(ctx context.Context, id int, req *models.UpdateVendorRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters *models.VendorListFilters) ([]models.Vendor, int64, error)

	// Vendor Products
	AddProduct(ctx context.Context, req *models.VendorProduct) (int, error)
	GetProducts(ctx context.Context, vendorID int) ([]models.VendorProduct, error)
	GetProductsByProductID(ctx context.Context, productID int) ([]models.VendorProduct, error)
	UpdateProduct(ctx context.Context, id int, req *models.VendorProduct) error
	DeleteProduct(ctx context.Context, id int) error

	// Vendor Discounts
	AddDiscount(ctx context.Context, req *models.VendorDiscount) (int, error)
	GetDiscounts(ctx context.Context, vendorID int) ([]models.VendorDiscount, error)
	UpdateDiscount(ctx context.Context, id int, req *models.VendorDiscount) error
	DeleteDiscount(ctx context.Context, id int) error
}

// ============================================
// Service Implementation
// ============================================

type vendorServiceImpl struct {
	db postgres.Executor
}

func New(db postgres.Executor) VendorService {
	return &vendorServiceImpl{db: db}
}

// ============================================
// Vendor CRUD
// ============================================

func (s *vendorServiceImpl) Create(ctx context.Context, req *models.CreateVendorRequest) (int, error) {
	// Set defaults
	if req.Currency == "" {
		req.Currency = "USD"
	}
	if req.PaymentTermsDays == 0 {
		req.PaymentTermsDays = 30
	}
	if req.LeadTimeDays == 0 {
		req.LeadTimeDays = 7
	}

	query := `
		INSERT INTO vendors (
			vendor_code, name, address_line1, address_line2, city, state,
			postal_code, country, phone, email, payment_terms_days, currency,
			lead_time_days, minimum_order, buyer_id, is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, true
		) RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.VendorCode, req.Name, req.AddressLine1, req.AddressLine2,
		req.City, req.State, req.PostalCode, req.Country,
		req.Phone, req.Email, req.PaymentTermsDays, req.Currency,
		req.LeadTimeDays, req.MinimumOrder, req.BuyerID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create vendor: %w", err)
	}
	return id, nil
}

func (s *vendorServiceImpl) GetByID(ctx context.Context, id int) (*models.VendorWithDetails, error) {
	vendor, err := s.getVendor(ctx, "id = $1", id)
	if err != nil {
		return nil, err
	}

	result := &models.VendorWithDetails{
		Vendor: *vendor,
	}

	// Get buyer name if exists
	if vendor.BuyerID != nil {
		var buyerName string
		buyerQuery := `SELECT COALESCE(english_name, email) FROM employees WHERE id = $1`
		s.db.QueryRow(ctx, buyerQuery, *vendor.BuyerID).Scan(&buyerName)
		result.BuyerName = buyerName
	}

	// Get products
	products, _ := s.GetProducts(ctx, id)
	result.Products = products

	// Get discounts
	discounts, _ := s.GetDiscounts(ctx, id)
	result.Discounts = discounts

	return result, nil
}

func (s *vendorServiceImpl) GetByCode(ctx context.Context, code string) (*models.Vendor, error) {
	return s.getVendor(ctx, "vendor_code = $1", code)
}

func (s *vendorServiceImpl) getVendor(ctx context.Context, whereClause string, arg interface{}) (*models.Vendor, error) {
	query := fmt.Sprintf(`
		SELECT id, vendor_code, name, address_line1, address_line2, city, state,
			   postal_code, country, phone, email, payment_terms_days, currency,
			   lead_time_days, minimum_order, buyer_id, is_active, created_at, updated_at
		FROM vendors
		WHERE %s`, whereClause)

	var v models.Vendor
	var addr1, addr2, city, state, postal, country, phone, email *string
	var minOrder *float64

	err := s.db.QueryRow(ctx, query, arg).Scan(
		&v.ID, &v.VendorCode, &v.Name, &addr1, &addr2, &city, &state,
		&postal, &country, &phone, &email, &v.PaymentTermsDays, &v.Currency,
		&v.LeadTimeDays, &minOrder, &v.BuyerID, &v.IsActive, &v.CreatedAt, &v.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("vendor not found")
		}
		return nil, fmt.Errorf("failed to get vendor: %w", err)
	}

	// Handle nullable fields
	if addr1 != nil {
		v.AddressLine1 = *addr1
	}
	if addr2 != nil {
		v.AddressLine2 = *addr2
	}
	if city != nil {
		v.City = *city
	}
	if state != nil {
		v.State = *state
	}
	if postal != nil {
		v.PostalCode = *postal
	}
	if country != nil {
		v.Country = *country
	}
	if phone != nil {
		v.Phone = *phone
	}
	if email != nil {
		v.Email = *email
	}
	if minOrder != nil {
		v.MinimumOrder = *minOrder
	}

	return &v, nil
}

func (s *vendorServiceImpl) Update(ctx context.Context, id int, req *models.UpdateVendorRequest) error {
	query := `
		UPDATE vendors SET
			name = COALESCE($1, name),
			address_line1 = COALESCE($2, address_line1),
			address_line2 = COALESCE($3, address_line2),
			city = COALESCE($4, city),
			state = COALESCE($5, state),
			postal_code = COALESCE($6, postal_code),
			country = COALESCE($7, country),
			phone = COALESCE($8, phone),
			email = COALESCE($9, email),
			payment_terms_days = COALESCE($10, payment_terms_days),
			currency = COALESCE($11, currency),
			lead_time_days = COALESCE($12, lead_time_days),
			minimum_order = COALESCE($13, minimum_order),
			buyer_id = COALESCE($14, buyer_id),
			is_active = COALESCE($15, is_active),
			updated_at = NOW()
		WHERE id = $16`

	result, err := s.db.Exec(ctx, query,
		req.Name, req.AddressLine1, req.AddressLine2, req.City, req.State,
		req.PostalCode, req.Country, req.Phone, req.Email,
		req.PaymentTermsDays, req.Currency, req.LeadTimeDays,
		req.MinimumOrder, req.BuyerID, req.IsActive, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update vendor: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("vendor not found")
	}

	return nil
}

func (s *vendorServiceImpl) Delete(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE vendors SET is_active = false, updated_at = NOW() WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vendor: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("vendor not found")
	}
	return nil
}

func (s *vendorServiceImpl) List(ctx context.Context, filters *models.VendorListFilters) ([]models.Vendor, int64, error) {
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
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR vendor_code ILIKE $%d)", argNum, argNum)
		args = append(args, "%"+filters.Search+"%")
		argNum++
	}

	if filters.BuyerID != nil {
		whereClause += fmt.Sprintf(" AND buyer_id = $%d", argNum)
		args = append(args, *filters.BuyerID)
		argNum++
	}

	if filters.IsActive != nil {
		whereClause += fmt.Sprintf(" AND is_active = $%d", argNum)
		args = append(args, *filters.IsActive)
		argNum++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM vendors %s", whereClause)
	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vendors: %w", err)
	}

	// Get paginated results
	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT id, vendor_code, name, address_line1, address_line2, city, state,
			   postal_code, country, phone, email, payment_terms_days, currency,
			   lead_time_days, minimum_order, buyer_id, is_active, created_at, updated_at
		FROM vendors
		%s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, filters.PageSize, offset)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		var v models.Vendor
		var addr1, addr2, city, state, postal, country, phone, email *string
		var minOrder *float64

		err := rows.Scan(
			&v.ID, &v.VendorCode, &v.Name, &addr1, &addr2, &city, &state,
			&postal, &country, &phone, &email, &v.PaymentTermsDays, &v.Currency,
			&v.LeadTimeDays, &minOrder, &v.BuyerID, &v.IsActive, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vendor: %w", err)
		}

		if addr1 != nil {
			v.AddressLine1 = *addr1
		}
		if addr2 != nil {
			v.AddressLine2 = *addr2
		}
		if city != nil {
			v.City = *city
		}
		if state != nil {
			v.State = *state
		}
		if postal != nil {
			v.PostalCode = *postal
		}
		if country != nil {
			v.Country = *country
		}
		if phone != nil {
			v.Phone = *phone
		}
		if email != nil {
			v.Email = *email
		}
		if minOrder != nil {
			v.MinimumOrder = *minOrder
		}

		vendors = append(vendors, v)
	}

	return vendors, total, nil
}

// ============================================
// Vendor Products CRUD
// ============================================

func (s *vendorServiceImpl) AddProduct(ctx context.Context, req *models.VendorProduct) (int, error) {
	query := `
		INSERT INTO vendor_products (
			vendor_id, product_id, vendor_sku, vendor_description,
			unit_of_measure, unit_cost, minimum_order_qty, lead_time_days, is_preferred
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.VendorID, req.ProductID, req.VendorSKU, req.VendorDescription,
		req.UnitOfMeasure, req.UnitCost, req.MinimumOrderQty, req.LeadTimeDays, req.IsPreferred,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to add vendor product: %w", err)
	}
	return id, nil
}

func (s *vendorServiceImpl) GetProducts(ctx context.Context, vendorID int) ([]models.VendorProduct, error) {
	return s.getVendorProducts(ctx, "vendor_id = $1", vendorID)
}

func (s *vendorServiceImpl) GetProductsByProductID(ctx context.Context, productID int) ([]models.VendorProduct, error) {
	return s.getVendorProducts(ctx, "product_id = $1", productID)
}

func (s *vendorServiceImpl) getVendorProducts(ctx context.Context, whereClause string, arg interface{}) ([]models.VendorProduct, error) {
	query := fmt.Sprintf(`
		SELECT id, vendor_id, product_id, vendor_sku, vendor_description,
			   unit_of_measure, unit_cost, minimum_order_qty, lead_time_days, is_preferred
		FROM vendor_products
		WHERE %s
		ORDER BY is_preferred DESC, unit_cost ASC`, whereClause)

	rows := s.db.Query(ctx, query, arg)
	defer rows.Close()

	var products []models.VendorProduct
	for rows.Next() {
		var p models.VendorProduct
		var sku, desc, uom *string
		var minQty *float64

		err := rows.Scan(
			&p.ID, &p.VendorID, &p.ProductID, &sku, &desc,
			&uom, &p.UnitCost, &minQty, &p.LeadTimeDays, &p.IsPreferred,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vendor product: %w", err)
		}

		if sku != nil {
			p.VendorSKU = *sku
		}
		if desc != nil {
			p.VendorDescription = *desc
		}
		if uom != nil {
			p.UnitOfMeasure = *uom
		}
		if minQty != nil {
			p.MinimumOrderQty = *minQty
		}

		products = append(products, p)
	}

	return products, nil
}

func (s *vendorServiceImpl) UpdateProduct(ctx context.Context, id int, req *models.VendorProduct) error {
	query := `
		UPDATE vendor_products SET
			vendor_sku = $1, vendor_description = $2, unit_of_measure = $3,
			unit_cost = $4, minimum_order_qty = $5, lead_time_days = $6, is_preferred = $7
		WHERE id = $8`

	result, err := s.db.Exec(ctx, query,
		req.VendorSKU, req.VendorDescription, req.UnitOfMeasure,
		req.UnitCost, req.MinimumOrderQty, req.LeadTimeDays, req.IsPreferred, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update vendor product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("vendor product not found")
	}

	return nil
}

func (s *vendorServiceImpl) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM vendor_products WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vendor product: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("vendor product not found")
	}
	return nil
}

// ============================================
// Vendor Discounts CRUD
// ============================================

func (s *vendorServiceImpl) AddDiscount(ctx context.Context, req *models.VendorDiscount) (int, error) {
	query := `
		INSERT INTO vendor_discounts (vendor_id, discount_days, discount_percent)
		VALUES ($1, $2, $3)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.VendorID, req.DiscountDays, req.DiscountPercent,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to add vendor discount: %w", err)
	}
	return id, nil
}

func (s *vendorServiceImpl) GetDiscounts(ctx context.Context, vendorID int) ([]models.VendorDiscount, error) {
	query := `
		SELECT id, vendor_id, discount_days, discount_percent
		FROM vendor_discounts
		WHERE vendor_id = $1
		ORDER BY discount_days ASC`

	rows := s.db.Query(ctx, query, vendorID)
	defer rows.Close()

	var discounts []models.VendorDiscount
	for rows.Next() {
		var d models.VendorDiscount
		err := rows.Scan(&d.ID, &d.VendorID, &d.DiscountDays, &d.DiscountPercent)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vendor discount: %w", err)
		}
		discounts = append(discounts, d)
	}

	return discounts, nil
}

func (s *vendorServiceImpl) UpdateDiscount(ctx context.Context, id int, req *models.VendorDiscount) error {
	query := `
		UPDATE vendor_discounts SET
			discount_days = $1, discount_percent = $2
		WHERE id = $3`

	result, err := s.db.Exec(ctx, query, req.DiscountDays, req.DiscountPercent, id)
	if err != nil {
		return fmt.Errorf("failed to update vendor discount: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("vendor discount not found")
	}

	return nil
}

func (s *vendorServiceImpl) DeleteDiscount(ctx context.Context, id int) error {
	query := `DELETE FROM vendor_discounts WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vendor discount: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("vendor discount not found")
	}
	return nil
}
