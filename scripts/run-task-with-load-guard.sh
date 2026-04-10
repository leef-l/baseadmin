#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/run-task-with-load-guard.sh [options] -- <command> [args...]

Purpose:
  Wait for CPU usage to fall below the configured threshold before starting a
  task, then pause and resume the whole task automatically while it is running
  so a 2c4g host is not overwhelmed.

Options:
  --nice N
  --env KEY=VALUE
  --systemd-env KEY=VALUE
  --systemd-prop NAME=VALUE
  --systemd-unit-prefix PREFIX
  --local-only
  --no-ionice

Environment overrides:
  RESOURCE_LOAD_MAX_PERCENT=80
  RESOURCE_LOAD_RESUME_PERCENT=50
  RESOURCE_LOAD_SAMPLE_SECONDS=1
  RESOURCE_LOAD_POLL_SECONDS=5
  LOAD_GUARD_TERMINATE_GRACE_SECONDS=5
  LOAD_GUARD_LOG_PREFIX=load-guard
EOF
}

if [ "$#" -eq 0 ]; then
  usage >&2
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/load-guard-common.sh"

load_guard_init_defaults

nice_value=""
use_ionice=1
local_only=0
systemd_unit_prefix="guarded-task"
declare -a child_env=()
declare -a systemd_env=()
declare -a systemd_props=()

require_assignment() {
  local flag="$1"
  local value="$2"

  case "$value" in
    *=*)
      ;;
    *)
      echo "${flag} expects KEY=VALUE" >&2
      exit 1
      ;;
  esac
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --nice)
      nice_value="${2:-}"
      shift 2
      ;;
    --env)
      require_assignment "$1" "${2:-}"
      child_env+=("${2}")
      shift 2
      ;;
    --systemd-env)
      require_assignment "$1" "${2:-}"
      systemd_env+=("${2}")
      shift 2
      ;;
    --systemd-prop)
      require_assignment "$1" "${2:-}"
      systemd_props+=("${2}")
      shift 2
      ;;
    --systemd-unit-prefix)
      systemd_unit_prefix="${2:-}"
      shift 2
      ;;
    --local-only)
      local_only=1
      shift
      ;;
    --no-ionice)
      use_ionice=0
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    --)
      shift
      break
      ;;
    *)
      echo "unknown option: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [ "$#" -eq 0 ]; then
  echo "missing command to run" >&2
  usage >&2
  exit 1
fi

if [ -n "$nice_value" ] && ! printf '%s\n' "$nice_value" | grep -Eq '^-?[0-9]+$'; then
  echo "--nice expects an integer" >&2
  exit 1
fi

load_guard_validate_config

declare -a guarded_command=("$@")

should_use_systemd=0
if [ "$local_only" -ne 1 ] && command -v systemd-run >/dev/null 2>&1 && command -v systemctl >/dev/null 2>&1 && [ -d /run/systemd/system ] && [ "$(id -u)" -eq 0 ]; then
  should_use_systemd=1
fi

runner_pid=""
monitor_pid=""
child_unit=""
child_pgid=""
is_paused=0

is_runner_alive() {
  [ -n "$runner_pid" ] && kill -0 "$runner_pid" >/dev/null 2>&1
}

is_child_alive() {
  if [ "$should_use_systemd" -eq 1 ]; then
    is_runner_alive
    return
  fi

  [ -n "$child_pgid" ] && kill -0 -- "-${child_pgid}" >/dev/null 2>&1
}

build_unit_name() {
  local prefix="$1"
  local sanitized

  sanitized="$(printf '%s' "$prefix" | tr -cs 'A-Za-z0-9_.-' '-')"
  sanitized="${sanitized#-}"
  sanitized="${sanitized%-}"
  if [ -z "$sanitized" ]; then
    sanitized="guarded-task"
  fi

  printf '%s-%s-%s' "$sanitized" "$$" "$(date +%s%N)"
}

pause_child() {
  if [ "$is_paused" -eq 1 ] || ! is_child_alive; then
    return 0
  fi

  if [ "$should_use_systemd" -eq 1 ]; then
    systemctl kill --kill-who=all --signal=STOP "$child_unit" >/dev/null 2>&1 || true
  else
    kill -STOP -- "-${child_pgid}" >/dev/null 2>&1 || kill -STOP "$runner_pid" >/dev/null 2>&1 || true
  fi
  is_paused=1
}

resume_child() {
  if [ "$is_paused" -ne 1 ]; then
    return 0
  fi

  if [ "$should_use_systemd" -eq 1 ]; then
    systemctl kill --kill-who=all --signal=CONT "$child_unit" >/dev/null 2>&1 || true
  else
    kill -CONT -- "-${child_pgid}" >/dev/null 2>&1 || kill -CONT "$runner_pid" >/dev/null 2>&1 || true
  fi
  is_paused=0
}

stop_monitor() {
  if [ -n "$monitor_pid" ] && kill -0 "$monitor_pid" >/dev/null 2>&1; then
    kill "$monitor_pid" >/dev/null 2>&1 || true
    wait "$monitor_pid" >/dev/null 2>&1 || true
  fi
  monitor_pid=""
}

