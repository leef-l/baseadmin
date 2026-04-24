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

Hard limits:
  timeout 固定 1h，超时后先 TERM，30s 后 KILL。
  MemoryMax 固定 1536M，防止 Node 任务拖垮宿主机。
  本地 fallback 固定 ulimit -v 1572864，约 1.5G。
  CPU guard 固定 80% 暂停、50% 恢复。
  FRONTEND_NODE_MAX_OLD_SPACE_SIZE 固定 1024M，低于硬上限。

Fixed execution policy:
  RESOURCE_TASK_TIMEOUT=1h
  RESOURCE_MEMORY_HIGH=1024M
  RESOURCE_MEMORY_MAX=1536M
  RESOURCE_MEMORY_SWAP_MAX=0
  RESOURCE_OOM_SCORE_ADJUST=500
  RESOURCE_CPU_QUOTA=60%
  RESOURCE_TASKS_MAX=256
  RESOURCE_NICE=10
  RESOURCE_LOAD_MAX_PERCENT=80
  RESOURCE_LOAD_RESUME_PERCENT=50
  RESOURCE_ULIMIT_V_KB=1572864
  PNPM_NETWORK_CONCURRENCY=2
  PNPM_CHILD_CONCURRENCY=1
  NPM_CONFIG_MAXSOCKETS=4
  FRONTEND_NODE_MAX_OLD_SPACE_SIZE=1024
EOF
}

if [ "$#" -eq 0 ]; then
  usage >&2
  exit 1
fi

if [ "${1:-}" = "-h" ] || [ "${1:-}" = "--help" ]; then
  usage
  exit 0
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

RESOURCE_TASK_TIMEOUT=1h
RESOURCE_MEMORY_HIGH=1024M
RESOURCE_MEMORY_MAX=1536M
RESOURCE_MEMORY_SWAP_MAX=0
RESOURCE_OOM_SCORE_ADJUST=500
RESOURCE_CPU_QUOTA=60%
RESOURCE_TASKS_MAX=256
RESOURCE_NICE=10
RESOURCE_LOAD_MAX_PERCENT=80
RESOURCE_LOAD_RESUME_PERCENT=50
RESOURCE_ULIMIT_V_KB=1572864
LOAD_GUARD_LOG_PREFIX=node-guard
PNPM_NETWORK_CONCURRENCY=2
PNPM_CHILD_CONCURRENCY=1
NPM_CONFIG_MAXSOCKETS=4
FRONTEND_NODE_MAX_OLD_SPACE_SIZE=1024

# 传给通用负载守卫，避免日志前缀在 exec 后退回默认值。
export RESOURCE_LOAD_MAX_PERCENT RESOURCE_LOAD_RESUME_PERCENT LOAD_GUARD_LOG_PREFIX
export BASEADMIN_ALLOW_NODE_TASK=1

if ! command -v timeout >/dev/null 2>&1; then
  echo "timeout command not found" >&2
  exit 1
fi

exec "$SCRIPT_DIR/run-task-with-load-guard.sh" \
  --nice "${RESOURCE_NICE}" \
  --systemd-unit-prefix node-task \
  --systemd-prop "MemoryHigh=${RESOURCE_MEMORY_HIGH}" \
  --systemd-prop "MemoryMax=${RESOURCE_MEMORY_MAX}" \
  --systemd-prop "MemorySwapMax=${RESOURCE_MEMORY_SWAP_MAX}" \
  --systemd-prop "OOMScoreAdjust=${RESOURCE_OOM_SCORE_ADJUST}" \
  --systemd-prop "CPUQuota=${RESOURCE_CPU_QUOTA}" \
  --systemd-prop "TasksMax=${RESOURCE_TASKS_MAX}" \
  --systemd-env "PATH=${PATH}" \
  --systemd-env "PNPM_NETWORK_CONCURRENCY=${PNPM_NETWORK_CONCURRENCY}" \
  --systemd-env "PNPM_CHILD_CONCURRENCY=${PNPM_CHILD_CONCURRENCY}" \
  --systemd-env "npm_config_maxsockets=${NPM_CONFIG_MAXSOCKETS}" \
  --systemd-env "npm_config_node_options=--max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --systemd-env "NODE_OPTIONS=--max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --systemd-env "FRONTEND_NODE_MAX_OLD_SPACE_SIZE=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --systemd-env "BASEADMIN_ALLOW_NODE_TASK=1" \
  --env "PNPM_NETWORK_CONCURRENCY=${PNPM_NETWORK_CONCURRENCY}" \
  --env "PNPM_CHILD_CONCURRENCY=${PNPM_CHILD_CONCURRENCY}" \
  --env "npm_config_maxsockets=${NPM_CONFIG_MAXSOCKETS}" \
  --env "npm_config_node_options=--max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --env "NODE_OPTIONS=--max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --env "FRONTEND_NODE_MAX_OLD_SPACE_SIZE=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --env "BASEADMIN_ALLOW_NODE_TASK=1" \
  -- bash -c 'ulimit -v "$1" 2>/dev/null || true; shift; exec timeout -k 30s 1h "$@"' bash "$RESOURCE_ULIMIT_V_KB" "$@"
