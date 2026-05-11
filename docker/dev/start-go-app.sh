#!/bin/sh
set -eu

cd /app

app_name="${1:-}"

if [ -z "$app_name" ]; then
  echo "usage: start-go-app.sh <system|upload>"
  exit 1
fi

should_migrate="${DB_MIGRATE_AUTO:-auto}"
normalized_migrate="$(printf '%s' "$should_migrate" | tr '[:upper:]' '[:lower:]')"
if [ "$normalized_migrate" = "auto" ]; then
  if [ "$app_name" = "system" ]; then
    normalized_migrate="1"
  else
    normalized_migrate="0"
  fi
fi

case "$normalized_migrate" in
  0|false|no|off)
    normalized_migrate="0"
    ;;
  *)
    normalized_migrate="1"
    ;;
esac

if [ "$normalized_migrate" != "0" ]; then
  ./migrate up
fi

case "$app_name" in
  system)
    exec ./system
    ;;
  upload)
    exec ./upload
    ;;
  *)
    echo "unsupported app: $app_name"
    exit 1
    ;;
esac
