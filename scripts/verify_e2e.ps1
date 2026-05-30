# Smoke E2E: register -> login -> manual expense -> dashboard categories (PostgreSQL).

$ErrorActionPreference = "Stop"
$Base = "http://127.0.0.1:8000/api/v1"
$Phone = "+7900" + (Get-Random -Minimum 1000000 -Maximum 9999999)

Write-Host "1. Register $Phone"
Invoke-RestMethod -Method Post -Uri "$Base/auth/register" `
    -ContentType "application/json" `
    -Body (@{ phone = $Phone; password = "secret12345" } | ConvertTo-Json) | Out-Null

Write-Host "2. Login"
$login = Invoke-RestMethod -Method Post -Uri "$Base/auth/login" `
    -ContentType "application/json" `
    -Body (@{ phone = $Phone; password = "secret12345" } | ConvertTo-Json)

$token = $login.access_token
if (-not $token) { $token = $login.token }
if (-not $token) { throw "No JWT in login response" }

$headers = @{ Authorization = "Bearer $token" }

Write-Host "3. POST /receipt/manual"
$manualJson = (@{
    store    = "TestStore"
    amount   = 1500
    category = "Products"
    date     = (Get-Date -Format "yyyy-MM-dd")
} | ConvertTo-Json -Compress)
$expense = Invoke-RestMethod -Method Post -Uri "$Base/receipt/manual" `
    -Headers $headers `
    -ContentType "application/json; charset=utf-8" `
    -Body ([System.Text.Encoding]::UTF8.GetBytes($manualJson))

Write-Host "   saved: $($expense.receipt_id) amount=$($expense.amount)"

Start-Sleep -Seconds 1

Write-Host "4. GET /dashboard/categories"
$cats = Invoke-RestMethod -Method Get -Uri "$Base/dashboard/categories" -Headers $headers

$total = 0
foreach ($c in $cats.categories) {
    $total += $c.total
    Write-Host "   category: $($c.name) total=$($c.total)"
}

if ($total -lt 1500) {
    throw "Expected dashboard total >= 1500, got $total"
}

Write-Host ""
Write-Host "E2E PASSED: expense persisted and visible on dashboard"
