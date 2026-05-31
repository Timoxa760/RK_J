#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PATH_IN_APP="${1:-/dashboard}"
SCALE="${SCREENCAST_SCALE:-2}"
PORT="${NUXT_PORT:-3000}"
BASE="http://localhost:${PORT}"

QUERY="path=${PATH_IN_APP}&scale=${SCALE}&clean=1"
URL="${BASE}/screencast.html?${QUERY}"

if ! curl -sf "${BASE}/" >/dev/null 2>&1; then
  echo "Сначала запустите dev-сервер:"
  echo "  cd frontend && npm run dev"
  exit 1
fi

FRAME_W=444
FRAME_H=921
WIN_W=$((FRAME_W * SCALE))
WIN_H=$((FRAME_H * SCALE))

CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
EDGE="/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"

open_window() {
  local browser="$1"
  "$browser" \
    --new-window \
    --app="${URL}" \
    --window-size="${WIN_W},${WIN_H}" \
    --window-position=120,60 \
    >/dev/null 2>&1 &
}

if [[ -x "$CHROME" ]]; then
  open_window "$CHROME"
elif [[ -x "$EDGE" ]]; then
  open_window "$EDGE"
else
  open "${URL}"
fi

cat <<EOF

Окно для записи открыто (${WIN_W}×${WIN_H}px, scale=${SCALE}).

Запись:
  • QuickTime → Файл → Новая запись экрана → это окно
  • OBS → Источник «Захват окна» → это окно

Настройки (перед записью):
  • ${BASE}/screencast.html?path=/onboarding  — другой экран
  • SCREENCAST_SCALE=3 ./scripts/open-screencast.sh  — крупнее (1332×2763)

Если видите логин — сначала войдите: ${BASE}/login

EOF
