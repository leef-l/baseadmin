#!/bin/sh
set -eu

cd /app

if ! command -v git >/dev/null 2>&1; then
  apk add --no-cache git
fi

if ! command -v pnpm >/dev/null 2>&1; then
  corepack enable
  corepack prepare pnpm@10.32.1 --activate
fi

: "${PNPM_NETWORK_CONCURRENCY:=1}"
: "${PNPM_CHILD_CONCURRENCY:=1}"
: "${FRONTEND_INSTALL_NODE_MAX_OLD_SPACE_SIZE:=256}"
: "${FRONTEND_NODE_MAX_OLD_SPACE_SIZE:=256}"

if [ "${NPM_REGISTRY:-}" != "" ]; then
  npm config set registry "${NPM_REGISTRY}"
  pnpm config set registry "${NPM_REGISTRY}"
fi

export NODE_OPTIONS="${NODE_OPTIONS:---max-old-space-size=${FRONTEND_NODE_MAX_OLD_SPACE_SIZE}}"

if [ ! -d node_modules/.pnpm ]; then
  echo "Bootstrapping frontend dependencies with constrained resources..."
  env \
    CI=1 \
    NODE_OPTIONS="--max-old-space-size=${FRONTEND_INSTALL_NODE_MAX_OLD_SPACE_SIZE}" \
    pnpm install \
      --frozen-lockfile \
      --network-concurrency="${PNPM_NETWORK_CONCURRENCY}" \
      --child-concurrency="${PNPM_CHILD_CONCURRENCY}" \
      --reporter=append-only
fi

exec pnpm dev:antd --host 0.0.0.0
