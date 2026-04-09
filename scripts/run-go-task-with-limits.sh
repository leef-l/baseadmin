#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/run-go-task-with-limits.sh <command> [args...]

Purpose:
  Run Go build/test tasks with conservative CPU and memory limits so a small
  server is less likely to be overwhelmed.

Environment overrides:
  RESOURCE_MEMORY_MAX=1200M
  RESOURCE_CPU_QUOTA=60%
  RESOURCE_TASKS_MAX=256
  RESOURCE_NICE=10
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

: "${RESOURCE_MEMORY_MAX:=1200M}"
: "${RESOURCE_CPU_QUOTA:=60%}"
: "${RESOURCE_TASKS_MAX:=256}"
: "${RESOURCE_NICE:=10}"
: "${GO_TASK_GOMAXPROCS:=1}"
: "${GO_TASK_GOROOT:=$(go env GOROOT)}"
: "${GO_TASK_GOPATH:=$(go env GOPATH)}"
: "${GO_TASK_GOMODCACHE:=$(go env GOMODCACHE)}"
: "${GO_TASK_GOCACHE:=$(go env GOCACHE)}"

if command -v systemd-run >/dev/null 2>&1 && [ -d /run/systemd/system ] && [ "$(id -u)" -eq 0 ]; then
  exec systemd-run \
    --pipe \
    --quiet \
    --wait \
    --collect \
    --same-dir \
    --nice="${RESOURCE_NICE}" \
    --property="MemoryMax=${RESOURCE_MEMORY_MAX}" \
    --property="CPUQuota=${RESOURCE_CPU_QUOTA}" \
    --property="TasksMax=${RESOURCE_TASKS_MAX}" \
    --setenv="PATH=${PATH}" \
    --setenv="GOMAXPROCS=${GO_TASK_GOMAXPROCS}" \
    --setenv="GOROOT=${GO_TASK_GOROOT}" \
    --setenv="GOPATH=${GO_TASK_GOPATH}" \
    --setenv="GOMODCACHE=${GO_TASK_GOMODCACHE}" \
    --setenv="GOCACHE=${GO_TASK_GOCACHE}" \
    "$@"
fi

if command -v ionice >/dev/null 2>&1; then
  set -- ionice -c3 "$@"
fi

if command -v nice >/dev/null 2>&1; then
  set -- nice -n "${RESOURCE_NICE}" "$@"
fi

export GOMAXPROCS="${GO_TASK_GOMAXPROCS}"
export GOROOT="${GO_TASK_GOROOT}"
export GOPATH="${GO_TASK_GOPATH}"
export GOMODCACHE="${GO_TASK_GOMODCACHE}"
export GOCACHE="${GO_TASK_GOCACHE}"

exec "$@"
