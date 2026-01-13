package v1

import (
	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/core/storage"
	"github.com/go-chi/chi/v5"

	// Routes - Currently implemented
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/ap"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/ar"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/bank"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/catch_weight"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/customer"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/department"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/employee"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/finance"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/gl"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/inventory"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/login"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/payroll"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/picking"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/pricing"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/product"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/purchase_order"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/role"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/sales_order"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/vendor"
	"github.com/anas-dev-92/FoodHive/registration/src/v1/routes/warehouse"
	// TODO: Uncomment as routes are implemented
	// "github.com/anas-dev-92/FoodHive/registration/src/v1/routes/wms"
)

// Suppress unused import warnings for storage
var _ = storage.StorageService(nil)

func Router(
	app chi.Router,
	jwtService jwt.JWTService,
	db postgres.Executor,
	storageService storage.StorageService,
	authService auth.AuthService,
) chi.Router {

	// ===========================================
	// Authentication (No auth required)
	// ===========================================
	app.Post("/login", login.Handler(jwtService, db))

	// ===========================================
	// Phase 1: Foundation - Master Data
	// ===========================================
	app.Mount("/employees", employee.Router(db, jwtService, authService))
	app.Mount("/departments", department.Router(db, jwtService, authService))
	app.Mount("/roles", role.Router(db, jwtService, authService))
	app.Mount("/customers", customer.Router(db, jwtService, authService))
	app.Mount("/vendors", vendor.Router(db, jwtService, authService))
	app.Mount("/warehouses", warehouse.Router(db, jwtService, authService))
	app.Mount("/inventory", inventory.Router(db, jwtService, authService))

	// ===========================================
	// Phase 2: Core ERP - Transactions
	// ===========================================
	app.Mount("/purchase-orders", purchase_order.Router(db, jwtService, authService))
	app.Mount("/ar", ar.Router(db, jwtService, authService))
	app.Mount("/ap", ap.Router(db, jwtService, authService))
	app.Mount("/sales-orders", sales_order.Router(db, jwtService, authService))

	// ===========================================
	// Phase 3: Advanced - Financial
	// ===========================================
	app.Mount("/gl", gl.Router(db, jwtService, authService))
	app.Mount("/pricing", pricing.Router(db, jwtService, authService))
	app.Mount("/bank", bank.Router(db, jwtService, authService))
	app.Mount("/payroll", payroll.Router(db, jwtService, authService))
	app.Mount("/finance", finance.Router(db, jwtService, authService))
	app.Mount("/products", product.Router(db, jwtService, authService))

	// ===========================================
	// Phase 4: WMS - Warehouse Operations
	// ===========================================
	app.Mount("/picking", picking.Router(db, jwtService, authService))
	app.Mount("/catch-weight", catch_weight.Router(db, jwtService, authService))

	// TODO: Uncomment as routes are implemented
	// app.Mount("/wms", wms.Router(db, jwtService, authService))

	return app
}
