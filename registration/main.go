package main

import (
	"log"
	"net/http"

	envconfig "github.com/Netflix/go-env"
	"github.com/joho/godotenv"

	// httpSwagger "github.com/swaggo/http-swagger/v2" // TODO: Enable after generating swagger docs

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/anas-dev-92/FoodHive/core/auth"
	"github.com/anas-dev-92/FoodHive/core/jwt"
	"github.com/anas-dev-92/FoodHive/core/postgres"
	"github.com/anas-dev-92/FoodHive/core/storage"
	"github.com/anas-dev-92/FoodHive/core/utils/env"

	// _ "github.com/anas-dev-92/FoodHive/registration/docs" // TODO: Enable after generating swagger docs
	v1 "github.com/anas-dev-92/FoodHive/registration/src/v1"

	// Middlewares - Currently implemented
	mAP "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/ap"
	mAR "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/ar"
	mCatchWeight "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/catch_weight"
	mCustomer "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/customer"
	mInventory "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/inventory"
	mPicking "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/picking"
	mPricing "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/pricing"
	mProduct "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/product"
	mPurchaseOrder "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/purchase_order"
	mSalesOrder "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/sales_order"
	mVendor "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/vendor"
	mWarehouse "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/warehouse"
	// TODO: Uncomment as middlewares are implemented
	// mEmployee "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/employee"
	// mBank "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/bank"
	// mGL "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/gl"
	// mPayroll "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/payroll"
	// mPricing "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/pricing"
	// mWMS "github.com/anas-dev-92/FoodHive/registration/src/v1/middlewares/wms"
)

// @title ERP System API Documentation
// @version 1.0
// @description Enterprise Resource Planning System API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@foodhive.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables (optional - can also use system env vars)
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	// Parse configuration
	var config env.Config
	_, err = envconfig.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	// Initialize database connection
	db, err := postgres.New(config.DBConfig.ConnString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✓ Database connected")

	// Initialize JWT service
	jwtService := jwt.New(config.JWTSecret.JWTSecret)
	log.Println("✓ JWT service initialized")

	// Initialize storage service (MinIO) - optional for development
	var storageService *storage.MinioStorageService
	storageService, err = storage.New(
		config.StorageHost,
		config.StorageKey,
		config.StorageSecret,
		config.StorageSSL,
	)
	if err != nil {
		log.Printf("⚠ Storage service not available (MinIO): %v", err)
		log.Println("  File uploads will not work until MinIO is configured")
	} else {
		log.Println("✓ Storage service initialized")
	}

	// Initialize auth service
	authService := auth.New(db)
	log.Println("✓ Auth service initialized")

	// Create router
	app := chi.NewRouter()

	// Global middlewares
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(middleware.RequestID)
	app.Use(middleware.RealIP)
	app.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ===========================================
	// Service Injection Middlewares
	// ===========================================

	// Currently implemented
	app.Use(mCustomer.New(db))
	app.Use(mProduct.New(db))
	app.Use(mVendor.New(db))
	app.Use(mWarehouse.New(db))
	app.Use(mInventory.New(db))
	app.Use(mPurchaseOrder.New(db))
	app.Use(mSalesOrder.New(db))
	app.Use(mPicking.New(db))
	app.Use(mPricing.New(db))
	app.Use(mAR.New(db))
	app.Use(mAP.New(db))
	app.Use(mCatchWeight.New(db))

	// TODO: Uncomment as middlewares are implemented
	// Phase 1: Foundation
	// app.Use(mEmployee.New(db))

	// Phase 3: Advanced
	// app.Use(mGL.New(db))
	// app.Use(mBank.New(db))
	// app.Use(mPayroll.New(db))

	// Phase 4: WMS
	// app.Use(mPicking.New(db))
	// app.Use(mWMS.New(db))

	// ===========================================
	// Health Check
	// ===========================================
	app.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"erp-api"}`))
	})

	// ===========================================
	// Swagger Documentation
	// ===========================================
	// TODO: Enable after generating swagger docs
	// app.Get("/swagger/*", httpSwagger.WrapHandler)

	// ===========================================
	// API Routes
	// ===========================================
	app.Mount("/v1", v1.Router(app, jwtService, db, storageService, authService))

	// Start server
	port := ":8080"
	log.Printf("🚀 ERP Server starting on %s", port)
	log.Println("📚 Swagger docs available at http://localhost:8080/swagger/index.html")

	err = http.ListenAndServe(port, app)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
