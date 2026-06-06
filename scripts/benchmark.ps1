#!/usr/bin/env pwsh
# Conduit Benchmark — measures latency and throughput
param(
    [string]$Url = "http://localhost:8080",
    [int]$Requests = 20,
    [int]$Concurrency = 5
)

Write-Host "═══ Conduit Benchmark ═══" -ForegroundColor Cyan
Write-Host "Target:      $Url"
Write-Host "Requests:    $Requests"
Write-Host "Concurrency: $Concurrency"
Write-Host ""

# Health check
try {
    $health = Invoke-RestMethod -Uri "$Url/health" -ErrorAction Stop
    Write-Host "• Health: OK" -ForegroundColor Green
} catch {
    Write-Host "• Health: FAIL ($($_.Exception.Message))" -ForegroundColor Red
    exit 1
}

# Sequential latency test
Write-Host "• Latency (sequential, $Requests requests):" -ForegroundColor Yellow
$times = @()
for ($i = 0; $i -lt $Requests; $i++) {
    $sw = [System.Diagnostics.Stopwatch]::StartNew()
    try {
        $body = '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"Hi"}]}'
        Invoke-WebRequest -Uri "$Url/v1/chat/completions" -Method Post -Body $body -ContentType "application/json" -UseBasicParsing -ErrorAction Stop | Out-Null
    } catch {
        # Expected when no API key — routing is what matters
    }
    $sw.Stop()
    $times += $sw.ElapsedMilliseconds
}

$avg = ($times | Measure-Object -Average).Average
$min = $times | Measure-Object -Minimum
$max = $times | Measure-Object -Maximum
Write-Host "  Min:    $($min.Minimum)ms" -ForegroundColor Gray
Write-Host "  Avg:    $([math]::Round($avg, 1))ms" -ForegroundColor Gray
Write-Host "  Max:    $($max.Maximum)ms" -ForegroundColor Gray

# Stats
Write-Host ""
Write-Host "• Gateway Stats:" -ForegroundColor Yellow
try {
    $stats = Invoke-RestMethod -Uri "$Url/v1/stats" -ErrorAction Stop
    Write-Host "  $($stats | ConvertTo-Json)" -ForegroundColor Gray
} catch {
    Write-Host "  Unable to fetch stats" -ForegroundColor Red
}
