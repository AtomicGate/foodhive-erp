# 🏭 ERP System

Enterprise Resource Planning System built with Go.

## 📋 Features

### Phase 1: Foundation
- ✅ Employee Management
- ✅ Warehouse Management
- ✅ Product Management
- ✅ Customer Management
- ✅ Vendor Management
- ✅ Inventory Control

### Phase 2: Core ERP
- 🔄 Sales Order Entry
- 🔄 Purchase Orders
- 🔄 Accounts Receivable (AR)
- 🔄 Accounts Payable (AP)

### Phase 3: Advanced
- ⏳ General Ledger
- ⏳ Pricing & Cost Management
- ⏳ Bank & Reconciliation
- ⏳ Payroll Integration

### Phase 4: WMS
- ⏳ Picking & Routing
- ⏳ Warehouse Management System
- ⏳ Cycle Counts
- ⏳ Skid Tracking

## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.23 |
| Web Framework | Chi v5 |
| Database | PostgreSQL |
| Object Storage | MinIO |
| Authentication | JWT |
| API Docs | Swagger |

## 🚀 Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL 14+
- MinIO (for file storage)

### Database Setup

```bash
# Create database user and database
psql -U postgres -c "CREATE USER erp PASSWORD 'erp_password' SUPERUSER"
psql -U postgres -c "CREATE DATABASE erp_db OWNER erp"

# Run migrations
cat ./sql/001_core_tables.sql | psql -U erp -d erp_db
cat ./sql/002_transactions_tables.sql | psql -U erp -d erp_db
cat ./sql/003_financial_tables.sql | psql -U erp -d erp_db
```

### Environment Setup

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your configuration
```

### Run the Application

```bash
cd registration

# Install dependencies
go mod tidy

# Run the server
go run main.go
```

### Access the API

- API: http://localhost:8080
- Swagger Docs: http://localhost:8080/swagger/index.html
- Health Check: http://localhost:8080/health

## 📁 Project Structure

```
erp/
├── core/                    # Shared infrastructure
│   ├── auth/               # Authentication service
│   ├── jwt/                # JWT service
│   ├── postgres/           # Database layer
│   ├── storage/            # MinIO storage
│   └── utils/              # Utilities
│
├── registration/           # Main application
│   ├── src/v1/
│   │   ├── middlewares/    # Service injection
│   │   ├── routes/         # HTTP handlers
│   │   ├── services/       # Business logic
│   │   └── models/         # Data models
│   ├── docs/               # Swagger docs
│   └── test/               # Tests
│
├── sql/                    # Database migrations
├── scripts/                # Automation scripts
└── docs/                   # Documentation
```

## 📚 API Modules

| Module | Base Path | Description |
|--------|-----------|-------------|
| Auth | `/v1/login` | Authentication |
| Employees | `/v1/employees` | Employee management |
| Warehouses | `/v1/warehouses` | Warehouse management |
| Products | `/v1/products` | Product catalog |
| Customers | `/v1/customers` | Customer management |
| Vendors | `/v1/vendors` | Vendor management |
| Inventory | `/v1/inventory` | Stock control |
| Sales Orders | `/v1/sales-orders` | Order entry |
| Purchase Orders | `/v1/purchase-orders` | Purchasing |
| AR | `/v1/ar` | Accounts Receivable |
| AP | `/v1/ap` | Accounts Payable |
| GL | `/v1/gl` | General Ledger |
| WMS | `/v1/wms` | Warehouse operations |

## 🔐 Authentication

All endpoints (except `/v1/login` and `/health`) require JWT authentication.

```bash
# Login to get token
curl -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@company.com", "password": "password"}'

# Use token in requests
curl http://localhost:8080/v1/customers \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 🧪 Testing

```bash
cd registration
go test ./test/... -v
```

## 📝 License

MIT License

## 👥 Team

Built with ❤️ for enterprise resource planning.

