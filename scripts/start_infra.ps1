# Поднимает инфраструктуру для полного локального стека (без Go-сервисов).
# Требуется Docker Desktop.

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

docker info *> $null
if ($LASTEXITCODE -ne 0) {
    Write-Error "Docker daemon is not running. Start Docker Desktop, then retry."
}

Write-Host "Starting infrastructure (postgres, redis, zookeeper, kafka, clickhouse)..."
docker compose up -d postgres redis zookeeper kafka clickhouse

Write-Host "Waiting for PostgreSQL..."
$ready = $false
for ($i = 0; $i -lt 60; $i++) {
    $status = docker inspect -f "{{.State.Health.Status}}" backend_postgres 2>$null
    if ($status -eq "healthy") {
        $ready = $true
        break
    }
    Start-Sleep -Seconds 2
}
if (-not $ready) {
    Write-Error "PostgreSQL did not become healthy in time. Is Docker Desktop running?"
}

Write-Host "Infrastructure is up."
Write-Host "Next: powershell -File scripts\migrate.ps1"
Write-Host "Then: powershell -File scripts\start_services.ps1"
