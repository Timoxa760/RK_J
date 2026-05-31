#!/usr/bin/env bash
# Быстрая проверка SMS.ru ключа и цепочки register → SMS
# Usage:
#   cd backend && source .env && SMOKE_PHONE=+7XXXXXXXXXX ./scripts/test_smsru.sh

set -euo pipefail

cd "$(dirname "$0")/.."
if [[ -f .env ]]; then set -a && source .env && set +a; fi

API="${API_BASE:-http://localhost:8000}/api/v1"
PHONE="${SMOKE_PHONE:-}"

echo "==> 1. SMS.ru balance"
BAL=$(curl -sf "https://sms.ru/my/balance?api_id=${SMSRU_API_ID}&json=1")
echo "$BAL" | python3 -m json.tool
STATUS=$(echo "$BAL" | python3 -c "import sys,json; print(json.load(sys.stdin).get('status_code',0))")
if [[ "$STATUS" != "100" ]]; then
  echo "FAIL: invalid SMSRU_API_ID"
  exit 1
fi

echo "==> 2. Gateway health"
curl -sf "${API%/api/v1}/health" | python3 -m json.tool

if [[ -z "$PHONE" ]]; then
  echo "SKIP register: set SMOKE_PHONE=+7... (номер аккаунта SMS.ru)"
  exit 0
fi

echo "==> 3. POST /auth/register → SMS на $PHONE"
REG=$(curl -s -w "\nHTTP:%{http_code}" -X POST "$API/auth/register" \
  -H 'Content-Type: application/json' \
  -d "{\"phone\":\"$PHONE\"}")
HTTP=$(echo "$REG" | tail -1 | cut -d: -f2)
BODY=$(echo "$REG" | sed '$d')
echo "HTTP $HTTP"
echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"

if [[ "$HTTP" != "200" ]]; then
  docker logs backend_user_service 2>&1 | tail -8
  exit 1
fi

echo "OK — проверь SMS на телефоне. Затем:"
echo "  SMOKE_OTP=<код> SMOKE_PHONE=$PHONE ./scripts/smoke_auth_chat.sh"
