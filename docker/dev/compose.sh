#!/usr/bin/env sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
PROJECT_ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/../.." && pwd)
ENV_SOURCE="$SCRIPT_DIR/.env"
ENV_TARGET="$PROJECT_ROOT/admin-go/.env"
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.yml"

contains_frontend_request() {
  prev=''
  for arg in "$@"; do
    case "$arg" in
      --profile=frontend|frontend)
        return 0
        ;;
    esac

    if [ "$prev" = "--profile" ] && [ "$arg" = "frontend" ]; then
      return 0
    fi

    prev="$arg"
  done

  return 1
}

get_available_memory_mb() {
  if [ -r /proc/meminfo ]; then
    awk '/MemAvailable:/ { printf "%d", $2 / 1024; exit }' /proc/meminfo
    return 0
  fi

  if command -v vm_stat >/dev/null 2>&1; then
    vm_stat | awk '
      BEGIN { pages = 0; page_size = 4096 }
      /page size of/ {
        for (i = 1; i <= NF; i++) {
          if ($i == "of") {
            page_size = $(i + 1)
            break
          }
        }
      }
      /^Pages free:/ || /^Pages speculative:/ || /^Pages inactive:/ {
        value = $3
        gsub("\\.", "", value)
        pages += value
      }
      END { printf "%d", pages * page_size / 1024 / 1024 }
    '
    return 0
  fi

  return 1
}

guard_frontend_start() {
  if [ "${ALLOW_LOW_MEMORY_FRONTEND:-0}" = "1" ]; then
    echo "[WARN] Skipping frontend memory guard because ALLOW_LOW_MEMORY_FRONTEND=1" >&2
    return 0
  fi

  min_mb="${FRONTEND_MIN_HOST_MEM_MB:-2048}"
  if ! available_mb=$(get_available_memory_mb); then
    echo "[WARN] Unable to determine host available memory, continuing without frontend guard" >&2
    return 0
  fi

  if [ -z "$available_mb" ]; then
    echo "[WARN] Unable to determine host available memory, continuing without frontend guard" >&2
    return 0
  fi

  if [ "$available_mb" -lt "$min_mb" ]; then
    echo "[ERROR] Refusing to start frontend: host available memory ${available_mb}MB is below FRONTEND_MIN_HOST_MEM_MB=${min_mb}MB" >&2
    echo "[ERROR] Run backend-only compose, or set ALLOW_LOW_MEMORY_FRONTEND=1 if you accept the risk" >&2
    exit 1
  fi
}

if [ "${1:-}" = "-China" ] || [ "${1:-}" = "--china" ]; then
  COMPOSE_FILE="$SCRIPT_DIR/docker-compose.cn.yml"
  shift
fi

if [ ! -f "$ENV_SOURCE" ]; then
  echo "Missing env file: $ENV_SOURCE" >&2
  exit 1
fi

cp "$ENV_SOURCE" "$ENV_TARGET"
echo "[INFO] Synced $ENV_SOURCE -> $ENV_TARGET"

if [ "$#" -eq 0 ]; then
  set -- up -d --build
fi

if contains_frontend_request "$@"; then
  guard_frontend_start
fi

exec docker compose --env-file "$ENV_SOURCE" -f "$COMPOSE_FILE" "$@"
