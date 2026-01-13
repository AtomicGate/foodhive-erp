# 🐝 FoodHive ERP System

A comprehensive Enterprise Resource Planning (ERP) system designed for food distribution businesses. Built with Go backend and React frontend.

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?logo=postgresql)
![License](https://img.shields.io/badge/license-MIT-green)

## 📋 Table of Contents

- [Overview](#-overview)
- [Features Status](#-features-status)
- [Tech Stack](#-tech-stack)
- [Getting Started](#-getting-started)
- [Database Setup](#-database-setup)
- [Running the Application](#-running-the-application)
- [Project Structure](#-project-structure)
- [API Documentation](#-api-documentation)
- [Recent Updates](#-recent-updates)
- [Known Issues & Solutions](#-known-issues--solutions)
- [Roadmap](#-roadmap)
- [Contributing](#-contributing)

---

## 🎯 Overview

FoodHive ERP is a full-stack enterprise application designed to manage:
- **Master Data**: Employees, Customers, Vendors, Products, Warehouses
- **Transactions**: Sales Orders, Purchase Orders, Inventory Management
- **Financials**: Accounts Receivable (AR), Accounts Payable (AP), General Ledger
- **Operations**: Warehouse Management System (WMS), Catch Weight handling

---

## ✅ Features Status

### Phase 1: Foundation (✅ Complete)

| Feature | Backend | Frontend | Status |
|---------|---------|----------|--------|
| User Authentication (JWT) | ✅ | ✅ | Complete |
| Employee Management | ✅ | ✅ | Complete |
| Role & Permission System | ✅ | ✅ | Complete |
| Customer Management | ✅ | ✅ | Complete |
| Vendor Management | ✅ | ✅ | Complete |
| Product Management | ✅ | ✅ | Complete |
| Product Units Management | ✅ | ✅ | Complete |
| Warehouse Management | ✅ | ✅ | Complete |
| Department Management | ✅ | ✅ | Complete |

### Phase 2: Core ERP (🔄 In Progress)

| Feature | Backend | Frontend | Status |
|---------|---------|----------|--------|
| Sales Order Entry | ✅ | 🔄 | In Progress |
| Purchase Orders | ✅ | 🔄 | In Progress |
| Inventory Control | ✅ | 🔄 | In Progress |
| Accounts Receivable (AR) | ✅ | ⏳ | Backend Ready |
| Accounts Payable (AP) | ✅ | ⏳ | Backend Ready |

### Phase 3: Advanced Features (⏳ Planned)

| Feature | Backend | Frontend | Status |
|---------|---------|----------|--------|
| General Ledger | ⏳ | ⏳ | Planned |
| Pricing & Cost Management | ⏳ | ⏳ | Planned |
| Bank Reconciliation | ⏳ | ⏳ | Planned |
| Payroll Integration | ⏳ | ⏳ | Planned |
| Reporting & Analytics | ⏳ | ⏳ | Planned |

### Phase 4: WMS (⏳ Planned)

| Feature | Backend | Frontend | Status |
|---------|---------|----------|--------|
| Picking & Routing | ⏳ | ⏳ | Planned |
| Cycle Counts | ⏳ | ⏳ | Planned |
| Skid Tracking | ⏳ | ⏳ | Planned |
| Barcode Integration | ⏳ | ⏳ | Planned |

---

## 🛠️ Tech Stack

### Backend
| Component | Technology |
|-----------|------------|
| Language | Go 1.25 |
| Web Framework | Chi v5 |
| Database | PostgreSQL 14+ |
| Database Driver | pgx v5 (with connection pooling) |
| Object Storage | MinIO |
| Authentication | JWT (HS256) |
| API Docs | Swagger/OpenAPI |

### Frontend
| Component | Technology |
|-----------|------------|
| Framework | React 18 |
| Build Tool | Vite |
| Language | TypeScript |
| UI Components | shadcn/ui |
| Styling | Tailwind CSS |
| State Management | React Query (TanStack Query) |
| Routing | React Router v6 |
| HTTP Client | Axios |
| Forms | React Hook Form |
| Tables | TanStack Table |

---

## 🚀 Getting Started

### Prerequisites

- **Go** 1.23+ ([Download](https://golang.org/dl/))
- **Node.js** 18+ ([Download](https://nodejs.org/))
- **PostgreSQL** 14+ ([Download](https://www.postgresql.org/download/))
- **MinIO** (optional, for file storage)

### Clone the Repository

```bash
git clone https://github.com/yourusername/FoodHive.git
cd FoodHive-main
```

---

## 🗄️ Database Setup

### 1. Create Database

```sql
-- Connect to PostgreSQL as superuser
CREATE USER erp WITH PASSWORD 'erp_password' SUPERUSER;
CREATE DATABASE erp_db OWNER erp;
```

### 2. Run Migrations (In Order!)

⚠️ **Important**: Run these SQL scripts in the exact order shown:

```bash
# 1. Core tables (employees, roles, pages, etc.)
psql -U erp -d erp_db -f sql/001_core_tables.sql

# 2. Transaction tables
psql -U erp -d erp_db -f sql/002_transactions_tables.sql

# 3. Financial tables
psql -U erp -d erp_db -f sql/003_financial_tables.sql

# 4. Insert default roles (MUST run before admin user)
psql -U erp -d erp_db -f sql/004_insert_roles.sql

# 5. Create default admin user
psql -U erp -d erp_db -f sql/005_insert_default_admin.sql

# 6. Schema updates & permissions (run AFTER admin exists)
psql -U erp -d erp_db -f sql/003_schema_updates.sql
```

### 3. Default Admin Credentials

After running migrations:
- **Email**: `admin@foodhive.com`
- **Password**: `admin123`

⚠️ **Change this password immediately in production!**

---

## ▶️ Running the Application

### Option 1: Using PowerShell Scripts (Windows)

```powershell
# Terminal 1 - Start Backend
.\start-server.ps1

# Terminal 2 - Start Frontend
cd foodhive-erp-frontend
.\start-frontend.ps1
```

### Option 2: Manual Start

**Backend:**
```bash
cd registration
go mod tidy
go run main.go
```

**Frontend:**
```bash
cd foodhive-erp-frontend
npm install --legacy-peer-deps
npm run dev
```

### Access Points

| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| Backend API | http://localhost:8080 |
| Swagger Docs | http://localhost:8080/swagger/index.html |
| Health Check | http://localhost:8080/health |

---

## 📁 Project Structure

```
FoodHive-main/
├── 📂 core/                        # Shared infrastructure
│   ├── auth/                       # Authentication service
│   ├── jwt/                        # JWT token service
│   ├── postgres/                   # Database connection (pgxpool)
│   ├── storage/                    # MinIO file storage
│   └── utils/                      # Shared utilities
│
├── 📂 registration/                # Backend application
│   ├── src/v1/
│   │   ├── middlewares/            # Auth & service injection
│   │   ├── routes/                 # HTTP handlers
│   │   │   ├── login/              # Authentication
│   │   │   ├── employee/           # Employee CRUD
│   │   │   ├── customer/           # Customer CRUD
│   │   │   ├── vendor/             # Vendor CRUD
│   │   │   ├── product/            # Product & Units
│   │   │   ├── warehouse/          # Warehouse management
│   │   │   └── ...                 # Other modules
│   │   ├── services/               # Business logic layer
│   │   ├── models/                 # Data models & DTOs
│   │   └── utils/                  # Route utilities
│   ├── docs/                       # Swagger documentation
│   └── main.go                     # Entry point
│
├── 📂 foodhive-erp-frontend/       # Frontend application
│   └── client/
│       ├── src/
│       │   ├── components/         # Reusable UI components
│       │   │   └── ui/             # shadcn/ui components
│       │   ├── pages/              # Page components
│       │   │   ├── employees/      # Employee management
│       │   │   ├── customers/      # Customer management
│       │   │   ├── vendors/        # Vendor management
│       │   │   ├── products/       # Product management
│       │   │   └── ...             # Other pages
│       │   ├── services/           # API service layer
│       │   ├── hooks/              # Custom React hooks
│       │   ├── lib/                # Utilities
│       │   └── App.tsx             # Main app component
│       └── index.html
│
├── 📂 sql/                         # Database migrations
│   ├── 001_core_tables.sql
│   ├── 002_transactions_tables.sql
│   ├── 003_financial_tables.sql
│   ├── 003_schema_updates.sql
│   ├── 004_insert_roles.sql
│   └── 005_insert_default_admin.sql
│
├── start-server.ps1                # Backend startup script
└── README.md
```

---

## 📚 API Documentation

### Authentication

All endpoints (except `/v1/login` and `/health`) require JWT authentication.

```bash
# Login
curl -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@foodhive.com", "password": "admin123"}'

# Response: { "token": "eyJhbG..." }

# Use token in requests
curl http://localhost:8080/v1/employees/list \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### API Modules

| Module | Base Path | Description |
|--------|-----------|-------------|
| Auth | `POST /v1/login` | User authentication |
| Employees | `/v1/employees/*` | Employee CRUD |
| Roles | `/v1/roles/*` | Role management |
| Customers | `/v1/customers/*` | Customer CRUD |
| Vendors | `/v1/vendors/*` | Vendor CRUD |
| Products | `/v1/products/*` | Product catalog |
| Product Units | `/v1/product-units/*` | Unit conversions |
| Warehouses | `/v1/warehouses/*` | Warehouse management |
| Inventory | `/v1/inventory/*` | Stock control |
| Sales Orders | `/v1/sales-orders/*` | Order entry |
| Purchase Orders | `/v1/purchase-orders/*` | Purchasing |
| AR | `/v1/ar/*` | Accounts Receivable |
| AP | `/v1/ap/*` | Accounts Payable |

---

## 🔄 Recent Updates

### January 2026

#### Backend Fixes
- ✅ Fixed `emp_status` enum mismatch (`ACTIVE` → `CONTINUED`)
- ✅ Implemented database connection pooling (`pgxpool`) to fix "conn busy" errors
- ✅ Fixed role service column name (`description` → `role_desc`)
- ✅ Added `base_unit` to product update endpoint
- ✅ Improved error logging in helper functions
- ✅ Fixed JSON parsing to allow unknown fields from frontend

#### Frontend Fixes
- ✅ Fixed employee list disappearing on edit (browser autofill issue)
- ✅ Added proper autocomplete prevention on all forms
- ✅ Fixed double-wrapped API response handling in services
- ✅ Implemented employee status management (Activate/Suspend/On Leave/Resign)
- ✅ Fixed product edit and unit management functionality
- ✅ Added product reactivation feature
- ✅ Improved DataTable with memoization to prevent unnecessary re-renders
- ✅ Added React Query optimizations (`staleTime`, `placeholderData`)

#### Database
- ✅ Created `005_insert_default_admin.sql` for proper migration order
- ✅ Fixed foreign key constraint issues in `emp_page` table
- ✅ Updated `003_schema_updates.sql` with conditional inserts

---

## ⚠️ Known Issues & Solutions

### 1. PowerShell Execution Policy Error
```powershell
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process -Force
```

### 2. npm Dependency Conflicts
```bash
npm install --legacy-peer-deps
```

### 3. "conn busy" Database Errors
- **Solution**: Already fixed with `pgxpool` connection pooling
- Restart the backend server after pulling latest changes

### 4. 403 Forbidden After Login
- **Cause**: Admin user exists but has no page permissions
- **Solution**: Run `003_schema_updates.sql` after creating admin user

### 5. Employee List Disappears on Edit Click
- **Cause**: Browser autofill filling search box with saved credentials
- **Solution**: Already fixed with autocomplete prevention attributes

---

## 🗺️ Roadmap

### Q1 2026
- [ ] Complete Sales Order workflow
- [ ] Complete Purchase Order workflow
- [ ] Inventory transactions UI
- [ ] Dashboard analytics widgets

### Q2 2026
- [ ] Accounts Receivable UI
- [ ] Accounts Payable UI
- [ ] Invoice generation & printing
- [ ] Payment processing

### Q3 2026
- [ ] General Ledger implementation
- [ ] Financial reporting
- [ ] Multi-warehouse support
- [ ] Batch/Lot tracking

### Q4 2026
- [ ] WMS: Picking & routing
- [ ] Mobile app for warehouse
- [ ] Barcode scanning integration
- [ ] Advanced analytics & BI

---

## 🧪 Testing

```bash
# Backend tests
cd registration
go test ./test/... -v

# Frontend tests (if configured)
cd foodhive-erp-frontend
npm test
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👥 Team

Built with ❤️ by the FoodHive Team

**Last Updated**: January 14, 2026
