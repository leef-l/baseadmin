#!/usr/bin/env bash

# Guard interactive/agent shells started inside this repository from running
# high-risk local commands directly. Heavy Node commands must go through
# scripts/run-node-task-with-limits.sh, which sets BASEADMIN_ALLOW_NODE_TASK=1.

baseadmin_guard_message() {
  local command_name="$1"

  case "$command_name" in
    docker|docker-compose)
      printf '%s\n' \
        "Blocked by baseadmin shell guard: local Docker execution is forbidden by CLAUDE.md." \
        "Use documented remote/container workflow instead of running ${command_name} here." >&2
      ;;
    *)
      printf '%s\n' \
        "Blocked by baseadmin shell guard: local ${command_name} is forbidden unless it is wrapped." \
        "Use: ./scripts/run-node-task-with-limits.sh ${command_name} ..." \
        "Hard wrapper limits: timeout 1h, MemoryMax 1536M, CPU guard 80%/50%." >&2
      ;;
  esac
}

baseadmin_guard_node_command() {
  local command_name="$1"
  shift

  if [ "${BASEADMIN_ALLOW_NODE_TASK:-}" = "1" ]; then
    command "$command_name" "$@"
    return $?
  fi

  baseadmin_guard_message "$command_name"
  return 126
}

baseadmin_guard_docker_command() {
  local command_name="$1"
  shift

  if [ "${BASEADMIN_ALLOW_DOCKER_TASK:-}" = "1" ]; then
    command "$command_name" "$@"
    return $?
  fi

  baseadmin_guard_message "$command_name"
  return 126
}

pnpm() { baseadmin_guard_node_command pnpm "$@"; }
npm() { baseadmin_guard_node_command npm "$@"; }
npx() { baseadmin_guard_node_command npx "$@"; }
pnpx() { baseadmin_guard_node_command pnpx "$@"; }
yarn() { baseadmin_guard_node_command yarn "$@"; }
corepack() { baseadmin_guard_node_command corepack "$@"; }
docker() { baseadmin_guard_docker_command docker "$@"; }
docker-compose() { baseadmin_guard_docker_command docker-compose "$@"; }

export -f baseadmin_guard_message
export -f baseadmin_guard_node_command
export -f baseadmin_guard_docker_command
export -f pnpm npm npx pnpx yarn corepack docker docker-compose
