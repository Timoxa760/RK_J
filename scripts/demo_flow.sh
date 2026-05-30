#!/usr/bin/env bash
# Demo flow — шаги по docs/api/API_Contract.md
# Требует: api-gateway :8000, JWT (demo code 0000)

set -euo pipefail

export DEMO_MODE="${DEMO_MODE:-true}"
export JWT_SECRET="${JWT_SECRET:-test-secret}"
export USER_SERVICE_URL="${USER_SERVICE_URL:-http://127.0.0.1:8001}"
export RECEIPT_SERVICE_URL="${RECEIPT_SERVICE_URL:-http://127.0.0.1:8002}"
export AI_PROCESSOR_URL="${AI_PROCESSOR_URL:-http://127.0.0.1:8100}"

API="${API_BASE:-http://localhost:8000}/api/v1"
PHONE="${DEMO_PHONE:-+79991234567}"
CODE="${DEMO_CODE:-0000}"

step() { echo; echo "==> $1"; }

step "1. Login (demo code $CODE)"
TOKEN=$(curl -sf -X POST "$API/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"phone\":\"$PHONE\",\"code\":\"$CODE\"}" | \
  python3 -c "import sys,json; print(json.load(sys.stdin).get('access_token',''))" 2>/dev/null || true)

if [[ -z "$TOKEN" ]]; then
  echo "WARN: login failed — продолжаем без JWT (DEMO_MODE)"
  AUTH=()
else
  echo "JWT ok"
  AUTH=(-H "Authorization: Bearer $TOKEN")
fi

step "2. Голосовой ввод расхода"
curl -sf -X POST "$API/expenses/manual" \
  -H 'Content-Type: application/json' \
  "${AUTH[@]}" \
  -d "{\"user_id\":\"$PHONE\",\"raw_text\":\"купил продукты на 5000 и кроссовки за 16000\",\"source\":\"voice\"}" \
  | python3 -m json.tool 2>/dev/null || echo "(skip if ai-processor down)"

step "3. Dashboard — sankey"
curl -sf "${AUTH[@]}" "$API/dashboard/sankey" | python3 -c "import sys,json; d=json.load(sys.stdin); print('nodes:', len(d.get('nodes',[])))" 2>/dev/null || echo "FAIL sankey"

step "4. Dashboard — timemachine"
curl -sf "${AUTH[@]}" "$API/dashboard/timemachine" | head -c 200; echo

step "5. Кредитный светофор"
curl -sf "${AUTH[@]}" "$API/credits/dashboard" | python3 -m json.tool 2>/dev/null | head -20 || echo "(credits stub — 404 ожидаем)"

step "6. Insights"
curl -sf "${AUTH[@]}" "$API/insights" | head -c 200; echo || echo "(analytics stub — 404 ожидаем)"

echo
echo "Demo flow complete."
