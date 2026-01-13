package customer

import (
	"context"
	"errors"
	"fmt"

	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound      = errors.New("customer not found")
	ErrDuplicateCode = errors.New("customer code already exists")
)

type CustomerService interface {
	Create(ctx context.Context, req models.CreateCustomerRequest, createdBy int) (int, error)
	GetByID(ctx context.Context, id int) (*models.CustomerWithDetails, error)
	GetByCode(ctx context.Context, code string) (*models.Customer, error)
	Update(ctx context.Context, id int, req models.UpdateCustomerRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filters models.CustomerListFilters) ([]models.Customer, int64, error)
	GetOrderGuide(ctx context.Context, customerID int) ([]models.CustomerOrderGuide, error)
	AddShipTo(ctx context.Context, customerID int, shipTo models.CustomerShipTo) (int, error)
}

type customerServiceImpl struct {
	db postgres.Connection
}

func New(db postgres.Connection) CustomerService {
	return &customerServiceImpl{db: db}
}

func (s *customerServiceImpl) Create(ctx context.Context, req models.CreateCustomerRequest, createdBy int) (int, error) {
	// Set defaults
	if req.Currency == "" {
		req.Currency = "USD"
	}
	if req.PaymentTermsDays == 0 {
		req.PaymentTermsDays = 30
	}

	query := `
		INSERT INTO customers (
			customer_code, name, credit_limit, payment_terms_days, 
			currency, sales_rep_id, default_warehouse_id, tax_exempt, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		req.CustomerCode,
		req.Name,
		req.CreditLimit,
		req.PaymentTermsDays,
		req.Currency,
		req.SalesRepID,
		req.DefaultWarehouseID,
		req.TaxExempt,
		createdBy,
	).Scan(&id)

	if err != nil {
		// Check for duplicate code
		if err.Error() == "ERROR: duplicate key value violates unique constraint" {
			return 0, ErrDuplicateCode
		}
		return 0, fmt.Errorf("failed to create customer: %w", err)
	}

	return id, nil
}

func (s *customerServiceImpl) GetByID(ctx context.Context, id int) (*models.CustomerWithDetails, error) {
	query := `
		SELECT 
			c.id, c.customer_code, c.name, c.billing_address_id,
			c.credit_limit, c.current_balance, c.payment_terms_days,
			c.currency, c.sales_rep_id, c.default_route_id, c.default_warehouse_id,
			c.tax_exempt, c.is_active, c.created_by, c.created_at, c.updated_at,
			COALESCE(e.english_name, '') as sales_rep_name,
			COALESCE(w.name, '') as warehouse_name
		FROM customers c
		LEFT JOIN employees e ON c.sales_rep_id = e.id
		LEFT JOIN warehouses w ON c.default_warehouse_id = w.id
		WHERE c.id = $1`

	var result models.CustomerWithDetails
	var cust models.Customer

	err := s.db.QueryRow(ctx, query, id).Scan(
		&cust.ID,
		&cust.CustomerCode,
		&cust.Name,
		&cust.BillingAddressID,
		&cust.CreditLimit,
		&cust.CurrentBalance,
		&cust.PaymentTermsDays,
		&cust.Currency,
		&cust.SalesRepID,
		&cust.DefaultRouteID,
		&cust.DefaultWarehouseID,
		&cust.TaxExempt,
		&cust.IsActive,
		&cust.CreatedBy,
		&cust.CreatedAt,
		&cust.UpdatedAt,
		&result.SalesRepName,
		&result.WarehouseName,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	result.Customer = cust

	// Get ship-to addresses
	shipToQuery := `
		SELECT id, customer_id, ship_to_code, name, address_line1, address_line2,
			city, state, postal_code, country, phone, is_default, warehouse_id, route_id
		FROM customer_ship_to
		WHERE customer_id = $1
		ORDER BY is_default DESC, name`

	rows := s.db.Query(ctx, shipToQuery, id)
	defer rows.Close()

	for rows.Next() {
		var st models.CustomerShipTo
		err := rows.Scan(
			&st.ID, &st.CustomerID, &st.ShipToCode, &st.Name,
			&st.AddressLine1, &st.AddressLine2, &st.City, &st.State,
			&st.PostalCode, &st.Country, &st.Phone, &st.IsDefault,
			&st.WarehouseID, &st.RouteID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ship-to: %w", err)
		}
		result.ShipToList = append(result.ShipToList, st)
	}

	return &result, nil
}

func (s *customerServiceImpl) GetByCode(ctx context.Context, code string) (*models.Customer, error) {
	query := `
		SELECT id, customer_code, name, billing_address_id, credit_limit, 
			current_balance, payment_terms_days, currency, sales_rep_id,
			default_route_id, default_warehouse_id, tax_exempt, is_active,
			created_by, created_at, updated_at
		FROM customers
		WHERE customer_code = $1`

	var cust models.Customer
	err := s.db.QueryRow(ctx, query, code).Scan(
		&cust.ID, &cust.CustomerCode, &cust.Name, &cust.BillingAddressID,
		&cust.CreditLimit, &cust.CurrentBalance, &cust.PaymentTermsDays,
		&cust.Currency, &cust.SalesRepID, &cust.DefaultRouteID,
		&cust.DefaultWarehouseID, &cust.TaxExempt, &cust.IsActive,
		&cust.CreatedBy, &cust.CreatedAt, &cust.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &cust, nil
}

func (s *customerServiceImpl) Update(ctx context.Context, id int, req models.UpdateCustomerRequest) error {
	query := `
		UPDATE customers SET
			name = COALESCE($2, name),
			credit_limit = COALESCE($3, credit_limit),
			payment_terms_days = COALESCE($4, payment_terms_days),
			currency = COALESCE($5, currency),
			sales_rep_id = COALESCE($6, sales_rep_id),
			default_warehouse_id = COALESCE($7, default_warehouse_id),
			tax_exempt = COALESCE($8, tax_exempt),
			is_active = COALESCE($9, is_active),
			updated_at = NOW()
		WHERE id = $1`

	result, err := s.db.Exec(ctx, query,
		id,
		req.Name,
		req.CreditLimit,
		req.PaymentTermsDays,
		req.Currency,
		req.SalesRepID,
		req.DefaultWarehouseID,
		req.TaxExempt,
		req.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *customerServiceImpl) Delete(ctx context.Context, id int) error {
	// Soft delete - just mark as inactive
	query := `UPDATE customers SET is_active = false, updated_at = NOW() WHERE id = $1`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *customerServiceImpl) List(ctx context.Context, filters models.CustomerListFilters) ([]models.Customer, int64, error) {
	// Set defaults
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM customers WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters.Search != "" {
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR customer_code ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}
	if filters.SalesRepID != nil {
		countQuery += fmt.Sprintf(" AND sales_rep_id = $%d", argIndex)
		args = append(args, *filters.SalesRepID)
		argIndex++
	}
	if filters.IsActive != nil {
		countQuery += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *filters.IsActive)
		argIndex++
	}

	var total int64
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	// Get page
	query := `
		SELECT id, customer_code, name, billing_address_id, credit_limit,
			current_balance, payment_terms_days, currency, sales_rep_id,
			default_route_id, default_warehouse_id, tax_exempt, is_active,
			created_by, created_at, updated_at
		FROM customers WHERE 1=1`

	args = []interface{}{}
	argIndex = 1

	if filters.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR customer_code ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}
	if filters.SalesRepID != nil {
		query += fmt.Sprintf(" AND sales_rep_id = $%d", argIndex)
		args = append(args, *filters.SalesRepID)
		argIndex++
	}
	if filters.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *filters.IsActive)
		argIndex++
	}

	query += " ORDER BY name"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.PageSize, (filters.Page-1)*filters.PageSize)

	rows := s.db.Query(ctx, query, args...)
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(
			&c.ID, &c.CustomerCode, &c.Name, &c.BillingAddressID,
			&c.CreditLimit, &c.CurrentBalance, &c.PaymentTermsDays,
			&c.Currency, &c.SalesRepID, &c.DefaultRouteID,
			&c.DefaultWarehouseID, &c.TaxExempt, &c.IsActive,
			&c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, c)
	}

	return customers, total, nil
}

func (s *customerServiceImpl) GetOrderGuide(ctx context.Context, customerID int) ([]models.CustomerOrderGuide, error) {
	query := `
		SELECT 
			cog.id, cog.customer_id, cog.product_id,
			p.name as product_name, p.sku as product_sku,
			cog.default_quantity, cog.last_ordered_date, cog.last_ordered_quantity,
			cog.avg_weekly_quantity, cog.times_ordered, cog.is_push_item, cog.custom_price
		FROM customer_order_guides cog
		JOIN products p ON cog.product_id = p.id
		WHERE cog.customer_id = $1
		ORDER BY cog.times_ordered DESC, p.name`

	rows := s.db.Query(ctx, query, customerID)
	defer rows.Close()

	var guides []models.CustomerOrderGuide
	for rows.Next() {
		var g models.CustomerOrderGuide
		err := rows.Scan(
			&g.ID, &g.CustomerID, &g.ProductID,
			&g.ProductName, &g.ProductSKU,
			&g.DefaultQuantity, &g.LastOrderedDate, &g.LastOrderedQuantity,
			&g.AvgWeeklyQuantity, &g.TimesOrdered, &g.IsPushItem, &g.CustomPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order guide: %w", err)
		}
		guides = append(guides, g)
	}

	return guides, nil
}

func (s *customerServiceImpl) AddShipTo(ctx context.Context, customerID int, shipTo models.CustomerShipTo) (int, error) {
	query := `
		INSERT INTO customer_ship_to (
			customer_id, ship_to_code, name, address_line1, address_line2,
			city, state, postal_code, country, phone, is_default, warehouse_id, route_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query,
		customerID,
		shipTo.ShipToCode,
		shipTo.Name,
		shipTo.AddressLine1,
		shipTo.AddressLine2,
		shipTo.City,
		shipTo.State,
		shipTo.PostalCode,
		shipTo.Country,
		shipTo.Phone,
		shipTo.IsDefault,
		shipTo.WarehouseID,
		shipTo.RouteID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to add ship-to: %w", err)
	}

	return id, nil
}
