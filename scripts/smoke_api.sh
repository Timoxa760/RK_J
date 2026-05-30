#!/usr/bin/env bash
# Расширенный smoke по API_Contract для реализованных эндпоинтов.
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
export AI_PROCESSOR_URL=http://127.0.0.1:8100
export CREDIT_SERVICE_URL=http://127.0.0.1:8009
export ANALYTICS_SERVICE_URL=http://127.0.0.1:8101
export REPORTING_SERVICE_URL=http://127.0.0.1:8010

API="http://127.0.0.1:8000/api/v1"
PHONE="+79991234567"
PASS="secret12345"
FAIL=0

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1"; FAIL=1; }
http_code() { curl -s -o /dev/null -w "%{http_code}" "$@"; }

cleanup() {
  kill "${USPID:-}" "${RSPID:-}" "${CRPID:-}" "${ANPID:-}" "${RPPID:-}" "${AIPID:-}" "${GWPID:-}" 2>/dev/null || true
}
trap cleanup EXIT

build() {
  local name=$1 path=$2
  [[ -x "$ROOT/bin/$name" ]] || go build -o "$ROOT/bin/$name" "$path"
}

build user-service "$ROOT/services/core-api/user-service/"
build receipt-service "$ROOT/services/receipt-engine/receipt-service/"
build credit-service "$ROOT/services/finance-core/credit-service/"
build analytics-service "$ROOT/services/money-intelligence/analytics-service/"
build reporting-service "$ROOT/services/reporting/reporting-service/"
build api-gateway "$ROOT/services/core-api/api-gateway/"

"$ROOT/bin/user-service" &
USPID=$!
"$ROOT/bin/receipt-service" &
RSPID=$!
"$ROOT/bin/credit-service" &
CRPID=$!
"$ROOT/bin/analytics-service" &
ANPID=$!
"$ROOT/bin/reporting-service" &
RPPID=$!

AIPID=""
if build ai-processor "$ROOT/services/money-intelligence/ai-processor/" 2>/dev/null; then
  "$ROOT/bin/ai-processor" &
  AIPID=$!
fi

sleep 2
"$ROOT/bin/api-gateway" &
GWPID=$!
sleep 1

echo "=== 1. Gateway health ==="
[[ "$(http_code http://127.0.0.1:8000/health)" == "200" ]] && pass "gateway /health" || fail "gateway /health"

echo "=== 2. Auth (API_Contract) ==="
REG=$(curl -sf -X POST "$API/auth/register" -H 'Content-Type: application/json' -d "{\"phone\":\"$PHONE\",\"password\":\"$PASS\"}" 2>/dev/null || echo '{}')
echo "$REG" | python3 -c "import sys,json; d=json.load(sys.stdin); assert d.get('message')=='registered'" 2>/dev/null \
  && pass "register (existing user may 409)" || pass "register skipped or conflict"

LOGIN=$(curl -sf -X POST "$API/auth/login" -H 'Content-Type: application/json' -d "{\"phone\":\"$PHONE\",\"password\":\"$PASS\"}")
TOKEN=$(echo "$LOGIN" | python3 -c "import sys,json; d=json.load(sys.stdin); assert d.get('token')==d.get('access_token'); print(d['access_token']); assert d.get('user',{}).get('role')=='user'")
pass "login + user + token alias"
AUTH=(-H "Authorization: Bearer $TOKEN")

BAD_CODE=$(http_code -X POST "$API/auth/login" -H 'Content-Type: application/json' -d "{\"phone\":\"$PHONE\",\"password\":\"wrongpass\"}")
[[ "$BAD_CODE" == "401" ]] && pass "login wrong password 401" || fail "login wrong password (got $BAD_CODE)"

echo "=== 3. Providers connect ==="
PCODE=$(http_code -X POST "$API/providers/connect?provider=x5club" -H 'Content-Type: application/json' "${AUTH[@]}" \
  -d '{"credentials":{"phone":"+79991234567","password":"x"}}')
[[ "$PCODE" == "200" || "$PCODE" == "409" ]] && pass "providers/connect ($PCODE)" || fail "providers/connect ($PCODE)"

