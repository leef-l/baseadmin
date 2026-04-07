#!/bin/bash

# This shell is executed before docker build.
set -eu

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
APP_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
ROOT_DIR="$(cd "$APP_DIR/../.." && pwd)"
TARGET_DIR="$APP_DIR/temp/linux_amd64"

mkdir -p "$TARGET_DIR/database"

cd "$ROOT_DIR"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$TARGET_DIR/migrate" ./cmd/migrate
rm -rf "$TARGET_DIR/database/migrations"
cp -R "$ROOT_DIR/database/migrations" "$TARGET_DIR/database/migrations"
