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

if [ "${NPM_REGISTRY:-}" != "" ]; then
  npm config set registry "${NPM_REGISTRY}"
  pnpm config set registry "${NPM_REGISTRY}"
fi

if [ ! -d node_modules/.pnpm ]; then
  echo "Bootstrapping frontend dependencies..."
  pnpm install
fi

exec pnpm dev:antd --host 0.0.0.0
