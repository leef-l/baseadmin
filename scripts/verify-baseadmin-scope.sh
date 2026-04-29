#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCOPE_FILE="$ROOT_DIR/contracts/baseadmin-scope.json"
MIGRATION_FILE="$ROOT_DIR/admin-go/database/migrations/000001_baseline_system_upload.up.sql"
CODEGEN_CONFIG="$ROOT_DIR/admin-go/codegen/codegen.yaml"
VIEWS_DIR="$ROOT_DIR/vue-vben-admin/apps/web-antd/src/views"

fail() {
  echo "[verify-baseadmin-scope] $1" >&2
  exit 1
}

if command -v rg >/dev/null 2>&1; then
  SEARCH_BIN="rg"
else
  SEARCH_BIN="grep"
fi

search_fixed() {
  local needle="$1"
  local file="$2"
  if [ "$SEARCH_BIN" = "rg" ]; then
    rg -n -F -- "$needle" "$file" >/dev/null
  else
    grep -n -F -- "$needle" "$file" >/dev/null
  fi
}

search_regex() {
  local pattern="$1"
  local file="$2"
  if [ "$SEARCH_BIN" = "rg" ]; then
    rg -n -- "$pattern" "$file" >/dev/null
  else
    grep -n -E -- "$pattern" "$file" >/dev/null
  fi
}

[ -f "$SCOPE_FILE" ] || fail "missing scope contract: $SCOPE_FILE"
[ -f "$MIGRATION_FILE" ] || fail "missing baseline migration: $MIGRATION_FILE"
[ -f "$CODEGEN_CONFIG" ] || fail "missing codegen config: $CODEGEN_CONFIG"

if ! command -v node >/dev/null 2>&1; then
  fail "missing node, cannot read scope contract"
fi

mapfile -t allowed_apps < <(node -e "const fs=require('fs'); const data=JSON.parse(fs.readFileSync(process.argv[1],'utf8')); for (const item of data.allowedApps || []) console.log(item);" "$SCOPE_FILE")
mapfile -t top_level_paths < <(node -e "const fs=require('fs'); const data=JSON.parse(fs.readFileSync(process.argv[1],'utf8')); for (const item of data.topLevelMenuPaths || []) console.log(item);" "$SCOPE_FILE")
mapfile -t forbidden_paths < <(node -e "const fs=require('fs'); const data=JSON.parse(fs.readFileSync(process.argv[1],'utf8')); for (const item of data.forbiddenMenuPaths || []) console.log(item);" "$SCOPE_FILE")

[ "${#allowed_apps[@]}" -gt 0 ] || fail "scope contract has no allowed apps"
[ "${#top_level_paths[@]}" -gt 0 ] || fail "scope contract has no top-level menu paths"

for path in "${forbidden_paths[@]}"; do
  if search_fixed "$path" "$MIGRATION_FILE"; then
    fail "baseline migration contains forbidden menu path: $path"
  fi
done

mapfile -t actual_top_level_paths < <(python3 - "$MIGRATION_FILE" <<'PY'
import re
import sys

path = sys.argv[1]
with open(path, 'r', encoding='utf-8') as fh:
    data = fh.read()

pattern = re.compile(r"\(\s*\d+\s*,\s*0\s*,\s*'[^']+'\s*,\s*1\s*,\s*'([^']+)'\s*,", re.MULTILINE)
for item in sorted(set(pattern.findall(data))):
    print(item)
PY
)

for path in "${actual_top_level_paths[@]}"; do
  found=0
  for allowed in "${top_level_paths[@]}"; do
    if [ "$path" = "$allowed" ]; then
      found=1
      break
    fi
  done
  [ "$found" -eq 1 ] || fail "baseline migration contains unexpected top-level menu path: $path"
done

for app in "${allowed_apps[@]}"; do
  search_regex "^[[:space:]]+$app:" "$CODEGEN_CONFIG" || fail "codegen allowed app missing from menu_apps: $app"
done

mapfile -t config_allowed_apps < <(awk '
  /^allowed_apps:/ { section="allowed_apps"; next }
  /^[^[:space:]]/ { section=""; next }
  section == "allowed_apps" && /^[[:space:]]*-[[:space:]]+/ {
    line = $0
    sub(/^[[:space:]]*-[[:space:]]+/, "", line)
    print line
  }
' "$CODEGEN_CONFIG")

mapfile -t menu_app_keys < <(awk '
  /^menu_apps:/ { section="menu_apps"; next }
  /^[^[:space:]]/ { section=""; next }
  section == "menu_apps" && /^[[:space:]]{2}[A-Za-z0-9_]+:/ {
    line = $0
    sub(/^[[:space:]]+/, "", line)
    sub(/:.*/, "", line)
    print line
  }
' "$CODEGEN_CONFIG")

for app in "${allowed_apps[@]}"; do
  found=0
  for current in "${config_allowed_apps[@]}"; do
    if [ "$current" = "$app" ]; then
      found=1
      break
    fi
  done
  [ "$found" -eq 1 ] || fail "allowed_apps missing required app: $app"
done

for app in "${menu_app_keys[@]}"; do
  found=0
  for allowed in "${allowed_apps[@]}"; do
    if [ "$app" = "$allowed" ]; then
      found=1
      break
    fi
  done
  [ "$found" -eq 1 ] || fail "menu_apps contains unexpected app: $app"
done

for dir in "$VIEWS_DIR"/*; do
  [ -d "$dir" ] || continue
  base="$(basename "$dir")"
  if [ "$base" = "_core" ]; then
    continue
  fi
  found=0
  for allowed in "${allowed_apps[@]}"; do
    if [ "$base" = "$allowed" ]; then
      found=1
      break
    fi
  done
  [ "$found" -eq 1 ] || fail "unexpected top-level view directory: $base"
done

echo "[verify-baseadmin-scope] ok"
