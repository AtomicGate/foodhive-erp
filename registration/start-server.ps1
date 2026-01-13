# FoodHive ERP Server Startup Script
# Run this script to start the backend server with all configuration

# Database Configuration
$env:DB_CONFIG = "postgres://postgres:bakri314@localhost:5432/FoodHive?sslmode=disable"

# JWT Secret Key
$env:JWT_SECRET = "foodhive-jwt-secret-2024"

# Storage Configuration (MinIO - optional)
$env:STORAGE_HOST = "localhost:9000"
$env:STORAGE_KEY = "minioadmin"
$env:STORAGE_SECRET = "minioadmin"
$env:STORAGE_SSL = "false"

Write-Host "Starting FoodHive ERP Server..." -ForegroundColor Green
Write-Host "Database: FoodHive @ localhost:5432" -ForegroundColor Cyan

# Check if port 8080 is already in use
$portInUse = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue
if ($portInUse) {
    Write-Host ""
    Write-Host "WARNING: Port 8080 is already in use. Attempting to free it..." -ForegroundColor Yellow
    try {
        $portInUse | ForEach-Object {
            $processId = $_.OwningProcess
            $process = Get-Process -Id $processId -ErrorAction SilentlyContinue
            if ($process) {
                Write-Host "  Killing process: $($process.ProcessName) (PID: $processId)" -ForegroundColor Yellow
                Stop-Process -Id $processId -Force -ErrorAction SilentlyContinue
            }
        }
        Start-Sleep -Seconds 2
        Write-Host "SUCCESS: Port 8080 is now free" -ForegroundColor Green
    } catch {
        Write-Host "  Error freeing port: $_" -ForegroundColor Red
        Write-Host ""
        Write-Host "Please manually stop the process using port 8080 and try again." -ForegroundColor Yellow
        Write-Host "You can find it with: Get-NetTCPConnection -LocalPort 8080" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host ""

# Run the Go server
go run main.go
