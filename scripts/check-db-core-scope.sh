#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="$ROOT_DIR/admin-go/.env"

fail() {
  echo "[check-db-core-scope] $1" >&2
  exit 1
}

[ -f "$ENV_FILE" ] || fail "missing env file: $ENV_FILE"
command -v mysql >/dev/null 2>&1 || fail "missing mysql client"

set -a
. "$ENV_FILE"
set +a

: "${MYSQL_HOST:?}"
: "${MYSQL_PORT:?}"
: "${MYSQL_USER:?}"
: "${MYSQL_PASSWORD:?}"
: "${MYSQL_DATABASE:?}"

mysql_exec() {
  MYSQL_PWD="$MYSQL_PASSWORD" mysql \
    -h"$MYSQL_HOST" \
    -P"$MYSQL_PORT" \
    -u"$MYSQL_USER" \
    -D"$MYSQL_DATABASE" \
    -N -B \
    -e "$1"
}

table_count="$(mysql_exec "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='${MYSQL_DATABASE}';")"
table_count="${table_count:-0}"

echo "[check-db-core-scope] database=${MYSQL_DATABASE} tables=${table_count}"

if [ "$table_count" -eq 0 ]; then
  echo "[check-db-core-scope] database is empty"
  exit 0
fi

unexpected_tables="$(mysql_exec "SELECT table_name FROM information_schema.tables WHERE table_schema='${MYSQL_DATABASE}' AND table_name NOT LIKE 'system\\_%' AND table_name NOT LIKE 'upload\\_%' AND table_name <> 'schema_migrations' ORDER BY table_name;")"
if [ -n "$unexpected_tables" ]; then
  echo "$unexpected_tables"
  fail "found non-core tables"
fi

if mysql_exec "SELECT 1 FROM information_schema.tables WHERE table_schema='${MYSQL_DATABASE}' AND table_name='system_menu' LIMIT 1;" | grep -q 1; then
  forbidden_menus="$(mysql_exec "SELECT CONCAT(id, '\t', IFNULL(path,''), '\t', title) FROM system_menu WHERE deleted_at IS NULL AND path IN ('/dashboard','/analytics','/workspace') ORDER BY id;")"
  if [ -n "$forbidden_menus" ]; then
    echo "$forbidden_menus"
    fail "found forbidden dashboard menus"
  fi

  top_level_paths="$(mysql_exec "SELECT path FROM system_menu WHERE deleted_at IS NULL AND parent_id = 0 AND type = 1 ORDER BY sort, id;")"
  if [ -n "$top_level_paths" ]; then
    while IFS= read -r path; do
      case "$path" in
        /system|/upload)
          ;;
        *)
          fail "unexpected top-level menu path: $path"
          ;;
      esac
    done <<< "$top_level_paths"
  fi
fi

echo "[check-db-core-scope] ok"
