# Полный локальный стек: infra + migrations + build + 15 Go-сервисов.
# Whisper и OnlySQ не обязательны: manual expenses работают через regex + PostgreSQL.

param(
    [switch]$SkipBuild,
    [switch]$SkipInfra
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

$dataDir = Join-Path $Root "data"
if (-not (Test-Path $dataDir)) {
    New-Item -ItemType Directory -Path $dataDir | Out-Null
}

if (-not $SkipInfra) {
    & "$PSScriptRoot\start_infra.ps1"
    Start-Sleep -Seconds 3
    & "$PSScriptRoot\migrate.ps1"
}

if (-not $SkipBuild) {
    Write-Host "Building Go services..."
    make build
    if ($LASTEXITCODE -ne 0) {
        Write-Error "make build failed"
    }
}

& "$PSScriptRoot\start_services.ps1"

Write-Host ""
Write-Host "=== Stack ready ==="
Write-Host "Gateway:  http://127.0.0.1:8000/api/v1"
Write-Host "Verify:   powershell -File scripts\verify_e2e.ps1"
Write-Host "Front:    NUXT_PUBLIC_API_BASE=http://127.0.0.1:8000 NUXT_PUBLIC_DEMO_MODE=false npm run dev"
Write-Host "Stop:     powershell -File scripts\stop_stack.ps1"