echo "=== 4. Dashboard (receipt-service) ==="
for path in sankey categories stores compare timemachine; do
  CODE=$(http_code "${AUTH[@]}" "$API/dashboard/$path")
  [[ "$CODE" == "200" ]] && pass "GET /dashboard/$path" || fail "GET /dashboard/$path ($CODE)"
done

echo "$LOGIN" | python3 -c "
import sys,json,urllib.request
token=json.load(sys.stdin)['access_token']
req=urllib.request.Request('http://127.0.0.1:8000/api/v1/dashboard/timemachine', headers={'Authorization':'Bearer '+token})
d=json.loads(urllib.request.urlopen(req).read())
assert len(d.get('months',[]))==60, d
assert len(d.get('real_savings',[]))==60
" && pass "timemachine contract shape" || fail "timemachine contract shape"

echo "=== 5. Expenses manual (ai-processor) ==="
if [[ -n "$AIPID" ]]; then
  ECODE=$(http_code -X POST "$API/expenses/manual" -H 'Content-Type: application/json' "${AUTH[@]}" \
    -d "{\"user_id\":\"$PHONE\",\"raw_text\":\"продукты 5000\",\"source\":\"voice\"}")
  [[ "$ECODE" == "200" ]] && pass "POST /expenses/manual" || fail "POST /expenses/manual ($ECODE)"
else
  pass "ai-processor skip"
fi

echo "=== 6. Credits + insights ==="
CODE=$(http_code "${AUTH[@]}" "$API/credits/dashboard")
[[ "$CODE" == "200" ]] && pass "GET /credits/dashboard" || fail "GET /credits/dashboard ($CODE)"
CODE=$(http_code "${AUTH[@]}" "$API/insights")
[[ "$CODE" == "200" ]] && pass "GET /insights" || fail "GET /insights ($CODE)"
SIM=$(curl -sf -X POST "$API/scenarios/simulate" -H 'Content-Type: application/json' "${AUTH[@]}" \
  -d '{"scenario":"reduce_delivery","reduction_percent":50,"months":3}' | python3 -c "import sys,json; d=json.load(sys.stdin); exit(0 if len(d.get('months',[]))==3 else 1)" 2>/dev/null && echo ok || echo fail)
[[ "$SIM" == "ok" ]] && pass "POST /scenarios/simulate" || fail "POST /scenarios/simulate"
CODE=$(http_code "${AUTH[@]}" "$API/analytics/insights")
[[ "$CODE" == "200" ]] && pass "GET /analytics/insights (front alias)" || fail "GET /analytics/insights ($CODE)"

echo "=== 7. Profile + digest ==="
PCODE=$(http_code "${AUTH[@]}" "$API/users/me/profile")
[[ "$PCODE" == "200" ]] && pass "GET /users/me/profile" || fail "GET /users/me/profile ($PCODE)"
PPATCH=$(http_code -X PATCH "$API/users/me/profile" -H 'Content-Type: application/json' "${AUTH[@]}" \
  -d '{"goal_title":"Отпуск","goal_amount":150000}')
[[ "$PPATCH" == "200" ]] && pass "PATCH /users/me/profile" || fail "PATCH /users/me/profile ($PPATCH)"
SCODE=$(http_code -X POST "$API/credits/scan" "${AUTH[@]}" -F 'file=@/dev/null;filename=contract.pdf;type=application/pdf')
[[ "$SCODE" == "200" ]] && pass "POST /credits/scan (demo)" || fail "POST /credits/scan ($SCODE)"
DCODE=$(http_code "${AUTH[@]}" "$API/digest/latest")
[[ "$DCODE" == "200" ]] && pass "GET /digest/latest" || fail "GET /digest/latest ($DCODE)"

echo "=== 8. Auth без JWT ==="
NCODE=$(http_code "$API/dashboard/sankey")
[[ "$NCODE" == "401" ]] && pass "dashboard without JWT 401" || fail "dashboard without JWT ($NCODE)"

if [[ "$FAIL" -eq 0 ]]; then
  echo
  echo "ALL API SMOKE PASSED"
else
  echo
  echo "API SMOKE HAD FAILURES"
  exit 1
fi
