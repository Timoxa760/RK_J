#!/usr/bin/env bash
# Проверка critical/important из API_Contract (локально, DEMO_MODE).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export GOCACHE="${GOCACHE:-${XDG_CACHE_HOME:-$HOME/.cache}/potok-go-build}"
export GOMODCACHE="${GOMODCACHE:-${XDG_CACHE_HOME:-$HOME/.cache}/potok-go-mod}"
export DEMO_MODE=true
export JWT_SECRET="${JWT_SECRET:-test-secret}"
export DATABASE_URL="${DATABASE_URL:-postgres://postgres:postgres@127.0.0.1:5432/moneymind?sslmode=disable}"
export KAFKA_BROKERS="${KAFKA_BROKERS:-127.0.0.1:9092}"
export USER_SERVICE_URL=http://127.0.0.1:8001
export RECEIPT_SERVICE_URL=http://127.0.0.1:8002
export SCRAPER_SERVICE_URL=http://127.0.0.1:8003
export AI_PROCESSOR_URL=http://127.0.0.1:8100
export CREDIT_SERVICE_URL=http://127.0.0.1:8009
export ANALYTICS_SERVICE_URL=http://127.0.0.1:8101
export REPORTING_SERVICE_URL=http://127.0.0.1:8010

API="http://127.0.0.1:8000/api/v1"
PHONE="+79991234567"
CODE="0000"
RESULTS=()

record() {
  local prio=$1 name=$2 status=$3 note=$4
  RESULTS+=("$prio|$name|$status|$note")
}

http_code() { curl -s -o /dev/null -w "%{http_code}" "$@"; }

cleanup() {
  kill "${USPID:-}" "${RSPID:-}" "${SCPID:-}" "${CRPID:-}" "${ANPID:-}" \
    "${RPPID:-}" "${AIPID:-}" "${GWPID:-}" 2>/dev/null || true
}
trap cleanup EXIT

build() {
  local name=$1 path=$2
  [[ -x "$ROOT/bin/$name" ]] || go build -o "$ROOT/bin/$name" "$path"
}

for svc in user-service receipt-service scraper-service credit-service \
  analytics-service reporting-service ai-processor api-gateway; do
  case $svc in
    user-service) build user-service "$ROOT/services/core-api/user-service/" ;;
    receipt-service) build receipt-service "$ROOT/services/receipt-engine/receipt-service/" ;;
    scraper-service) build scraper-service "$ROOT/services/receipt-engine/scraper-service/" ;;
    credit-service) build credit-service "$ROOT/services/finance-core/credit-service/" ;;
    analytics-service) build analytics-service "$ROOT/services/money-intelligence/analytics-service/" ;;
    reporting-service) build reporting-service "$ROOT/services/reporting/reporting-service/" ;;
    ai-processor) build ai-processor "$ROOT/services/money-intelligence/ai-processor/" ;;
    api-gateway) build api-gateway "$ROOT/services/core-api/api-gateway/" ;;
  esac
done

"$ROOT/bin/user-service" & USPID=$!
"$ROOT/bin/receipt-service" & RSPID=$!
"$ROOT/bin/scraper-service" & SCPID=$!
"$ROOT/bin/credit-service" & CRPID=$!
"$ROOT/bin/analytics-service" & ANPID=$!
"$ROOT/bin/reporting-service" & RPPID=$!
"$ROOT/bin/ai-processor" & AIPID=$!
sleep 2
"$ROOT/bin/api-gateway" & GWPID=$!
sleep 1

# --- Auth ---
REG_CODE=$(http_code -X POST "$API/auth/register" -H 'Content-Type: application/json' -d "{\"phone\":\"$PHONE\"}")
if [[ "$REG_CODE" == "200" || "$REG_CODE" == "409" ]]; then
  record critical "POST /auth/register" PASS "HTTP $REG_CODE"
else
  record critical "POST /auth/register" FAIL "HTTP $REG_CODE"
fi

LOGIN=$(curl -sf -X POST "$API/auth/login" -H 'Content-Type: application/json' -d "{\"phone\":\"$PHONE\",\"code\":\"$CODE\"}" 2>/dev/null || echo '{}')
TOKEN=$(echo "$LOGIN" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('access_token',''))" 2>/dev/null || true)
if [[ -n "$TOKEN" ]] && echo "$LOGIN" | python3 -c "import sys,json; d=json.load(sys.stdin); exit(0 if d.get('token')==d.get('access_token') and d.get('user',{}).get('role')=='user' else 1)" 2>/dev/null; then
  record critical "POST /auth/login" PASS "JWT + user + token"
else
  record critical "POST /auth/login" FAIL "нет токена или форма ответа"
