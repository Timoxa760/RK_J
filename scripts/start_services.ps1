# Запуск всех 15 Go-сервисов против локальной инфраструктуры (PostgreSQL, Kafka, ClickHouse).
# Предварительно: scripts\start_infra.ps1 и scripts\migrate.ps1

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$bin = Join-Path $Root "bin"
$logs = Join-Path $Root "logs"
$dataDir = Join-Path $Root "data"

if (-not (Test-Path $logs)) { New-Item -ItemType Directory -Path $logs | Out-Null }
if (-not (Test-Path $dataDir)) { New-Item -ItemType Directory -Path $dataDir | Out-Null }

$env:DEMO_MODE = "false"
$env:JWT_SECRET = "test-secret"
$env:JWT_ACCESS_TTL = "4h"
$env:DATABASE_URL = "postgres://postgres:postgres@127.0.0.1:5432/moneymind?sslmode=disable"
$env:EXPENSE_STORE_PATH = Join-Path $dataDir "expenses.json"

$env:KAFKA_BROKERS = "127.0.0.1:9092"
$env:KAFKA_URL = "127.0.0.1:9092"
$env:REDIS_URL = "127.0.0.1:6379"

$env:CLICKHOUSE_HOST = "127.0.0.1"
$env:CLICKHOUSE_URL = "http://127.0.0.1:8123"
$env:CLICKHOUSE_USER = "clickhouse_user"
$env:CLICKHOUSE_PASSWORD = "clickhouse_password"
$env:CLICKHOUSE_DB = "default"

# OnlySQ и Whisper — опционально; без ключей manual работает через regex
$env:ONLYSQ_BASE_URL = "https://api.onlysq.ru/v1"
if (-not $env:ONLYSQ_API_KEY) { $env:ONLYSQ_API_KEY = "" }
$env:WHISPER_URL = "http://127.0.0.1:9001"

# URL upstream-сервисов для api-gateway (локальный запуск без Docker DNS)
$env:USER_SERVICE_URL = "http://127.0.0.1:8001"
$env:RECEIPT_SERVICE_URL = "http://127.0.0.1:8002"
$env:SCRAPER_SERVICE_URL = "http://127.0.0.1:8003"
$env:CATEGORY_SERVICE_URL = "http://127.0.0.1:8004"
$env:BUDGET_SERVICE_URL = "http://127.0.0.1:8005"
$env:GOAL_SERVICE_URL = "http://127.0.0.1:8006"
$env:CREDIT_SERVICE_URL = "http://127.0.0.1:8009"
$env:REPORTING_SERVICE_URL = "http://127.0.0.1:8010"
$env:BANK_SERVICE_URL = "http://127.0.0.1:8011"
$env:AI_PROCESSOR_URL = "http://127.0.0.1:8100"
$env:ANALYTICS_SERVICE_URL = "http://127.0.0.1:8101"
$env:SOCIAL_SERVICE_URL = "http://127.0.0.1:8102"

function Start-Svc($name, $file) {
    if (-not (Test-Path $file)) {
        Write-Error "Missing binary: $file (run: make build)"
        exit 1
    }
    $log = Join-Path $logs "$name.log"
    $err = Join-Path $logs "$name.err"
    $process = Start-Process -FilePath $file -WindowStyle Hidden -PassThru `
        -RedirectStandardOutput $log -RedirectStandardError $err
    Write-Host "$name started (PID: $($process.Id))"
}

$services = @(
    @("user-service", "user-service.exe"),
    @("receipt-service", "receipt-service.exe"),
    @("scraper-service", "scraper-service.exe"),
    @("category-service", "category-service.exe"),
    @("budget-service", "budget-service.exe"),
    @("goal-service", "goal-service.exe"),
    @("credit-service", "credit-service.exe"),
    @("bank-service", "bank-service.exe"),
    @("ai-processor", "ai-processor.exe"),
    @("analytics-service", "analytics-service.exe"),
    @("gamification", "gamification.exe"),
    @("social-service", "social-service.exe"),
    @("notification-service", "notification-service.exe"),
    @("reporting-service", "reporting-service.exe"),
    @("api-gateway", "api-gateway.exe")
)

Write-Host "Starting 15 backend services (DEMO_MODE=false)..."
foreach ($svc in $services) {
    Start-Svc $svc[0] (Join-Path $bin $svc[1])
    Start-Sleep -Milliseconds 800
}

Write-Host ""
Write-Host "Gateway: http://127.0.0.1:8000/api/v1"
Write-Host "Health:  http://127.0.0.1:8000/health"
Write-Host "Logs:    $logs"
Write-Host "Stop:    powershell -File scripts\stop_stack.ps1"
