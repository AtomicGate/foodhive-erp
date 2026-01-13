# FoodHive ERP Frontend Startup Script
# Run this script to start the frontend development server

Write-Host "Starting FoodHive ERP Frontend..." -ForegroundColor Green
Write-Host "Backend API: http://localhost:8080" -ForegroundColor Cyan
Write-Host "Frontend: http://localhost:3000" -ForegroundColor Cyan

# Check if port 3000 is already in use
$portInUse = Get-NetTCPConnection -LocalPort 3000 -ErrorAction SilentlyContinue
if ($portInUse) {
    Write-Host ""
    Write-Host "WARNING: Port 3000 is already in use. Attempting to free it..." -ForegroundColor Yellow
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
        Write-Host "SUCCESS: Port 3000 is now free" -ForegroundColor Green
    } catch {
        Write-Host "  Error freeing port: $_" -ForegroundColor Red
        Write-Host ""
        Write-Host "Please manually stop the process using port 3000 and try again." -ForegroundColor Yellow
        exit 1
    }
}

Write-Host ""

# Check if node_modules exists or if vite is missing
$needsInstall = (-not (Test-Path "node_modules")) -or (-not (Test-Path "node_modules\.bin\vite.exe") -and -not (Test-Path "node_modules\.bin\vite"))
if ($needsInstall) {
    Write-Host "Installing dependencies..." -ForegroundColor Yellow
    # Check if pnpm is available first
    $pnpmAvailable = $false
    try {
        $pnpmVersion = pnpm --version 2>$null
        if ($pnpmVersion) {
            $pnpmAvailable = $true
            Write-Host "Using pnpm to install..." -ForegroundColor Cyan
            pnpm install
        }
    } catch {
        Write-Host "pnpm not found, using npm with --legacy-peer-deps..." -ForegroundColor Yellow
        npm install --legacy-peer-deps
    }
    Write-Host ""
}

# Check if pnpm is available, otherwise use npm
$usePnpm = $false
try {
    $pnpmVersion = pnpm --version 2>$null
    if ($pnpmVersion) {
        $usePnpm = $true
        Write-Host "Using pnpm..." -ForegroundColor Cyan
    }
} catch {
    Write-Host "pnpm not found, using npm..." -ForegroundColor Yellow
}

Write-Host "Starting development server..." -ForegroundColor Green
Write-Host ""

if ($usePnpm) {
    pnpm dev
} else {
    npm run dev
}
