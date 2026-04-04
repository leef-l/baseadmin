#!/usr/bin/env sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
PROJECT_ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/../.." && pwd)
ENV_SOURCE="$SCRIPT_DIR/.env"
ENV_TARGET="$PROJECT_ROOT/admin-go/.env"
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.yml"

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

exec docker compose --env-file "$ENV_SOURCE" -f "$COMPOSE_FILE" "$@"
