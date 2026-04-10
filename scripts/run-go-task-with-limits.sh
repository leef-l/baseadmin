#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/run-go-task-with-limits.sh <command> [args...]

Purpose:
  Run Go build/test tasks with conservative CPU and memory limits so a small
  server is less likely to be overwhelmed. If CPU usage reaches the pause
  threshold while the task is running, the task is paused and resumed
  automatically after the host recovers.

Environment overrides:
  RESOURCE_MEMORY_MAX=1200M
  RESOURCE_CPU_QUOTA=60%
  RESOURCE_TASKS_MAX=256
  RESOURCE_NICE=10
  RESOURCE_LOAD_MAX_PERCENT=80
  RESOURCE_LOAD_RESUME_PERCENT=50
  GO_TASK_GOMAXPROCS=1
  GO_TASK_GOROOT=<go env GOROOT>
  GO_TASK_GOPATH=<go env GOPATH>
  GO_TASK_GOMODCACHE=<go env GOMODCACHE>
  GO_TASK_GOCACHE=<go env GOCACHE>
EOF
}

if [ "$#" -eq 0 ]; then
  usage >&2
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "go command not found" >&2
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

: "${RESOURCE_MEMORY_MAX:=1200M}"
: "${RESOURCE_CPU_QUOTA:=60%}"
: "${RESOURCE_TASKS_MAX:=256}"
: "${RESOURCE_NICE:=10}"
: "${LOAD_GUARD_LOG_PREFIX:=go-guard}"
: "${GO_TASK_GOMAXPROCS:=1}"
: "${GO_TASK_GOROOT:=$(go env GOROOT)}"
: "${GO_TASK_GOPATH:=$(go env GOPATH)}"
: "${GO_TASK_GOMODCACHE:=$(go env GOMODCACHE)}"
: "${GO_TASK_GOCACHE:=$(go env GOCACHE)}"

# 传给通用负载守卫，避免日志前缀在 exec 后退回默认值。
export RESOURCE_LOAD_MAX_PERCENT RESOURCE_LOAD_RESUME_PERCENT LOAD_GUARD_LOG_PREFIX

exec "$SCRIPT_DIR/run-task-with-load-guard.sh" \
  --nice "${RESOURCE_NICE}" \
  --systemd-unit-prefix go-task \
  --systemd-prop "MemoryMax=${RESOURCE_MEMORY_MAX}" \
  --systemd-prop "CPUQuota=${RESOURCE_CPU_QUOTA}" \
  --systemd-prop "TasksMax=${RESOURCE_TASKS_MAX}" \
  --systemd-env "PATH=${PATH}" \
  --systemd-env "GOMAXPROCS=${GO_TASK_GOMAXPROCS}" \
  --systemd-env "GOROOT=${GO_TASK_GOROOT}" \
  --systemd-env "GOPATH=${GO_TASK_GOPATH}" \
  --systemd-env "GOMODCACHE=${GO_TASK_GOMODCACHE}" \
  --systemd-env "GOCACHE=${GO_TASK_GOCACHE}" \
  --env "GOMAXPROCS=${GO_TASK_GOMAXPROCS}" \
  --env "GOROOT=${GO_TASK_GOROOT}" \
  --env "GOPATH=${GO_TASK_GOPATH}" \
  --env "GOMODCACHE=${GO_TASK_GOMODCACHE}" \
  --env "GOCACHE=${GO_TASK_GOCACHE}" \
  -- "$@"
