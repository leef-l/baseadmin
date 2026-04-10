#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/run-node-task-with-limits.sh <command> [args...]

Purpose:
  Run npm/pnpm build tasks with conservative CPU and memory limits so a small
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

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

: "${RESOURCE_MEMORY_MAX:=1200M}"
: "${RESOURCE_CPU_QUOTA:=60%}"
: "${RESOURCE_TASKS_MAX:=256}"
: "${RESOURCE_NICE:=10}"
: "${LOAD_GUARD_LOG_PREFIX:=node-guard}"
: "${PNPM_NETWORK_CONCURRENCY:=3}"
: "${PNPM_CHILD_CONCURRENCY:=1}"
: "${NPM_CONFIG_MAXSOCKETS:=6}"
: "${FRONTEND_NODE_MAX_OLD_SPACE_SIZE:=1024}"

# 传给通用负载守卫，避免日志前缀在 exec 后退回默认值。
export RESOURCE_LOAD_MAX_PERCENT RESOURCE_LOAD_RESUME_PERCENT LOAD_GUARD_LOG_PREFIX

exec "$SCRIPT_DIR/run-task-with-load-guard.sh" \
  --nice "${RESOURCE_NICE}" \
  --systemd-unit-prefix node-task \
  --systemd-prop "MemoryMax=${RESOURCE_MEMORY_MAX}" \
  --systemd-prop "CPUQuota=${RESOURCE_CPU_QUOTA}" \
  --systemd-prop "TasksMax=${RESOURCE_TASKS_MAX}" \
  --systemd-env "PATH=${PATH}" \
  --systemd-env "PNPM_NETWORK_CONCURRENCY=${PNPM_NETWORK_CONCURRENCY}" \
  --systemd-env "PNPM_CHILD_CONCURRENCY=${PNPM_CHILD_CONCURRENCY}" \
  --systemd-env "npm_config_maxsockets=${NPM_CONFIG_MAXSOCKETS}" \
  --systemd-env "npm_config_node_options=--max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --systemd-env "FRONTEND_NODE_MAX_OLD_SPACE_SIZE=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --env "PNPM_NETWORK_CONCURRENCY=${PNPM_NETWORK_CONCURRENCY}" \
  --env "PNPM_CHILD_CONCURRENCY=${PNPM_CHILD_CONCURRENCY}" \
  --env "npm_config_maxsockets=${NPM_CONFIG_MAXSOCKETS}" \
  --env "npm_config_node_options=--max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --env "FRONTEND_NODE_MAX_OLD_SPACE_SIZE=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  -- "$@"
