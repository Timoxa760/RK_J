#!/usr/bin/env bash
# Demo flow для жюри — 6 шагов «Поток»
# Требует: api-gateway :8000, JWT (demo login 0000)
# Использование: ./demo_flow.sh
# Копия для back: scripts/demo_flow.sh

set -euo pipefail

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
curl -sf "${AUTH[@]}" "$API/credits/dashboard" | python3 -m json.tool 2>/dev/null | head -20 || echo "(credits stub)"

step "6. Insights"
curl -sf "${AUTH[@]}" "$API/insights" | head -c 200; echo

echo
echo "Demo flow complete. UI: http://localhost:3000 (front, demo code 0000)"