send_signal_to_child() {
  local signal_name="$1"

  if [ "$is_paused" -eq 1 ]; then
    resume_child
  fi

  if ! is_child_alive; then
    return 0
  fi

  if [ "$should_use_systemd" -eq 1 ]; then
    systemctl kill --kill-who=all --signal="$signal_name" "$child_unit" >/dev/null 2>&1 || true
  else
    kill "-${signal_name}" -- "-${child_pgid}" >/dev/null 2>&1 || kill "-${signal_name}" "$runner_pid" >/dev/null 2>&1 || true
  fi
}

wait_for_child_exit() {
  local remaining="$1"

  while is_child_alive; do
    if [ "$remaining" -le 0 ]; then
      return 1
    fi
    sleep 1
    remaining=$((remaining - 1))
  done
}

reap_runner() {
  if [ -n "$runner_pid" ]; then
    wait "$runner_pid" >/dev/null 2>&1 || true
  fi
}

shutdown_child() {
  local signal_name="$1"

  if ! is_child_alive; then
    reap_runner
    return 0
  fi

  send_signal_to_child "$signal_name"
  if wait_for_child_exit "$LOAD_GUARD_TERMINATE_GRACE_SECONDS"; then
    reap_runner
    return 0
  fi

  load_guard_log "子任务在 ${LOAD_GUARD_TERMINATE_GRACE_SECONDS}s 内未退出，发送 SIGKILL 强制结束"
  send_signal_to_child KILL
  wait_for_child_exit "$LOAD_GUARD_TERMINATE_GRACE_SECONDS" >/dev/null 2>&1 || true
  reap_runner
}

handle_signal() {
  local signal_name="$1"
  local exit_code="$2"

  stop_monitor
  shutdown_child "$signal_name"
  trap - EXIT
  exit "$exit_code"
}

monitor_load() {
  local usage_percent

  while is_child_alive; do
    usage_percent="$(load_guard_cpu_usage_percent)"
    if ! is_child_alive; then
      break
    fi

    if [ "$is_paused" -eq 0 ] && [ "$usage_percent" -ge "$RESOURCE_LOAD_MAX_PERCENT" ]; then
      pause_child
      load_guard_log "CPU 使用率 ${usage_percent}% 达到阈值 ${RESOURCE_LOAD_MAX_PERCENT}%，已暂停任务，等待回落到 ${RESOURCE_LOAD_RESUME_PERCENT}% 后自动继续"
    elif [ "$is_paused" -eq 1 ] && [ "$usage_percent" -le "$RESOURCE_LOAD_RESUME_PERCENT" ]; then
      resume_child
      load_guard_log "CPU 使用率已回落到 ${usage_percent}%，恢复任务"
    fi

    sleep "$RESOURCE_LOAD_POLL_SECONDS"
  done
}

launch_local() {
  local -a runner=()

  if [ "$use_ionice" -eq 1 ] && command -v ionice >/dev/null 2>&1; then
    runner+=(ionice -c3)
  fi

  if [ -n "$nice_value" ] && command -v nice >/dev/null 2>&1; then
    runner+=(nice -n "$nice_value")
  fi

  runner+=(env)
  if [ "${#child_env[@]}" -gt 0 ]; then
    runner+=("${child_env[@]}")
  fi
  runner+=("${guarded_command[@]}")

  if command -v setsid >/dev/null 2>&1; then
    setsid "${runner[@]}" &
  else
    "${runner[@]}" &
  fi
  runner_pid=$!
  child_pgid="$runner_pid"
}

launch_with_systemd() {
  local -a runner=()

  child_unit="$(build_unit_name "$systemd_unit_prefix")"
  runner+=(
    systemd-run
    --pipe
    --quiet
    --wait
    --collect
    --same-dir
    --unit "$child_unit"
  )

  if [ -n "$nice_value" ]; then
    runner+=("--nice=${nice_value}")
  fi

  if [ "${#systemd_props[@]}" -gt 0 ]; then
    local prop
    for prop in "${systemd_props[@]}"; do
      runner+=("--property=${prop}")
    done
  fi

  if [ "${#systemd_env[@]}" -gt 0 ]; then
    local assignment
    for assignment in "${systemd_env[@]}"; do
      runner+=("--setenv=${assignment}")
    done
  fi

  runner+=("${guarded_command[@]}")

  if [ -t 0 ]; then
    "${runner[@]}" &
  else
    "${runner[@]}" 2> >(grep -Fvx "stty: 'standard input': Inappropriate ioctl for device" >&2) &
  fi
  runner_pid=$!
}

trap 'handle_signal HUP 129' HUP
trap 'handle_signal INT 130' INT
trap 'handle_signal TERM 143' TERM
trap 'stop_monitor' EXIT

load_guard_wait_for_idle

if [ "$should_use_systemd" -eq 1 ]; then
  launch_with_systemd
else
  launch_local
fi

monitor_load &
monitor_pid=$!

set +e
wait "$runner_pid"
command_status=$?
set -e

stop_monitor
exit "$command_status"
