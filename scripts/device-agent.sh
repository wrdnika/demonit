#!/usr/bin/env bash
# Demonit device agent for macOS / Linux (bash + curl + python3 helper).
# Prefer scripts/device_agent.py when Python is available.
#
# Usage:
#   chmod +x device-agent.sh
#   ./device-agent.sh --device-id <UUID>
#   ./device-agent.sh --device-id <UUID> --api-base http://192.168.1.10:8080

set -euo pipefail

DEVICE_ID=""
API_BASE="${DEMONIT_API_BASE:-http://localhost:8080}"
API_KEY="${DEMONIT_DEVICE_API_KEY:-dev-device-key-change-me}"
INTERVAL=10

while [[ $# -gt 0 ]]; do
  case "$1" in
    --device-id) DEVICE_ID="$2"; shift 2 ;;
    --api-base) API_BASE="$2"; shift 2 ;;
    --api-key) API_KEY="$2"; shift 2 ;;
    --interval) INTERVAL="$2"; shift 2 ;;
    *) echo "Unknown arg: $1" >&2; exit 1 ;;
  esac
done

if [[ -z "$DEVICE_ID" ]]; then
  echo "Usage: $0 --device-id <UUID> [--api-base URL] [--api-key KEY] [--interval SEC]" >&2
  exit 1
fi

if ! command -v curl >/dev/null 2>&1; then
  echo "curl is required" >&2
  exit 1
fi

# Delegate metrics + POST to the Python agent when possible (best accuracy).
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if command -v python3 >/dev/null 2>&1; then
  exec python3 "$SCRIPT_DIR/device_agent.py" \
    --device-id "$DEVICE_ID" \
    --api-base "$API_BASE" \
    --api-key "$API_KEY" \
    --interval "$INTERVAL"
fi

echo "python3 not found; using curl-only fallback (cpu/ram mocked)."
echo "Demonit agent started"
echo "  device_id : $DEVICE_ID"
echo "  api       : $API_BASE"
echo "  interval  : ${INTERVAL}s"
echo "Ctrl+C to stop."
echo

HOST="$(hostname 2>/dev/null || echo unknown)"
OS="$(uname -s 2>/dev/null | tr '[:upper:]' '[:lower:]')"

while true; do
  CPU=$((RANDOM % 40 + 10))
  RAM=$((RANDOM % 40 + 30))
  TS="$(date +%H:%M:%S)"
  BODY=$(cat <<EOF
{"device_id":"$DEVICE_ID","cpu_usage":$CPU,"ram_usage":$RAM,"status_payload":{"hostname":"$HOST","os":"$OS","agent":"bash-curl"}}
EOF
)
  HTTP_CODE=$(curl -sS -o /tmp/demonit_hb.json -w "%{http_code}" \
    -X POST "$API_BASE/api/v1/heartbeat" \
    -H "Content-Type: application/json" \
    -H "X-Device-API-Key: $API_KEY" \
    -d "$BODY" || true)

  if [[ "$HTTP_CODE" == "200" ]]; then
    echo "[$TS] OK  cpu=${CPU}% ram=${RAM}%  http=$HTTP_CODE"
  else
    echo "[$TS] FAIL http=$HTTP_CODE  $(cat /tmp/demonit_hb.json 2>/dev/null || true)" >&2
  fi
  sleep "$INTERVAL"
done
