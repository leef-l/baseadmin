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

Memory policy (cgroup v2 soft-limit + swap spillover):
  RESOURCE_MEMORY_HIGH 是软限制：超过后内核在 memcg 回收路径上会主动扫描
  匿名 LRU 并把冷页换出到 swap，同时节流进程，但不会触发 OOM kill。
  RESOURCE_MEMORY_MAX 是硬上限，仅作为"失控兜底"，一般设为 HIGH 的 1.5–2×，
  中间留的这段空间就是软着陆区。
  RESOURCE_MEMORY_SWAP_MAX 让该 cgroup 能吃到多少 swap。
  FRONTEND_NODE_MAX_OLD_SPACE_SIZE 必须 ≥ RESOURCE_MEMORY_HIGH，否则 V8 会
  先于 cgroup 软限制 abort，swap 溢出机制就失效了。
  注：systemd 不暴露 MemorySwappiness 单元属性，如需更激进的换出倾向，请
  在宿主上调整 /proc/sys/vm/swappiness（全局），而不是在此脚本里设置。

Environment overrides:
  RESOURCE_MEMORY_HIGH=1200M
  RESOURCE_MEMORY_MAX=2400M
  RESOURCE_MEMORY_SWAP_MAX=1800M
  RESOURCE_OOM_SCORE_ADJUST=500
  RESOURCE_CPU_QUOTA=60%
  RESOURCE_TASKS_MAX=256
  RESOURCE_NICE=10
  RESOURCE_LOAD_MAX_PERCENT=80
  RESOURCE_LOAD_RESUME_PERCENT=50
  PNPM_NETWORK_CONCURRENCY=3
  PNPM_CHILD_CONCURRENCY=1
  NPM_CONFIG_MAXSOCKETS=6
  FRONTEND_NODE_MAX_OLD_SPACE_SIZE=2048
EOF
}

if [ "$#" -eq 0 ]; then
  usage >&2
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

: "${RESOURCE_MEMORY_HIGH:=1200M}"
: "${RESOURCE_MEMORY_MAX:=2400M}"
: "${RESOURCE_MEMORY_SWAP_MAX:=1800M}"
: "${RESOURCE_OOM_SCORE_ADJUST:=500}"
: "${RESOURCE_CPU_QUOTA:=60%}"
: "${RESOURCE_TASKS_MAX:=256}"
: "${RESOURCE_NICE:=10}"
: "${LOAD_GUARD_LOG_PREFIX:=node-guard}"
: "${PNPM_NETWORK_CONCURRENCY:=3}"
: "${PNPM_CHILD_CONCURRENCY:=1}"
: "${NPM_CONFIG_MAXSOCKETS:=6}"
: "${FRONTEND_NODE_MAX_OLD_SPACE_SIZE:=2048}"

# 传给通用负载守卫，避免日志前缀在 exec 后退回默认值。
export RESOURCE_LOAD_MAX_PERCENT RESOURCE_LOAD_RESUME_PERCENT LOAD_GUARD_LOG_PREFIX

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
  --env "PNPM_NETWORK_CONCURRENCY=${PNPM_NETWORK_CONCURRENCY}" \
  --env "PNPM_CHILD_CONCURRENCY=${PNPM_CHILD_CONCURRENCY}" \
  --env "npm_config_maxsockets=${NPM_CONFIG_MAXSOCKETS}" \
  --env "npm_config_node_options=--max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --env "NODE_OPTIONS=--max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  --env "FRONTEND_NODE_MAX_OLD_SPACE_SIZE=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}" \
  -- "$@"
