# Останавливает Go-бинарники из bin/ и опционально docker infra.

param(
    [switch]$WithInfra
)

$ErrorActionPreference = "SilentlyContinue"
$Root = Split-Path -Parent $PSScriptRoot
$bin = Join-Path $Root "bin"

Get-Process | Where-Object { $_.Path -like "$bin\*" } | Stop-Process -Force
Write-Host "Stopped Go services from $bin"

if ($WithInfra) {
    Set-Location $Root
    docker compose down
    Write-Host "Docker infrastructure stopped"
}
