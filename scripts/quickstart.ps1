#!/usr/bin/env pwsh
# Conduit Quickstart — runs the gateway with all features demonstrated
# Requires: Go 1.23+, API keys set in environment

param(
    [switch]$NoBuild
)

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot

Write-Host "╔═══════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║     Conduit AI Gateway — Quickstart          ║" -ForegroundColor Cyan
Write-Host "╚═══════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Check for API keys
if (-not $env:OPENAI_API_KEY -and -not $env:ANTHROPIC_API_KEY) {
    Write-Host "⚠  No API keys found. Set at least one:" -ForegroundColor Yellow
    Write-Host "   `$env:OPENAI_API_KEY = `"sk-...`"" -ForegroundColor Gray
    Write-Host "   `$env:ANTHROPIC_API_KEY = `"sk-ant-...`"" -ForegroundColor Gray
    Write-Host ""
    Write-Host "   Using test config (no real API calls will succeed)" -ForegroundColor Yellow
    $configPath = Join-Path $root "config.test.yaml"
} else {
    $configPath = Join-Path $root "config.example.yaml"
}

# Build if needed
if (-not $NoBuild) {
    Write-Host "▸ Building Conduit..." -ForegroundColor Green
    Set-Location $root
    go build -o (Join-Path $root "build\conduit.exe") ./cmd\gateway\
    if ($LASTEXITCODE -ne 0) { Write-Host "✗ Build failed" -ForegroundColor Red; exit 1 }
    Write-Host "  ✓ Build complete" -ForegroundColor Green
}

# Start the gateway
$gatewayPath = Join-Path $root "build\conduit.exe"
Write-Host "▸ Starting Conduit on localhost:8080..." -ForegroundColor Green
Write-Host "  Config: $configPath" -ForegroundColor Gray
Write-Host ""

$proc = Start-Process -NoNewWindow -FilePath $gatewayPath -ArgumentList "--config", $configPath -PassThru
Start-Sleep -Seconds 2

try {
    # Health check
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -ErrorAction Stop
    Write-Host "  ✓ Health: $($health.status)" -ForegroundColor Green

    # Model listing
    $models = Invoke-RestMethod -Uri "http://localhost:8080/v1/models" -ErrorAction Stop
    Write-Host "  ✓ Models: $($models.data.Count) available" -ForegroundColor Green

    # Stats
    $stats = Invoke-RestMethod -Uri "http://localhost:8080/v1/stats" -ErrorAction Stop
    Write-Host "  ✓ Stats: $($stats.total_requests) requests, $($stats.total_cost)" -ForegroundColor Green

    Write-Host ""
    Write-Host "╔═══════════════════════════════════════════════╗" -ForegroundColor Green
    Write-Host "║     Conduit is running! Try it:               ║" -ForegroundColor Green
    Write-Host "║                                               ║" -ForegroundColor Green
    Write-Host "║  curl http://localhost:8080/v1/chat/          ║" -ForegroundColor Green
    Write-Host "║    completions \                              ║" -ForegroundColor Green
    Write-Host "║    -H 'Content-Type: application/json' \      ║" -ForegroundColor Green
    Write-Host "║    -d '{ \"model\": \"gpt-4o-mini\",          ║" -ForegroundColor Green
    Write-Host "║           \"messages\": [{\"role\":           ║" -ForegroundColor Green
    Write-Host "║           \"user\",\"content\":               ║" -ForegroundColor Green
    Write-Host "║           \"Hello!\"}] }'                     ║" -ForegroundColor Green
    Write-Host "║                                               ║" -ForegroundColor Green
    Write-Host "║  Press Ctrl+C to stop                         ║" -ForegroundColor Green
    Write-Host "╚═══════════════════════════════════════════════╝" -ForegroundColor Green

    # Keep running until user presses Ctrl+C
    Wait-Process -Id $proc.Id
}
catch {
    Write-Host "✗ Error: $_" -ForegroundColor Red
}
finally {
    if ($proc -and !$proc.HasExited) {
        Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
        Write-Host "  Conduit stopped." -ForegroundColor Yellow
    }
}
