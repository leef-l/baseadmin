#!/bin/sh
set -eu

cd /app

app_name="${1:-}"

if [ -z "$app_name" ]; then
  echo "usage: start-go-app.sh <system|upload>"
  exit 1
fi

should_migrate="${DB_MIGRATE_AUTO:-auto}"
if [ "$should_migrate" = "auto" ]; then
  if [ "$app_name" = "system" ]; then
    should_migrate="1"
  else
    should_migrate="0"
  fi
fi

if [ "$should_migrate" != "0" ]; then
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
