#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/run-node-task-with-limits.sh <command> [args...]

Purpose:
  Run npm/pnpm build tasks with conservative CPU and memory limits so a small
  server is less likely to be overwhelmed.

Environment overrides:
  RESOURCE_MEMORY_MAX=1200M
  RESOURCE_CPU_QUOTA=60%
  RESOURCE_TASKS_MAX=256
  RESOURCE_NICE=10
  PNPM_NETWORK_CONCURRENCY=3
  PNPM_CHILD_CONCURRENCY=1
  NPM_CONFIG_MAXSOCKETS=6
  FRONTEND_NODE_MAX_OLD_SPACE_SIZE=1024
EOF
}

if [ "$#" -eq 0 ]; then
  usage >&2
  exit 1
fi

: "${RESOURCE_MEMORY_MAX:=1200M}"
: "${RESOURCE_CPU_QUOTA:=60%}"
: "${RESOURCE_TASKS_MAX:=256}"
: "${RESOURCE_NICE:=10}"
: "${PNPM_NETWORK_CONCURRENCY:=3}"
: "${PNPM_CHILD_CONCURRENCY:=1}"
: "${NPM_CONFIG_MAXSOCKETS:=6}"
: "${FRONTEND_NODE_MAX_OLD_SPACE_SIZE:=1024}"

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
    --setenv="PNPM_NETWORK_CONCURRENCY=${PNPM_NETWORK_CONCURRENCY}" \
    --setenv="PNPM_CHILD_CONCURRENCY=${PNPM_CHILD_CONCURRENCY}" \
    --setenv="npm_config_maxsockets=${NPM_CONFIG_MAXSOCKETS}" \
    --setenv="npm_config_node_options=--max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
    --setenv="FRONTEND_NODE_MAX_OLD_SPACE_SIZE=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
    "$@"
fi

if command -v ionice >/dev/null 2>&1; then
  set -- ionice -c3 "$@"
fi

if command -v nice >/dev/null 2>&1; then
  set -- nice -n "${RESOURCE_NICE}" "$@"
fi

export PNPM_NETWORK_CONCURRENCY
export PNPM_CHILD_CONCURRENCY
export FRONTEND_NODE_MAX_OLD_SPACE_SIZE
export npm_config_maxsockets="${NPM_CONFIG_MAXSOCKETS}"
export npm_config_node_options="--max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}"

exec "$@"