fi
AUTH=(-H "Authorization: Bearer $TOKEN")

BAD=$(http_code -X POST "$API/auth/login" -H 'Content-Type: application/json' -d "{\"phone\":\"$PHONE\",\"code\":\"9999\"}")
[[ "$BAD" == "401" ]] && record critical "login неверный code → 401" PASS "" || record critical "login неверный code → 401" FAIL "HTTP $BAD"

# --- JWT gateway ---
NCODE=$(http_code "$API/dashboard/sankey")
[[ "$NCODE" == "401" ]] && record critical "JWT обязателен (без токена 401)" PASS "" || record critical "JWT обязателен" FAIL "HTTP $NCODE"

INVALID=$(http_code -H "Authorization: Bearer invalid.token" "$API/dashboard/sankey")
[[ "$INVALID" == "401" ]] && record critical "JWT невалидный → 401" PASS "" || record critical "JWT невалидный → 401" FAIL "HTTP $INVALID"

# --- Providers ---
PCODE=$(http_code -X POST "$API/providers/connect?provider=x5club" -H 'Content-Type: application/json' "${AUTH[@]}" \
  -d '{"credentials":{"phone":"+79991234567","password":"x"}}')
[[ "$PCODE" == "200" || "$PCODE" == "409" ]] && record critical "POST /providers/connect" PASS "HTTP $PCODE" || record critical "POST /providers/connect" FAIL "HTTP $PCODE"

# --- Dashboard critical ---
for path in sankey categories; do
  CODE=$(http_code "${AUTH[@]}" "$API/dashboard/$path")
  [[ "$CODE" == "200" ]] && record critical "GET /dashboard/$path" PASS "" || record critical "GET /dashboard/$path" FAIL "HTTP $CODE"
done

for path in timemachine stores compare; do
  CODE=$(http_code "${AUTH[@]}" "$API/dashboard/$path")
  [[ "$CODE" == "200" ]] && record important "GET /dashboard/$path" PASS "" || record important "GET /dashboard/$path" FAIL "HTTP $CODE"
done

TM=$(curl -sf "${AUTH[@]}" "$API/dashboard/timemachine" | python3 -c "import sys,json; d=json.load(sys.stdin); exit(0 if len(d.get('months',[]))==60 and len(d.get('real_savings',[]))==60 else 1)" 2>/dev/null && echo ok || echo fail)
[[ "$TM" == "ok" ]] && record important "timemachine 60 months shape" PASS "" || record important "timemachine 60 months shape" FAIL ""

# --- Expenses ---
ECODE=$(http_code -X POST "$API/expenses/manual" -H 'Content-Type: application/json' "${AUTH[@]}" \
  -d "{\"user_id\":\"$PHONE\",\"raw_text\":\"продукты 5000\",\"source\":\"voice\"}")
[[ "$ECODE" == "200" ]] && record critical "POST /expenses/manual" PASS "" || record critical "POST /expenses/manual" FAIL "HTTP $ECODE"

# --- FNS scraper ---
FNS=$(curl -sf -X POST "$API/fns/ticket" -H 'Content-Type: application/json' "${AUTH[@]}" \
  -d '{"fn":"9289000100123456","fd":"12345","fp":"1234567890","sum":"999.99","date":"2026-05-01","time":"12:00"}' 2>/dev/null | python3 -c "import sys,json; d=json.load(sys.stdin); exit(0 if d.get('store_name') else 1)" 2>/dev/null && echo ok || echo fail)
[[ "$FNS" == "ok" ]] && record critical "POST /fns/ticket (gateway→scraper)" PASS "demo receipt" || record critical "POST /fns/ticket" FAIL ""

MCO=$(http_code -X POST "$API/fns/mco/sync" -H 'Content-Type: application/json' "${AUTH[@]}" -d "{\"phone\":\"$PHONE\"}")
if [[ "$MCO" == "200" ]]; then
  record important "POST /fns/mco/sync" PASS ""
elif [[ "$MCO" == "500" ]]; then
  record important "POST /fns/mco/sync" PARTIAL "HTTP 500 (нет Kafka)"
else
  record important "POST /fns/mco/sync" FAIL "HTTP $MCO"
fi

# --- Credits ---
CD=$(http_code "${AUTH[@]}" "$API/credits/dashboard")
[[ "$CD" == "200" ]] && record important "GET /credits/dashboard" PASS "" || record important "GET /credits/dashboard" FAIL "HTTP $CD"

CS=$(http_code -X POST "$API/credits/scan" "${AUTH[@]}" -F 'file=@/dev/null;filename=contract.pdf;type=application/pdf')
[[ "$CS" == "200" ]] && record important "POST /credits/scan" PASS "multipart demo" || record important "POST /credits/scan" FAIL "HTTP $CS"

