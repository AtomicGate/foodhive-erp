@echo off
REM FoodHive ERP Frontend Startup Script
REM Run this batch file to start the frontend development server

echo Starting FoodHive ERP Frontend...
echo Backend API: http://localhost:8080
echo Frontend: http://localhost:3000
echo.

REM Check if node_modules exists
if not exist "node_modules" (
    echo Installing dependencies...
    call npm install
    echo.
)

REM Check if pnpm is available
where pnpm >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    echo Using pnpm...
    echo Starting development server...
    echo.
    pnpm dev
) else (
    echo pnpm not found, using npm...
    echo Starting development server...
    echo.
    npm run dev
)
