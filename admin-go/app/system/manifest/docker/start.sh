#!/bin/sh
set -eu

cd /app

should_migrate="${DB_MIGRATE_AUTO:-auto}"
if [ "$should_migrate" = "auto" ]; then
  should_migrate="1"
fi

if [ "$should_migrate" != "0" ]; then
  ./migrate up
fi

exec ./main