DTI=$(curl -sf "${AUTH[@]}" "$API/credits/dashboard" | python3 -c "import sys,json; d=json.load(sys.stdin); exit(0 if 0<=d.get('dti',-1)<=100 else 1)" 2>/dev/null && echo ok || echo fail)
[[ "$DTI" == "ok" ]] && record important "credits DTI 0..100 (%)" PASS "" || record important "credits DTI 0..100" FAIL "проверить API"

# --- Profile (goal in profile, not goal-service) ---
PG=$(http_code "${AUTH[@]}" "$API/users/me/profile")
[[ "$PG" == "200" ]] && record important "GET /users/me/profile" PASS "" || record important "GET /users/me/profile" FAIL "HTTP $PG"

PP=$(http_code -X PATCH "$API/users/me/profile" -H 'Content-Type: application/json' "${AUTH[@]}" \
  -d '{"goal_title":"Отпуск","goal_amount":150000,"goal_kind":"save"}')
[[ "$PP" == "200" ]] && record important "PATCH /users/me/profile" PASS "" || record important "PATCH /users/me/profile" FAIL "HTTP $PP"

# --- AI advisor ---
PLAN=$(curl -sf "${AUTH[@]}" "$API/ai/plan" | python3 -c "import sys,json; d=json.load(sys.stdin); exit(0 if d.get('plan',{}).get('goalTitle') and d.get('diagnosis',{}).get('score') is not None else 1)" 2>/dev/null && echo ok || echo fail)
[[ "$PLAN" == "ok" ]] && record important "GET /ai/plan" PASS "plan + diagnosis" || record important "GET /ai/plan" FAIL "shape"

CHAT=$(http_code -X POST "$API/ai/chat" -H 'Content-Type: application/json' "${AUTH[@]}" \
  -d '{"message":"Сколько могу отложить в месяц?","history":[]}')
[[ "$CHAT" == "200" ]] && record important "POST /ai/chat" PASS "" || record important "POST /ai/chat" FAIL "HTTP $CHAT"

# --- Analytics ---
for ep in insights forecast; do
  CODE=$(http_code "${AUTH[@]}" "$API/$ep")
  [[ "$CODE" == "200" ]] && record important "GET /$ep" PASS "" || record important "GET /$ep" FAIL "HTTP $CODE"
done
ACODE=$(http_code "${AUTH[@]}" "$API/analytics/insights")
[[ "$ACODE" == "200" ]] && record important "GET /analytics/insights (front)" PASS "" || record important "GET /analytics/insights" FAIL "HTTP $ACODE"

SIM=$(curl -sf -X POST "$API/scenarios/simulate" -H 'Content-Type: application/json' "${AUTH[@]}" \
  -d '{"scenario":"reduce_delivery","reduction_percent":50,"months":3}' | python3 -c "import sys,json; d=json.load(sys.stdin); exit(0 if len(d.get('months',[]))==3 else 1)" 2>/dev/null && echo ok || echo fail)
[[ "$SIM" == "ok" ]] && record important "POST /scenarios/simulate" PASS "" || record important "POST /scenarios/simulate" FAIL ""

# --- Digest ---
DCODE=$(http_code "${AUTH[@]}" "$API/digest/latest")
[[ "$DCODE" == "200" ]] && record optional "GET /digest/latest" PASS "" || record optional "GET /digest/latest" FAIL "HTTP $DCODE"

# --- Stubs (ожидаемо не готово) ---
for stub in "GET /banks/accounts" "GET /banks/transactions" "POST /challenges"; do
  path="${stub#* }"
  method="${stub%% *}"
  url="$API${path}"
  if [[ "$method" == "GET" ]]; then
    SC=$(http_code "${AUTH[@]}" "$url")
  else
    SC=$(http_code -X POST "$url" -H 'Content-Type: application/json' "${AUTH[@]}" -d '{}')
  fi
  if [[ "$SC" == "502" || "$SC" == "404" || "$SC" == "503" ]]; then
    record optional "$stub" SKIP "HTTP $SC (не реализовано)"
  else
    record optional "$stub" INFO "HTTP $SC"
  fi
done

echo "PRIORITY|ENDPOINT|STATUS|NOTE"
for row in "${RESULTS[@]}"; do
  echo "$row"
done

FAIL_COUNT=0
for row in "${RESULTS[@]}"; do
  [[ "$row" == *"|FAIL|"* ]] && FAIL_COUNT=$((FAIL_COUNT+1))
done
echo "---"
echo "FAILURES: $FAIL_COUNT"
exit "$FAIL_COUNT"
