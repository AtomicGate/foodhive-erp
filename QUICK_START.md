# FoodHive ERP - Quick Start Guide

## 🚀 Quick Start Options

### Option 1: Start Everything at Once (Recommended)
```powershell
.\start-all.ps1
```
This will start both backend and frontend in separate windows.

### Option 2: Start Backend Only
```powershell
cd registration
.\start-server.ps1
```

### Option 3: Start Frontend Only
```powershell
cd foodhive-erp-frontend
.\start-frontend.ps1
```

## 📋 Prerequisites

1. **PostgreSQL** must be running with:
   - Database: `FoodHive`
   - User: `postgres`
   - Password: `bakri314`
   - Port: `5432`

2. **Go** installed (for backend)

3. **Node.js** and **npm** installed (for frontend)

## 🔧 Configuration

### Backend Configuration
Edit `registration/start-server.ps1` or `registration/env.example` to change:
- Database connection string
- JWT secret
- Storage (MinIO) settings

### Frontend Configuration
The frontend automatically connects to `http://localhost:8080` (backend API).

## 🌐 Access Points

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Health Check**: http://localhost:8080/health
- **Swagger Docs**: http://localhost:8080/swagger/index.html

## 🔐 Default Login

- **Email**: `admin@foodhive.com`
- **Password**: `password123`

## 📝 Notes

- Backend uses `.env` file if available, otherwise uses system environment variables
- Frontend will auto-install dependencies if `node_modules` doesn't exist
- Both servers support hot-reload during development

## 🛠️ Troubleshooting

### Backend won't start
- Check PostgreSQL is running
- Verify database credentials in `start-server.ps1`
- Check port 8080 is not in use

### Frontend won't start
- Run `npm install` manually in `foodhive-erp-frontend` folder
- Check port 3000 is not in use
- Verify backend is running first

### Can't login
- Ensure backend is running
- Check database has admin user: `admin@foodhive.com` / `password123`
- Verify permissions are set in `emp_page` table
