load_guard_init_defaults() {
  : "${RESOURCE_LOAD_MAX_PERCENT:=80}"
  : "${RESOURCE_LOAD_RESUME_PERCENT:=50}"
  : "${RESOURCE_LOAD_SAMPLE_SECONDS:=1}"
  : "${RESOURCE_LOAD_POLL_SECONDS:=5}"
  : "${LOAD_GUARD_TERMINATE_GRACE_SECONDS:=5}"
  : "${LOAD_GUARD_LOG_PREFIX:=load-guard}"
}

load_guard_log() {
  echo "[${LOAD_GUARD_LOG_PREFIX}] $*" >&2
}

load_guard_require_non_negative_integer() {
  local name="$1"
  local value="$2"

  case "$value" in
    ''|*[!0-9]*)
      echo "${name} must be a non-negative integer" >&2
      return 1
      ;;
  esac
}

load_guard_require_positive_integer() {
  local name="$1"
  local value="$2"

  load_guard_require_non_negative_integer "$name" "$value"
  if [ "$value" -le 0 ]; then
    echo "${name} must be greater than 0" >&2
    return 1
  fi
}

load_guard_validate_config() {
  load_guard_require_non_negative_integer "RESOURCE_LOAD_MAX_PERCENT" "$RESOURCE_LOAD_MAX_PERCENT"
  load_guard_require_non_negative_integer "RESOURCE_LOAD_RESUME_PERCENT" "$RESOURCE_LOAD_RESUME_PERCENT"
  load_guard_require_positive_integer "RESOURCE_LOAD_SAMPLE_SECONDS" "$RESOURCE_LOAD_SAMPLE_SECONDS"
  load_guard_require_positive_integer "RESOURCE_LOAD_POLL_SECONDS" "$RESOURCE_LOAD_POLL_SECONDS"
  load_guard_require_positive_integer "LOAD_GUARD_TERMINATE_GRACE_SECONDS" "$LOAD_GUARD_TERMINATE_GRACE_SECONDS"

  if [ "$RESOURCE_LOAD_MAX_PERCENT" -lt 0 ] || [ "$RESOURCE_LOAD_MAX_PERCENT" -gt 100 ]; then
    echo "RESOURCE_LOAD_MAX_PERCENT must be between 0 and 100" >&2
    return 1
  fi

  if [ "$RESOURCE_LOAD_RESUME_PERCENT" -lt 0 ] || [ "$RESOURCE_LOAD_RESUME_PERCENT" -gt 100 ]; then
    echo "RESOURCE_LOAD_RESUME_PERCENT must be between 0 and 100" >&2
    return 1
  fi

  if [ "$RESOURCE_LOAD_RESUME_PERCENT" -gt "$RESOURCE_LOAD_MAX_PERCENT" ]; then
    echo "RESOURCE_LOAD_RESUME_PERCENT cannot be greater than RESOURCE_LOAD_MAX_PERCENT" >&2
    return 1
  fi
}

load_guard_cpu_usage_percent() {
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

load_guard_wait_for_idle() {
  local usage_percent

  usage_percent="$(load_guard_cpu_usage_percent)"
  if [ "$usage_percent" -le "$RESOURCE_LOAD_MAX_PERCENT" ]; then
    return 0
  fi

  load_guard_log "CPU 使用率 ${usage_percent}% 达到阈值 ${RESOURCE_LOAD_MAX_PERCENT}%，等待回落到 ${RESOURCE_LOAD_RESUME_PERCENT}% 后继续..."
  while [ "$usage_percent" -gt "$RESOURCE_LOAD_RESUME_PERCENT" ]; do
    sleep "$RESOURCE_LOAD_POLL_SECONDS"
    usage_percent="$(load_guard_cpu_usage_percent)"
  done

  load_guard_log "CPU 使用率已回落到 ${usage_percent}%，继续执行任务"
}
