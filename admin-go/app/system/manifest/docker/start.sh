#!/bin/sh
set -eu

cd /app

should_migrate="${DB_MIGRATE_AUTO:-auto}"
normalized_migrate="$(printf '%s' "$should_migrate" | tr '[:upper:]' '[:lower:]')"
if [ "$normalized_migrate" = "auto" ]; then
  normalized_migrate="1"
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

exec ./main
