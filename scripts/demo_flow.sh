#!/usr/bin/env bash
# Demo flow для жюри — критический путь «Поток»
# Требует: api-gateway :8000, JWT (demo login 0000)
# Использование: ./demo_flow.sh

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

step "2. Profile (онбординг-данные для demo)"
curl -sf -X PATCH "$API/users/me/profile" \
  -H 'Content-Type: application/json' \
  "${AUTH[@]}" \
  -d '{"active_income":150000,"passive_income":30000,"emergency_fund":340000,"goal_title":"Отпуск","goal_amount":150000,"goal_kind":"save","onboarding_completed":true}' \
  | python3 -m json.tool 2>/dev/null | head -15 || echo "(skip if user-service down)"

step "3. Голосовой ввод расхода"
curl -sf -X POST "$API/expenses/manual" \
  -H 'Content-Type: application/json' \
  "${AUTH[@]}" \
  -d "{\"user_id\":\"$PHONE\",\"raw_text\":\"купил продукты на 5000 и кроссовки за 16000\",\"source\":\"voice\"}" \
  | python3 -m json.tool 2>/dev/null || echo "(skip if ai-processor down)"

step "4. ФНС ticket (demo receipt)"
curl -sf -X POST "$API/fns/ticket" \
  -H 'Content-Type: application/json' \
  "${AUTH[@]}" \
  -d '{"fn":"9289000100123456","fd":"12345","fp":"1234567890","sum":"999.99","date":"2026-05-01","time":"12:00"}' \
  | python3 -c "import sys,json; d=json.load(sys.stdin); print('store:', d.get('store_name','?'))" 2>/dev/null || echo "(skip FNS)"

step "5. Dashboard — sankey + timemachine"
curl -sf "${AUTH[@]}" "$API/dashboard/sankey" | python3 -c "import sys,json; d=json.load(sys.stdin); print('nodes:', len(d.get('nodes',[])))" 2>/dev/null || echo "FAIL sankey"
curl -sf "${AUTH[@]}" "$API/dashboard/timemachine" | head -c 120; echo

step "6. ИИ-план и чат"
curl -sf "${AUTH[@]}" "$API/ai/plan" | python3 -c "import sys,json; d=json.load(sys.stdin); print('plan:', d.get('plan',{}).get('goalTitle','?'), 'score:', d.get('diagnosis',{}).get('score','?'))" 2>/dev/null || echo "(skip ai/plan)"
curl -sf -X POST "$API/ai/chat" \
  -H 'Content-Type: application/json' \
  "${AUTH[@]}" \
  -d '{"message":"Сколько могу отложить в месяц?","history":[]}' \
  | python3 -c "import sys,json; d=json.load(sys.stdin); print('reply:', (d.get('reply') or '')[:80])" 2>/dev/null || echo "(skip ai/chat)"

step "7. Кредиты + ипотечный разбор"
curl -sf "${AUTH[@]}" "$API/credits/dashboard" | python3 -m json.tool 2>/dev/null | head -12 || echo "(credits stub)"
curl -sf -X POST "$API/credits/mortgage/analyze" \
  -H 'Content-Type: application/json' \
  "${AUTH[@]}" \
  -d '{"mortgage_amount":12000000}' \
  | python3 -c "import sys,json; d=json.load(sys.stdin); print('approval:', d.get('approval_level'), 'banks:', len(d.get('banks',[])))" 2>/dev/null || echo "(skip mortgage)"

step "8. Insights"
curl -sf "${AUTH[@]}" "$API/insights" | head -c 200; echo

echo
echo "Demo flow complete. UI: http://localhost:3000 (demo code 0000, ?tour=1 для тура)"
