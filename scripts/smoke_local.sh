#!/usr/bin/env bash
# Локальный smoke: user + receipt + gateway (без полного compose).
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

API="http://127.0.0.1:8000/api/v1"

cleanup() {
  kill "${USPID:-}" "${RSPID:-}" "${GWPID:-}" 2>/dev/null || true
}
trap cleanup EXIT

[[ -x "$ROOT/bin/user-service" ]] || go build -o "$ROOT/bin/user-service" "$ROOT/services/core-api/user-service/"
[[ -x "$ROOT/bin/receipt-service" ]] || go build -o "$ROOT/bin/receipt-service" "$ROOT/services/receipt-engine/receipt-service/"
[[ -x "$ROOT/bin/api-gateway" ]] || go build -o "$ROOT/bin/api-gateway" "$ROOT/services/core-api/api-gateway/"

"$ROOT/bin/user-service" &
USPID=$!
"$ROOT/bin/receipt-service" &
RSPID=$!
sleep 2
"$ROOT/bin/api-gateway" &
GWPID=$!
sleep 1

echo "==> health gateway"
curl -sf "http://127.0.0.1:8000/health" | grep -q ok

echo "==> login via gateway"
TOKEN=$(curl -sf -X POST "$API/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"phone":"+79991234567","code":"0000"}' | \
  python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])")

echo "==> dashboard sankey via gateway"
NODES=$(curl -sf -H "Authorization: Bearer $TOKEN" "$API/dashboard/sankey" | \
  python3 -c "import sys,json; print(len(json.load(sys.stdin).get('nodes',[])))")
echo "sankey nodes: $NODES"
[[ "$NODES" -ge 3 ]]

echo "==> dashboard timemachine via gateway"
curl -sf -H "Authorization: Bearer $TOKEN" "$API/dashboard/timemachine" | \
  python3 -c "import sys,json; d=json.load(sys.stdin); assert len(d.get('months',[]))==60"

echo "SMOKE OK"
