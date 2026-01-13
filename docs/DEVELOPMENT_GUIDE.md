# 🏭 ERP System - Development Guide

## Complete Reference for Continuing Development

---

## 📋 Table of Contents

1. [Project Overview](#project-overview)
2. [What's Already Created](#whats-already-created)
3. [What Needs Implementation](#what-needs-implementation)
4. [Setup Instructions](#setup-instructions)
5. [Module Templates](#module-templates)
6. [Database Schema](#database-schema)
7. [API Endpoints Reference](#api-endpoints-reference)
8. [Development Workflow](#development-workflow)

---

## Project Overview

### Tech Stack
| Component | Technology |
|-----------|------------|
| Language | Go 1.23 |
| Router | Chi v5 |
| Database | PostgreSQL (pgx/v5) |
| Storage | MinIO |
| Auth | JWT |
| Docs | Swagger |

### Project Location
```
C:\Users\AnasB\Desktop\ERP-System\
```

### Module Path
```
github.com/anas-dev-92/FoodHive
```

---

## What's Already Created

### ✅ Core Infrastructure (100% Complete)
```
core/
├── auth/auth.go           ✅ Authentication service
├── jwt/jwt.go             ✅ JWT token service
├── postgres/db.go         ✅ Database connection & transactions
├── storage/storage.go     ✅ MinIO file storage
├── utils/env/env.go       ✅ Environment config
├── utils/helpers/         ✅ JSON helpers, pagination
└── go.mod                 ✅ Module definition
```

### ✅ Application Structure (100% Complete)
```
registration/
├── main.go                ✅ Entry point with all middleware imports
├── go.mod                 ✅ Dependencies
└── src/v1/
    ├── v1.go              ✅ Router with all routes defined
    ├── middlewares/       ✅ 16 module folders created
    ├── routes/            ✅ 17 module folders created
    ├── services/          ✅ 16 module folders created
    └── models/            ✅ Core models created
```

### ✅ Models Created
| File | Status | Contents |
|------|--------|----------|
| `models/common.go` | ✅ | CustomDate, Validator, Pagination, Responses |
| `models/customer.go` | ✅ | Customer, ShipTo, OrderGuide, Requests |
| `models/product.go` | ✅ | Product, Category, Units, Requests |
| `models/inventory.go` | ✅ | Inventory, Transactions, Lot, Requests |
| `models/sales_order.go` | ✅ | SalesOrder, Lines, LostSales, Requests |

### ✅ Services Created
| Service | Status | Location |
|---------|--------|----------|
| Customer Service | ✅ Complete | `services/customer/customer.go` |

### ✅ Middlewares Created
| Middleware | Status | Location |
|------------|--------|----------|
| Auth | ✅ Complete | `middlewares/auth/auth.go` |
| Customer | ✅ Complete | `middlewares/customer/customer.go` |

### ✅ SQL Schema
| File | Status | Contents |
|------|--------|----------|
| `sql/001_core_tables.sql` | ✅ | Enums, Users, Products, Customers, Vendors, Inventory |
| `sql/002_transactions_tables.sql` | ✅ | Sales Orders, Purchase Orders, Picking |

---

## What Needs Implementation

### 🔴 Priority 1: Complete Middleware Stubs

Each module needs a middleware file. Use this template:

```go
// File: middlewares/{module}/{module}.go
package {module}

import (
    "context"
    "net/http"

    "github.com/anas-dev-92/FoodHive/core/postgres"
    "github.com/anas-dev-92/FoodHive/registration/src/v1/services/{module}"
)

type contextKey string

const {module}Key = contextKey("{module}_service")

func New(db postgres.Connection) func(http.Handler) http.Handler {
    svc := {module}.New(db)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), {module}Key, svc)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func Instance(ctx context.Context) ({module}.{Module}Service, bool) {
    svc, ok := ctx.Value({module}Key).({module}.{Module}Service)
    return svc, ok
}
```

**Modules needing middleware:**
- [ ] `employee`
- [ ] `warehouse`
- [ ] `product`
- [ ] `vendor`
- [ ] `inventory`
- [ ] `sales_order`
- [ ] `purchase_order`
- [ ] `ar`
- [ ] `ap`
- [ ] `gl`
- [ ] `pricing`
- [ ] `picking`
- [ ] `wms`
- [ ] `bank`
- [ ] `payroll`

### 🔴 Priority 2: Complete Route Files

Each module needs a router. Use this template:

```go
// File: routes/{module}/router.go
package {module}

import (
    "github.com/go-chi/chi/v5"
    "github.com/anas-dev-92/FoodHive/core/auth"
    "github.com/anas-dev-92/FoodHive/core/jwt"
    "github.com/anas-dev-92/FoodHive/core/postgres"
    authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
    app := chi.NewRouter()

    // Apply authentication middleware
    app.Use(authMiddleware.Authenticate(jwtService))

    // Routes with authorization
    app.With(authMiddleware.Authorize(jwtService)).Post("/create", HandlerCreate())
    app.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", HandlerGet())
    app.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", HandlerUpdate())
    app.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", HandlerDelete())
    app.With(authMiddleware.Authorize(jwtService)).Get("/list", HandlerList())

    return app
}
```

### 🔴 Priority 3: Complete Service Files

Each module needs a service. Use this template:

```go
// File: services/{module}/{module}.go
package {module}

import (
    "context"
    "errors"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/anas-dev-92/FoodHive/core/postgres"
    "github.com/anas-dev-92/FoodHive/registration/src/v1/models"
)

var (
    ErrNotFound = errors.New("{module} not found")
)

type {Module}Service interface {
    Create(ctx context.Context, req models.Create{Module}Request, createdBy int) (int, error)
    GetByID(ctx context.Context, id int) (*models.{Module}, error)
    Update(ctx context.Context, id int, req models.Update{Module}Request) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, filters models.{Module}ListFilters) ([]models.{Module}, int64, error)
}

type {module}ServiceImpl struct {
    db postgres.Connection
}

func New(db postgres.Connection) {Module}Service {
    return &{module}ServiceImpl{db: db}
}

// Implement interface methods...
```

---

## Setup Instructions

### 1. Database Setup

```powershell
# Create database user and database
psql -U postgres -c "CREATE USER erp PASSWORD 'erp_password' SUPERUSER"
psql -U postgres -c "CREATE DATABASE erp_db OWNER erp"

# Run migrations
cd C:\Users\AnasB\Desktop\ERP-System
Get-Content .\sql\001_core_tables.sql | psql -U erp -d erp_db
Get-Content .\sql\002_transactions_tables.sql | psql -U erp -d erp_db
```

### 2. Environment Setup

```powershell
# Create .env file
cd C:\Users\AnasB\Desktop\ERP-System
Copy-Item env.example.txt .env

# Edit .env with your settings
notepad .env
```

**.env contents:**
```env
DB_CONFIG=postgres://erp:erp_password@localhost:5432/erp_db?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-in-production
STORAGE_HOST=localhost:9000
STORAGE_KEY=minioadmin
STORAGE_SECRET=minioadmin
STORAGE_SSL=false
```

### 3. Run the Application

```powershell
cd C:\Users\AnasB\Desktop\ERP-System\registration

# Install dependencies
go mod tidy

# Run
go run main.go
```

---

## Module Templates

### Complete Middleware Template (Copy for each module)

```go
// middlewares/product/product.go
package product

import (
    "context"
    "net/http"

    "github.com/anas-dev-92/FoodHive/core/postgres"
    "github.com/anas-dev-92/FoodHive/registration/src/v1/services/product"
)

type contextKey string

const productKey = contextKey("product_service")

func New(db postgres.Connection) func(http.Handler) http.Handler {
    productService := product.New(db)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), productKey, productService)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func Instance(ctx context.Context) (product.ProductService, bool) {
    svc, ok := ctx.Value(productKey).(product.ProductService)
    return svc, ok
}
```

### Complete Router Template

```go
// routes/product/router.go
package product

import (
    "github.com/go-chi/chi/v5"
    "github.com/anas-dev-92/FoodHive/core/auth"
    "github.com/anas-dev-92/FoodHive/core/jwt"
    "github.com/anas-dev-92/FoodHive/core/postgres"
    authMiddleware "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/auth"
)

func Router(db postgres.Executor, jwtService jwt.JWTService, authService auth.AuthService) chi.Router {
    app := chi.NewRouter()

    app.Use(authMiddleware.Authenticate(jwtService))

    app.With(authMiddleware.Authorize(jwtService)).Post("/create", HandlerCreate())
    app.With(authMiddleware.Authorize(jwtService)).Get("/get/{id}", HandlerGet())
    app.With(authMiddleware.Authorize(jwtService)).Put("/update/{id}", HandlerUpdate())
    app.With(authMiddleware.Authorize(jwtService)).Delete("/delete/{id}", HandlerDelete())
    app.With(authMiddleware.Authorize(jwtService)).Get("/list", HandlerList())

    return app
}

// handlers.go
package product

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    middlewareProduct "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/product"
    "github.com/anas-dev-92/FoodHive/registration/src/v1/models"
)

func HandlerCreate() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req models.CreateProductRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Validate
        v := models.NewValidator()
        models.ValidateProduct(v, &req)
        if !v.Valid() {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Validation failed", Details: v.Errors})
            return
        }

        svc, ok := middlewareProduct.Instance(r.Context())
        if !ok {
            http.Error(w, "Service unavailable", http.StatusInternalServerError)
            return
        }

        // Get user ID from context (set by auth middleware)
        // userID := r.Context().Value(authMiddleware.UserIDKey).(int)

        id, err := svc.Create(r.Context(), req, 1) // Replace 1 with actual userID
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(models.IDResponse{ID: id, Message: "Created successfully"})
    }
}

func HandlerGet() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(chi.URLParam(r, "id"))
        if err != nil {
            http.Error(w, "Invalid ID", http.StatusBadRequest)
            return
        }

        svc, ok := middlewareProduct.Instance(r.Context())
        if !ok {
            http.Error(w, "Service unavailable", http.StatusInternalServerError)
            return
        }

        result, err := svc.GetByID(r.Context(), id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    }
}

func HandlerUpdate() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(chi.URLParam(r, "id"))
        if err != nil {
            http.Error(w, "Invalid ID", http.StatusBadRequest)
            return
        }

        var req models.UpdateProductRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        svc, ok := middlewareProduct.Instance(r.Context())
        if !ok {
            http.Error(w, "Service unavailable", http.StatusInternalServerError)
            return
        }

        if err := svc.Update(r.Context(), id, req); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(models.SuccessResponse{Message: "Updated successfully"})
    }
}

func HandlerDelete() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(chi.URLParam(r, "id"))
        if err != nil {
            http.Error(w, "Invalid ID", http.StatusBadRequest)
            return
        }

        svc, ok := middlewareProduct.Instance(r.Context())
        if !ok {
            http.Error(w, "Service unavailable", http.StatusInternalServerError)
            return
        }

        if err := svc.Delete(r.Context(), id); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func HandlerList() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse query parameters for filtering/pagination
        page, _ := strconv.Atoi(r.URL.Query().Get("page"))
        pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
        if page < 1 {
            page = 1
        }
        if pageSize < 1 {
            pageSize = 20
        }

        filters := models.ProductListFilters{
            Search:   r.URL.Query().Get("search"),
            Page:     page,
            PageSize: pageSize,
        }

        svc, ok := middlewareProduct.Instance(r.Context())
        if !ok {
            http.Error(w, "Service unavailable", http.StatusInternalServerError)
            return
        }

        items, total, err := svc.List(r.Context(), filters)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        totalPages := int(total) / pageSize
        if int(total)%pageSize > 0 {
            totalPages++
        }

        response := models.PaginatedResponse{
            Data: items,
            Pagination: models.Pagination{
                Page:       page,
                PageSize:   pageSize,
                TotalItems: total,
                TotalPages: totalPages,
            },
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}
```

### Complete Service Template

```go
// services/product/product.go
package product

import (
    "context"
    "errors"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/anas-dev-92/FoodHive/core/postgres"
    "github.com/anas-dev-92/FoodHive/registration/src/v1/models"
)

var (
    ErrNotFound     = errors.New("product not found")
    ErrDuplicateSKU = errors.New("SKU already exists")
)

type ProductService interface {
    Create(ctx context.Context, req models.CreateProductRequest, createdBy int) (int, error)
    GetByID(ctx context.Context, id int) (*models.ProductWithDetails, error)
    GetBySKU(ctx context.Context, sku string) (*models.Product, error)
    Update(ctx context.Context, id int, req models.UpdateProductRequest) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, filters models.ProductListFilters) ([]models.Product, int64, error)
    AddUnit(ctx context.Context, productID int, unit models.ProductUnit) (int, error)
}

type productServiceImpl struct {
    db postgres.Connection
}

func New(db postgres.Connection) ProductService {
    return &productServiceImpl{db: db}
}

func (s *productServiceImpl) Create(ctx context.Context, req models.CreateProductRequest, createdBy int) (int, error) {
    query := `
        INSERT INTO products (
            sku, barcode, upc, name, description, category_id,
            base_unit, is_catch_weight, catch_weight_unit, country_of_origin,
            shelf_life_days, min_shelf_life_days, is_lot_tracked, is_serialized,
            haccp_category, qc_required
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
        RETURNING id`

    var id int
    err := s.db.QueryRow(ctx, query,
        req.SKU, req.Barcode, req.UPC, req.Name, req.Description, req.CategoryID,
        req.BaseUnit, req.IsCatchWeight, req.CatchWeightUnit, req.CountryOfOrigin,
        req.ShelfLifeDays, req.MinShelfLifeDays, req.IsLotTracked, req.IsSerialized,
        req.HACCPCategory, req.QCRequired,
    ).Scan(&id)

    if err != nil {
        return 0, fmt.Errorf("failed to create product: %w", err)
    }

    return id, nil
}

func (s *productServiceImpl) GetByID(ctx context.Context, id int) (*models.ProductWithDetails, error) {
    query := `
        SELECT 
            p.id, p.sku, p.barcode, p.upc, p.name, p.description, p.category_id,
            p.base_unit, p.is_catch_weight, p.catch_weight_unit, p.country_of_origin,
            p.shelf_life_days, p.min_shelf_life_days, p.is_lot_tracked, p.is_serialized,
            p.haccp_category, p.qc_required, p.is_active, p.created_at, p.updated_at,
            COALESCE(SUM(i.quantity_on_hand), 0) as current_stock,
            COALESCE(AVG(i.average_cost), 0) as avg_cost
        FROM products p
        LEFT JOIN inventory i ON p.id = i.product_id
        WHERE p.id = $1
        GROUP BY p.id`

    var result models.ProductWithDetails
    var prod models.Product

    err := s.db.QueryRow(ctx, query, id).Scan(
        &prod.ID, &prod.SKU, &prod.Barcode, &prod.UPC, &prod.Name, &prod.Description,
        &prod.CategoryID, &prod.BaseUnit, &prod.IsCatchWeight, &prod.CatchWeightUnit,
        &prod.CountryOfOrigin, &prod.ShelfLifeDays, &prod.MinShelfLifeDays,
        &prod.IsLotTracked, &prod.IsSerialized, &prod.HACCPCategory, &prod.QCRequired,
        &prod.IsActive, &prod.CreatedAt, &prod.UpdatedAt,
        &result.CurrentStock, &result.AvgCost,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("failed to get product: %w", err)
    }

    result.Product = prod

    // Get units
    unitsQuery := `
        SELECT id, product_id, unit_name, description, conversion_factor,
            barcode, weight, is_purchase_unit, is_sales_unit
        FROM product_units
        WHERE product_id = $1`

    rows := s.db.Query(ctx, unitsQuery, id)
    defer rows.Close()

    for rows.Next() {
        var unit models.ProductUnit
        err := rows.Scan(
            &unit.ID, &unit.ProductID, &unit.UnitName, &unit.Description,
            &unit.ConversionFactor, &unit.Barcode, &unit.Weight,
            &unit.IsPurchaseUnit, &unit.IsSalesUnit,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan unit: %w", err)
        }
        result.Units = append(result.Units, unit)
    }

    return &result, nil
}

func (s *productServiceImpl) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
    query := `
        SELECT id, sku, barcode, upc, name, description, category_id,
            base_unit, is_catch_weight, catch_weight_unit, country_of_origin,
            shelf_life_days, min_shelf_life_days, is_lot_tracked, is_serialized,
            haccp_category, qc_required, is_active, created_at, updated_at
        FROM products
        WHERE sku = $1`

    var prod models.Product
    err := s.db.QueryRow(ctx, query, sku).Scan(
        &prod.ID, &prod.SKU, &prod.Barcode, &prod.UPC, &prod.Name, &prod.Description,
        &prod.CategoryID, &prod.BaseUnit, &prod.IsCatchWeight, &prod.CatchWeightUnit,
        &prod.CountryOfOrigin, &prod.ShelfLifeDays, &prod.MinShelfLifeDays,
        &prod.IsLotTracked, &prod.IsSerialized, &prod.HACCPCategory, &prod.QCRequired,
        &prod.IsActive, &prod.CreatedAt, &prod.UpdatedAt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("failed to get product: %w", err)
    }

    return &prod, nil
}

func (s *productServiceImpl) Update(ctx context.Context, id int, req models.UpdateProductRequest) error {
    query := `
        UPDATE products SET
            name = COALESCE($2, name),
            description = COALESCE($3, description),
            category_id = COALESCE($4, category_id),
            country_of_origin = COALESCE($5, country_of_origin),
            shelf_life_days = COALESCE($6, shelf_life_days),
            min_shelf_life_days = COALESCE($7, min_shelf_life_days),
            haccp_category = COALESCE($8, haccp_category),
            qc_required = COALESCE($9, qc_required),
            is_active = COALESCE($10, is_active),
            updated_at = NOW()
        WHERE id = $1`

    result, err := s.db.Exec(ctx, query,
        id, req.Name, req.Description, req.CategoryID, req.CountryOfOrigin,
        req.ShelfLifeDays, req.MinShelfLifeDays, req.HACCPCategory,
        req.QCRequired, req.IsActive,
    )

    if err != nil {
        return fmt.Errorf("failed to update product: %w", err)
    }

    if result.RowsAffected() == 0 {
        return ErrNotFound
    }

    return nil
}

func (s *productServiceImpl) Delete(ctx context.Context, id int) error {
    query := `UPDATE products SET is_active = false, updated_at = NOW() WHERE id = $1`
    result, err := s.db.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete product: %w", err)
    }

    if result.RowsAffected() == 0 {
        return ErrNotFound
    }

    return nil
}

func (s *productServiceImpl) List(ctx context.Context, filters models.ProductListFilters) ([]models.Product, int64, error) {
    if filters.Page < 1 {
        filters.Page = 1
    }
    if filters.PageSize < 1 {
        filters.PageSize = 20
    }

    // Count
    countQuery := `SELECT COUNT(*) FROM products WHERE 1=1`
    args := []interface{}{}
    argIndex := 1

    if filters.Search != "" {
        countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argIndex, argIndex)
        args = append(args, "%"+filters.Search+"%")
        argIndex++
    }
    if filters.CategoryID != nil {
        countQuery += fmt.Sprintf(" AND category_id = $%d", argIndex)
        args = append(args, *filters.CategoryID)
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
        return nil, 0, fmt.Errorf("failed to count: %w", err)
    }

    // Get page
    query := `
        SELECT id, sku, barcode, upc, name, description, category_id,
            base_unit, is_catch_weight, catch_weight_unit, country_of_origin,
            shelf_life_days, min_shelf_life_days, is_lot_tracked, is_serialized,
            haccp_category, qc_required, is_active, created_at, updated_at
        FROM products WHERE 1=1`

    args = []interface{}{}
    argIndex = 1

    if filters.Search != "" {
        query += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d)", argIndex, argIndex)
        args = append(args, "%"+filters.Search+"%")
        argIndex++
    }
    if filters.CategoryID != nil {
        query += fmt.Sprintf(" AND category_id = $%d", argIndex)
        args = append(args, *filters.CategoryID)
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

    var products []models.Product
    for rows.Next() {
        var p models.Product
        err := rows.Scan(
            &p.ID, &p.SKU, &p.Barcode, &p.UPC, &p.Name, &p.Description,
            &p.CategoryID, &p.BaseUnit, &p.IsCatchWeight, &p.CatchWeightUnit,
            &p.CountryOfOrigin, &p.ShelfLifeDays, &p.MinShelfLifeDays,
            &p.IsLotTracked, &p.IsSerialized, &p.HACCPCategory, &p.QCRequired,
            &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan: %w", err)
        }
        products = append(products, p)
    }

    return products, total, nil
}

func (s *productServiceImpl) AddUnit(ctx context.Context, productID int, unit models.ProductUnit) (int, error) {
    query := `
        INSERT INTO product_units (
            product_id, unit_name, description, conversion_factor,
            barcode, weight, is_purchase_unit, is_sales_unit
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`

    var id int
    err := s.db.QueryRow(ctx, query,
        productID, unit.UnitName, unit.Description, unit.ConversionFactor,
        unit.Barcode, unit.Weight, unit.IsPurchaseUnit, unit.IsSalesUnit,
    ).Scan(&id)

    if err != nil {
        return 0, fmt.Errorf("failed to add unit: %w", err)
    }

    return id, nil
}
```

---

## Database Schema

### Additional Tables Needed

Create `sql/003_financial_tables.sql`:

```sql
-- =====================================================
-- ACCOUNTS RECEIVABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS ar_invoices (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(20) UNIQUE NOT NULL,
    customer_id INT REFERENCES customers(id),
    ship_to_id INT REFERENCES customer_ship_to(id),
    sales_order_id INT REFERENCES sales_orders(id),
    invoice_date DATE DEFAULT CURRENT_DATE,
    due_date DATE,
    subtotal DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    freight_amount DECIMAL(12,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    paid_amount DECIMAL(12,2) DEFAULT 0,
    balance DECIMAL(12,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    status ar_invoice_status DEFAULT 'DRAFT',
    notes TEXT,
    gl_posted BOOLEAN DEFAULT false,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ar_payments (
    id SERIAL PRIMARY KEY,
    payment_number VARCHAR(20) UNIQUE NOT NULL,
    customer_id INT REFERENCES customers(id),
    payment_date DATE DEFAULT CURRENT_DATE,
    amount DECIMAL(12,2),
    currency VARCHAR(3) DEFAULT 'USD',
    payment_method VARCHAR(20),
    reference_number VARCHAR(50),
    check_number VARCHAR(50),
    bank_account_id INT,
    notes TEXT,
    gl_posted BOOLEAN DEFAULT false,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ar_payment_allocations (
    id SERIAL PRIMARY KEY,
    payment_id INT REFERENCES ar_payments(id) ON DELETE CASCADE,
    invoice_id INT REFERENCES ar_invoices(id),
    amount_applied DECIMAL(12,2),
    discount_taken DECIMAL(12,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- ACCOUNTS PAYABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS ap_invoices (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) NOT NULL,
    vendor_id INT REFERENCES vendors(id),
    po_id INT REFERENCES purchase_orders(id),
    invoice_date DATE,
    received_date DATE DEFAULT CURRENT_DATE,
    due_date DATE,
    subtotal DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    freight_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    paid_amount DECIMAL(12,2) DEFAULT 0,
    balance DECIMAL(12,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    status ap_invoice_status DEFAULT 'PENDING',
    notes TEXT,
    gl_posted BOOLEAN DEFAULT false,
    approved_by INT REFERENCES employees(id),
    approved_at TIMESTAMP,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(vendor_id, invoice_number)
);

CREATE TABLE IF NOT EXISTS ap_payments (
    id SERIAL PRIMARY KEY,
    payment_number VARCHAR(20) UNIQUE NOT NULL,
    vendor_id INT REFERENCES vendors(id),
    payment_date DATE DEFAULT CURRENT_DATE,
    amount DECIMAL(12,2),
    currency VARCHAR(3) DEFAULT 'USD',
    payment_method VARCHAR(20),
    check_number VARCHAR(50),
    bank_account_id INT,
    notes TEXT,
    gl_posted BOOLEAN DEFAULT false,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- GENERAL LEDGER
-- =====================================================

CREATE TABLE IF NOT EXISTS gl_fiscal_years (
    id SERIAL PRIMARY KEY,
    year_name VARCHAR(20) UNIQUE NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_closed BOOLEAN DEFAULT false,
    closed_by INT REFERENCES employees(id),
    closed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gl_periods (
    id SERIAL PRIMARY KEY,
    fiscal_year_id INT REFERENCES gl_fiscal_years(id),
    period_number INT,
    period_name VARCHAR(20),
    start_date DATE,
    end_date DATE,
    is_closed BOOLEAN DEFAULT false,
    UNIQUE(fiscal_year_id, period_number)
);

CREATE TABLE IF NOT EXISTS gl_accounts (
    id SERIAL PRIMARY KEY,
    account_number VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    account_type gl_account_type NOT NULL,
    parent_id INT REFERENCES gl_accounts(id),
    is_header BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    normal_balance VARCHAR(10) DEFAULT 'DEBIT',
    currency VARCHAR(3) DEFAULT 'USD',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gl_transactions (
    id SERIAL PRIMARY KEY,
    transaction_date DATE NOT NULL,
    period_id INT REFERENCES gl_periods(id),
    account_id INT REFERENCES gl_accounts(id),
    debit DECIMAL(14,2) DEFAULT 0,
    credit DECIMAL(14,2) DEFAULT 0,
    description TEXT,
    reference_type VARCHAR(30),
    reference_id INT,
    reference_number VARCHAR(50),
    is_posted BOOLEAN DEFAULT false,
    posted_by INT REFERENCES employees(id),
    posted_at TIMESTAMP,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- BANK
-- =====================================================

CREATE TABLE IF NOT EXISTS bank_accounts (
    id SERIAL PRIMARY KEY,
    account_code VARCHAR(20) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    bank_name TEXT,
    account_number VARCHAR(50),
    routing_number VARCHAR(50),
    currency VARCHAR(3) DEFAULT 'USD',
    current_balance DECIMAL(14,2) DEFAULT 0,
    gl_account_id INT REFERENCES gl_accounts(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bank_transactions (
    id SERIAL PRIMARY KEY,
    bank_account_id INT REFERENCES bank_accounts(id),
    transaction_date DATE,
    transaction_type bank_transaction_type,
    amount DECIMAL(12,2),
    reference_number VARCHAR(50),
    description TEXT,
    is_reconciled BOOLEAN DEFAULT false,
    reconciled_date DATE,
    source_type VARCHAR(20),
    source_id INT,
    created_by INT REFERENCES employees(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- INDEXES
CREATE INDEX idx_ar_invoices_customer ON ar_invoices(customer_id);
CREATE INDEX idx_ar_invoices_status ON ar_invoices(status);
CREATE INDEX idx_ap_invoices_vendor ON ap_invoices(vendor_id);
CREATE INDEX idx_ap_invoices_status ON ap_invoices(status);
CREATE INDEX idx_gl_transactions_date ON gl_transactions(transaction_date);
CREATE INDEX idx_gl_transactions_account ON gl_transactions(account_id);
```

---

## API Endpoints Reference

### Planned Endpoints

| Module | Method | Endpoint | Description |
|--------|--------|----------|-------------|
| **Auth** | POST | `/v1/login` | Login |
| **Employees** | POST | `/v1/employees/create` | Create employee |
| | GET | `/v1/employees/get/{id}` | Get by ID |
| | PUT | `/v1/employees/update/{id}` | Update |
| | DELETE | `/v1/employees/delete/{id}` | Delete |
| | GET | `/v1/employees/list` | List all |
| **Warehouses** | POST | `/v1/warehouses/create` | Create |
| | GET | `/v1/warehouses/get/{id}` | Get by ID |
| | PUT | `/v1/warehouses/update/{id}` | Update |
| | DELETE | `/v1/warehouses/delete/{id}` | Delete |
| | GET | `/v1/warehouses/list` | List all |
| **Products** | POST | `/v1/products/create` | Create |
| | GET | `/v1/products/get/{id}` | Get by ID |
| | GET | `/v1/products/sku/{sku}` | Get by SKU |
| | PUT | `/v1/products/update/{id}` | Update |
| | DELETE | `/v1/products/delete/{id}` | Delete |
| | GET | `/v1/products/list` | List all |
| | POST | `/v1/products/{id}/units` | Add unit |
| **Customers** | POST | `/v1/customers/create` | Create |
| | GET | `/v1/customers/get/{id}` | Get by ID |
| | PUT | `/v1/customers/update/{id}` | Update |
| | DELETE | `/v1/customers/delete/{id}` | Delete |
| | GET | `/v1/customers/list` | List all |
| | GET | `/v1/customers/{id}/order-guide` | Get order guide |
| | POST | `/v1/customers/{id}/ship-to` | Add ship-to |
| **Vendors** | POST | `/v1/vendors/create` | Create |
| | GET | `/v1/vendors/get/{id}` | Get by ID |
| | PUT | `/v1/vendors/update/{id}` | Update |
| | DELETE | `/v1/vendors/delete/{id}` | Delete |
| | GET | `/v1/vendors/list` | List all |
| **Inventory** | GET | `/v1/inventory/get/{id}` | Get by ID |
| | GET | `/v1/inventory/product/{productId}` | Get by product |
| | POST | `/v1/inventory/adjust` | Adjust |
| | POST | `/v1/inventory/transfer` | Transfer |
| | GET | `/v1/inventory/list` | List all |
| **Sales Orders** | POST | `/v1/sales-orders/create` | Create |
| | GET | `/v1/sales-orders/get/{id}` | Get by ID |
| | PUT | `/v1/sales-orders/update/{id}` | Update |
| | DELETE | `/v1/sales-orders/delete/{id}` | Cancel |
| | GET | `/v1/sales-orders/list` | List all |
| | POST | `/v1/sales-orders/{id}/confirm` | Confirm |
| | POST | `/v1/sales-orders/{id}/ship` | Ship |
| **Purchase Orders** | POST | `/v1/purchase-orders/create` | Create |
| | GET | `/v1/purchase-orders/get/{id}` | Get by ID |
| | PUT | `/v1/purchase-orders/update/{id}` | Update |
| | DELETE | `/v1/purchase-orders/delete/{id}` | Cancel |
| | GET | `/v1/purchase-orders/list` | List all |
| | POST | `/v1/purchase-orders/{id}/receive` | Receive |
| **AR** | POST | `/v1/ar/invoices/create` | Create invoice |
| | GET | `/v1/ar/invoices/get/{id}` | Get by ID |
| | POST | `/v1/ar/payments/create` | Create payment |
| | GET | `/v1/ar/customers/{id}/aging` | Aging report |
| **AP** | POST | `/v1/ap/invoices/create` | Create invoice |
| | GET | `/v1/ap/invoices/get/{id}` | Get by ID |
| | POST | `/v1/ap/payments/create` | Create payment |
| | GET | `/v1/ap/vendors/{id}/aging` | Aging report |
| **GL** | POST | `/v1/gl/accounts/create` | Create account |
| | GET | `/v1/gl/accounts/tree` | Chart of accounts |
| | POST | `/v1/gl/journal-entries` | Journal entry |
| | GET | `/v1/gl/trial-balance` | Trial balance |

---

## Development Workflow

### Step-by-Step for Each Module

1. **Create Middleware** (`middlewares/{module}/{module}.go`)
2. **Create Service Interface & Implementation** (`services/{module}/{module}.go`)
3. **Create Router** (`routes/{module}/router.go`)
4. **Create Handlers** (`routes/{module}/handlers.go`)
5. **Add Models if needed** (`models/{module}.go`)
6. **Test the endpoints**

### Quick Commands

```powershell
# Navigate to project
cd C:\Users\AnasB\Desktop\ERP-System\registration

# Run with live reload (if you have air installed)
air

# Or regular run
go run main.go

# Run tests
go test ./test/... -v

# Generate swagger docs
swag init -g main.go
```

---

## Checklist

### Phase 1: Foundation ⏳
- [ ] Employee middleware, service, routes
- [ ] Warehouse middleware, service, routes
- [ ] Product middleware, service, routes
- [ ] Vendor middleware, service, routes
- [ ] Inventory middleware, service, routes
- [x] Customer middleware, service, routes ✅

### Phase 2: Core ERP
- [ ] Sales Order middleware, service, routes
- [ ] Purchase Order middleware, service, routes
- [ ] AR middleware, service, routes
- [ ] AP middleware, service, routes

### Phase 3: Advanced
- [ ] GL middleware, service, routes
- [ ] Pricing middleware, service, routes
- [ ] Bank middleware, service, routes
- [ ] Payroll middleware, service, routes

### Phase 4: WMS
- [ ] Picking middleware, service, routes
- [ ] WMS middleware, service, routes

---

**Last Updated:** January 9, 2026  
**Project Location:** `C:\Users\AnasB\Desktop\ERP-System\`

