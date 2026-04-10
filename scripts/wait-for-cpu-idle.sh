#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/wait-for-cpu-idle.sh

Purpose:
  Wait until host CPU usage drops below the configured resume threshold before
  starting build/test work.

Environment overrides:
  RESOURCE_LOAD_MAX_PERCENT=80
  RESOURCE_LOAD_RESUME_PERCENT=50
  RESOURCE_LOAD_SAMPLE_SECONDS=1
  RESOURCE_LOAD_POLL_SECONDS=5
EOF
}

if [ "${1:-}" = "-h" ] || [ "${1:-}" = "--help" ]; then
  usage
  exit 0
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/load-guard-common.sh"

load_guard_init_defaults
load_guard_validate_config
load_guard_wait_for_idle
