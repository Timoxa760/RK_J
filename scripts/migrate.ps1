# Применяет миграции PostgreSQL и ClickHouse через docker compose profile tools.

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

Write-Host "Running PostgreSQL migrations..."
docker compose --profile tools run --rm migrate-pg
if ($LASTEXITCODE -ne 0) {
    Write-Error "PostgreSQL migration failed"
}

Write-Host "Running ClickHouse migrations..."
docker compose --profile tools run --rm migrate-ch
if ($LASTEXITCODE -ne 0) {
    Write-Error "ClickHouse migration failed"
}

Write-Host "Migrations complete."
