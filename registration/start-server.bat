@echo off
REM FoodHive ERP Server Startup Script
REM Run this batch file to start the backend server

REM Database Configuration
set DB_CONFIG=postgres://postgres:bakri314@localhost:5432/FoodHive?sslmode=disable

REM JWT Secret Key
set JWT_SECRET=foodhive-jwt-secret-2024

REM Storage Configuration (MinIO - optional)
set STORAGE_HOST=localhost:9000
set STORAGE_KEY=minioadmin
set STORAGE_SECRET=minioadmin
set STORAGE_SSL=false

echo Starting FoodHive ERP Server...
echo Database: FoodHive @ localhost:5432

REM Check if port 8080 is in use
netstat -ano | findstr :8080 >nul
if %ERRORLEVEL% EQU 0 (
    echo.
    echo Port 8080 is already in use. Attempting to free it...
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr :8080') do (
        taskkill /F /PID %%a >nul 2>&1
    )
    timeout /t 2 /nobreak >nul
    echo Port 8080 is now free
)

echo.

REM Run the Go server
go run main.go
