#!/usr/bin/env bash
# Диагностика ответов /ai/chat: source, reply, blocks, подозрительные артефакты LLM.
# Минимум токенов — 2 коротких сообщения по умолчанию.
#
# Usage:
#   ./scripts/advisor_chat_probe.sh
#   PROBE_MSGS='Сколько откладывать?|Где урезать?' ./scripts/advisor_chat_probe.sh
#   PROBE_CLEAR=1 ./scripts/advisor_chat_probe.sh
#
# Env: API_BASE, SMOKE_PHONE, SMOKE_PASSWORD, PROBE_MSGS, PROBE_CLEAR

set -euo pipefail

API="${API_BASE:-http://localhost:8000}/api/v1"
PHONE="${SMOKE_PHONE:-+79991234567}"
PASS="${SMOKE_PASSWORD:-secret12345}"
IFS='|' read -r -a MSGS <<< "${PROBE_MSGS:-Сколько откладывать?|Где урезать?}"

step() { echo; echo "==> $1"; }

analyze_response() {
  local label="$1"
  local payload="$2"
  PAYLOAD="$payload" python3 - "$label" <<'PY'
import json, os, re, sys

label = sys.argv[1]
raw = os.environ["PAYLOAD"]
d = json.loads(raw)

KNOWN_BROKEN = [
    "доход ы", "расход ы", "доход а", "расход а", "расход ов", "доход ов",
    "по это му", "за по лнении", "де по зит", "го до вых", "кр из ис",
    "пред по чтительнее", "месяц ев", "месяцона", "увасуже", "заменын",
]
RE_SUFFIX = re.compile(r"(?i)\b([а-яё]{4,}) ([а-яё]{1,3})\b")
RE_CHAIN = re.compile(
    r"(?i)(?:^|[\s(«\"'—-])([а-яё]{1,3}(?: [а-яё]{1,3}){2,})(?:[\s,.!?;:»)\]\"'—-]|$)"
)
STOP = {
    "по", "на", "за", "от", "до", "не", "ни", "из", "в", "и", "или", "как", "что", "это",
    "уже", "ещё", "еще", "для", "без", "при", "над", "под", "про", "там", "тут", "все",
    "всё", "нет", "есть", "она", "они", "мне", "вас", "нас", "ему", "ей", "им", "бы", "ли",
    "же", "он", "мы", "вы", "ты", "я", "д", "г", "т", "п",
}


def collect_text(payload):
    parts = []
    if payload.get("title"):
        parts.append(payload["title"])
    if payload.get("reply"):
        parts.append(payload["reply"])
    for block in payload.get("blocks") or []:
        if block.get("text"):
            parts.append(block["text"])
        parts.extend(block.get("items") or [])
    return "\n".join(parts)


text = collect_text(d)
issues = []
low = text.lower()
for frag in KNOWN_BROKEN:
    if frag in low:
        issues.append(f"known:{frag!r}")

for stem, tail in RE_SUFFIX.findall(text)[:24]:
    if tail.lower() in STOP:
        continue
    issues.append(f"suffix:{stem!r}+{tail!r}")

for match in RE_CHAIN.finditer(text):
    chain = match.group(1)
    if chain.lower().replace(" ", "") in {"итд", "тп"}:
        continue
    issues.append(f"chain:{chain!r}")

print(f"--- {label} ---")
print("source:", d.get("source", "?"))
print("id:", d.get("id", ""))
blocks = d.get("blocks")
print("blocks:", "none" if not blocks else len(blocks))
if blocks:
    for i, block in enumerate(blocks[:4]):
        preview = (block.get("text") or " / ".join(block.get("items") or []))[:120]
        print(f"  block[{i}] type={block.get('type')} {preview!r}")

reply = (d.get("reply") or "").strip()
print("reply_len:", len(reply))
print("reply_preview:")
print(reply[:500] + ("…" if len(reply) > 500 else ""))

if issues:
    print("ISSUES (" + str(len(issues)) + "):")
    for item in issues[:12]:
        print(" ", item)
    if len(issues) > 12:
        print("  …", len(issues) - 12, "more")
else:
    print("ISSUES: none detected")
print()
print(len(issues))
PY
}

step "Login ($PHONE)"
LOGIN=$(curl -sf -X POST "$API/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"phone\":\"$PHONE\",\"password\":\"$PASS\"}") || {
  echo "FAIL: login — API недоступен? ($API)"
  exit 1
}
TOKEN=$(echo "$LOGIN" | python3 -c "import sys,json; print(json.load(sys.stdin).get('access_token',''))")
if [[ -z "$TOKEN" ]]; then
  echo "FAIL login: $LOGIN"
  exit 1
fi
AUTH=(-H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json')

if [[ "${PROBE_CLEAR:-0}" == "1" ]]; then
  step "Clear chat history"
  curl -sf -X DELETE "${AUTH[@]}" "$API/ai/chat/history" >/dev/null || true
fi

total_issues=0
history_json='[]'

for i in "${!MSGS[@]}"; do
  msg="${MSGS[$i]}"
  step "POST /ai/chat #$((i + 1)): ${msg}"
  body=$(MSG="$msg" HIST="$history_json" python3 <<'PY'
import json, os
print(json.dumps({"message": os.environ["MSG"], "history": json.loads(os.environ["HIST"])}, ensure_ascii=False))
PY
)
  resp=$(curl -sf -X POST "$API/ai/chat" "${AUTH[@]}" -d "$body")
  out=$(analyze_response "chat #$((i + 1))" "$resp")
  count=$(echo "$out" | tail -1)
  echo "$out" | sed '$d'
  total_issues=$((total_issues + count))

  history_json=$(MSG="$msg" RESP="$resp" python3 <<'PY'
import json, os
resp = json.loads(os.environ["RESP"])
assistant = (resp.get("reply") or "")[:800]
print(json.dumps([
    {"role": "user", "content": os.environ["MSG"]},
    {"role": "assistant", "content": assistant},
], ensure_ascii=False))
PY
)
done

step "Summary"
echo "messages: ${#MSGS[@]}"
echo "total_issues: $total_issues"

if [[ "$total_issues" -gt 0 ]]; then
  echo
  echo "Hint: если go test ./internal/advisor/ -run RepairSplit проходит, а ISSUES остаются —"
  echo "      пересоберите ai-processor:"
  echo "      cd backend && docker compose build ai-processor && docker compose up -d ai-processor"
  exit 2
fi

echo "OK advisor_chat_probe"
