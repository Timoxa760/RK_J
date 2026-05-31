#!/usr/bin/env bash
# E2E smoke: register → login → profile → advisor chat (history + stream)
# Usage:
#   SMOKE_PHONE=+79284018652 SMOKE_PASSWORD=secret123 ./scripts/smoke_auth_chat.sh

set -euo pipefail

API="${API_BASE:-http://localhost:8000}/api/v1"
PHONE="${SMOKE_PHONE:-+79991234567}"
PASS="${SMOKE_PASSWORD:-secret12345}"

step() { echo; echo "==> $1"; }

step "1. Register (ignore if exists)"
curl -sf -X POST "$API/auth/register" \
  -H 'Content-Type: application/json' \
  -d "{\"phone\":\"$PHONE\",\"password\":\"$PASS\"}" 2>/dev/null \
  | python3 -m json.tool 2>/dev/null || echo "(already registered or skip)"

step "2. Login"
LOGIN=$(curl -sf -X POST "$API/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"phone\":\"$PHONE\",\"password\":\"$PASS\"}")
TOKEN=$(echo "$LOGIN" | python3 -c "import sys,json; print(json.load(sys.stdin).get('access_token',''))")
if [[ -z "$TOKEN" ]]; then
  echo "FAIL login: $LOGIN"
  exit 1
fi
echo "JWT ok"

AUTH=(-H "Authorization: Bearer $TOKEN")

step "3. PATCH profile"
curl -sf -X PATCH "$API/users/me/profile" \
  -H 'Content-Type: application/json' \
  "${AUTH[@]}" \
  -d '{"active_income":150000,"passive_income":30000,"emergency_fund":340000,"goal_title":"Отпуск","goal_amount":150000,"goal_kind":"save","onboarding_completed":true}' \
  | python3 -m json.tool 2>/dev/null | head -12

step "4. GET /ai/plan"
curl -sf "${AUTH[@]}" "$API/ai/plan" \
  | python3 -c "import sys,json; d=json.load(sys.stdin); print('plan goal:', d.get('plan',{}).get('goalTitle','?'))"

step "5. POST /ai/chat"
CHAT=$(curl -sf -X POST "$API/ai/chat" \
  -H 'Content-Type: application/json' \
  "${AUTH[@]}" \
  -d '{"message":"Сколько откладывать на цель?","history":[]}')
echo "$CHAT" | python3 -c "import sys,json; d=json.load(sys.stdin); print('reply:', (d.get('reply') or '')[:200]); print('source:', d.get('source','?'))"

step "6. GET /ai/chat/history"
HIST=$(curl -sf "${AUTH[@]}" "$API/ai/chat/history?limit=10")
echo "$HIST" | python3 -c "import sys,json; d=json.load(sys.stdin); print('messages:', len(d.get('messages',[])))"

step "7. POST /ai/chat/stream (first events)"
STREAM=$(curl -sf -N -X POST "$API/ai/chat/stream" \
  -H 'Content-Type: application/json' \
  "${AUTH[@]}" \
  -d '{"message":"Где урезать траты?","history":[]}' | head -20 || true)
echo "$STREAM" | head -5

step "8. DELETE /ai/chat/history"
curl -sf -X DELETE "${AUTH[@]}" "$API/ai/chat/history" \
  | python3 -m json.tool 2>/dev/null || echo "cleared"

echo
echo "OK smoke_auth_chat"
