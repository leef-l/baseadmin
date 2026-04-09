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
  RESOURCE_LOAD_RESUME_PERCENT=60
  RESOURCE_LOAD_SAMPLE_SECONDS=1
  RESOURCE_LOAD_POLL_SECONDS=5
EOF
}

if [ "${1:-}" = "-h" ] || [ "${1:-}" = "--help" ]; then
  usage
  exit 0
fi

: "${RESOURCE_LOAD_MAX_PERCENT:=80}"
: "${RESOURCE_LOAD_RESUME_PERCENT:=60}"
: "${RESOURCE_LOAD_SAMPLE_SECONDS:=1}"
: "${RESOURCE_LOAD_POLL_SECONDS:=5}"

if [ "$RESOURCE_LOAD_MAX_PERCENT" -lt 0 ] || [ "$RESOURCE_LOAD_MAX_PERCENT" -gt 100 ]; then
  echo "RESOURCE_LOAD_MAX_PERCENT must be between 0 and 100" >&2
  exit 1
fi

if [ "$RESOURCE_LOAD_RESUME_PERCENT" -lt 0 ] || [ "$RESOURCE_LOAD_RESUME_PERCENT" -gt 100 ]; then
  echo "RESOURCE_LOAD_RESUME_PERCENT must be between 0 and 100" >&2
  exit 1
fi

if [ "$RESOURCE_LOAD_RESUME_PERCENT" -gt "$RESOURCE_LOAD_MAX_PERCENT" ]; then
  echo "RESOURCE_LOAD_RESUME_PERCENT cannot be greater than RESOURCE_LOAD_MAX_PERCENT" >&2
  exit 1
fi

cpu_usage_percent() {
  local user1 nice1 system1 idle1 iowait1 irq1 softirq1 steal1 idle_total1 total1
  local user2 nice2 system2 idle2 iowait2 irq2 softirq2 steal2 idle_total2 total2
  local total_delta idle_delta busy_delta

  read -r _ user1 nice1 system1 idle1 iowait1 irq1 softirq1 steal1 _ < /proc/stat
  idle_total1=$((idle1 + iowait1))
  total1=$((user1 + nice1 + system1 + idle_total1 + irq1 + softirq1 + steal1))

  sleep "$RESOURCE_LOAD_SAMPLE_SECONDS"

  read -r _ user2 nice2 system2 idle2 iowait2 irq2 softirq2 steal2 _ < /proc/stat
  idle_total2=$((idle2 + iowait2))
  total2=$((user2 + nice2 + system2 + idle_total2 + irq2 + softirq2 + steal2))

  total_delta=$((total2 - total1))
  idle_delta=$((idle_total2 - idle_total1))
  busy_delta=$((total_delta - idle_delta))

  if [ "$total_delta" -le 0 ]; then
    echo 0
    return
  fi

  echo $(((busy_delta * 100 + total_delta / 2) / total_delta))
}

usage_percent="$(cpu_usage_percent)"
if [ "$usage_percent" -le "$RESOURCE_LOAD_MAX_PERCENT" ]; then
  exit 0
fi

echo "[load-guard] CPU 使用率 ${usage_percent}% 超过阈值 ${RESOURCE_LOAD_MAX_PERCENT}%，等待降到 ${RESOURCE_LOAD_RESUME_PERCENT}% 后继续..." >&2
while [ "$usage_percent" -gt "$RESOURCE_LOAD_RESUME_PERCENT" ]; do
  sleep "$RESOURCE_LOAD_POLL_SECONDS"
  usage_percent="$(cpu_usage_percent)"
done

echo "[load-guard] CPU 使用率已回落到 ${usage_percent}%，继续执行任务" >&2
