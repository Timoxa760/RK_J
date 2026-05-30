#!/usr/bin/env bash
# Health-check сервисов «Поток» (api-gateway + infra)
# Использование: ./health_check.sh
# Копия для back: scripts/health_check.sh

set -euo pipefail

API="${API_BASE:-http://localhost:8000}"
FAIL=0

check() {
  local name="$1" url="$2"
  if curl -sf --max-time 3 "$url" >/dev/null; then
    echo "OK  $name"
  else
    echo "FAIL $name ($url)"
    FAIL=1
  fi
}

echo "=== API Gateway ==="
check "api-gateway /health" "$API/health"

echo "=== Dashboard (auth may 401 — ожидаемо без token) ==="
code=$(curl -so /dev/null -w '%{http_code}' --max-time 3 "$API/api/v1/dashboard/sankey" || true)
if [[ "$code" == "200" || "$code" == "401" ]]; then
  echo "OK  dashboard/sankey (HTTP $code)"
else
  echo "FAIL dashboard/sankey (HTTP $code)"
  FAIL=1
fi

echo "=== Infra ==="
check "ClickHouse" "http://localhost:8123/?query=SELECT%201"

if command -v pg_isready >/dev/null 2>&1; then
  pg_isready -h localhost -p 5432 -U postgres >/dev/null 2>&1 && echo "OK  PostgreSQL" || { echo "FAIL PostgreSQL"; FAIL=1; }
else
  echo "SKIP PostgreSQL (pg_isready not installed)"
fi

if command -v redis-cli >/dev/null 2>&1; then
  redis-cli -h localhost ping >/dev/null 2>&1 && echo "OK  Redis" || { echo "FAIL Redis"; FAIL=1; }
else
  echo "SKIP Redis (redis-cli not installed)"
fi

exit "$FAIL"
