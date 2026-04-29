#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VIEWS_DIR="$ROOT_DIR/vue-vben-admin/apps/web-antd/src/views"

fail() {
  echo "[verify-vben-pages] $1" >&2
  exit 1
}

if command -v rg >/dev/null 2>&1; then
  SEARCH_BIN="rg"
else
  SEARCH_BIN="grep"
fi

find_named_files() {
  local name="$1"
  shift
  find "$@" -type f -name "$name" -print 2>/dev/null
}

first_match_in_files() {
  local pattern="$1"
  local name="$2"
  shift 2
  local file

  while IFS= read -r file; do
    [ -n "$file" ] || continue
    if [ "$SEARCH_BIN" = "rg" ]; then
      if rg -n -- "$pattern" "$file" >/dev/null; then
        printf '%s\n' "$file"
        return 0
      fi
    else
      if grep -n -E -- "$pattern" "$file" >/dev/null; then
        printf '%s\n' "$file"
        return 0
      fi
    fi
  done < <(find_named_files "$name" "$@")

  return 1
}

first_file_without_match() {
  local pattern="$1"
  local name="$2"
  shift 2
  local file

  while IFS= read -r file; do
    [ -n "$file" ] || continue
    if [ "$SEARCH_BIN" = "rg" ]; then
      if ! rg -n -- "$pattern" "$file" >/dev/null; then
        printf '%s\n' "$file"
        return 0
      fi
    else
      if ! grep -n -E -- "$pattern" "$file" >/dev/null; then
        printf '%s\n' "$file"
        return 0
      fi
    fi
  done < <(find_named_files "$name" "$@")

  return 1
}

# 业务页必须走当前仓库的 adapter/form + vben modal/grid 方案。
business_pages=(
  "$VIEWS_DIR/system"
  "$VIEWS_DIR/upload"
  "$VIEWS_DIR/demo"
)

for dir in "${business_pages[@]}"; do
  if first_match_in_files "rules:[[:space:]]*\\[" "*.vue" "$dir" >/dev/null; then
    fail "业务页存在数组 rules 写法，请统一使用当前 vben form 规则约定: $dir"
  fi
  if first_match_in_files "document\\.createElement|createElement\\('input'\\)|createElement\\('a'\\)" "*.vue" "$dir" >/dev/null; then
    fail "业务页存在裸 DOM 交互写法，请改用当前 vben 组件/工具链: $dir"
  fi
done

if first_file_without_match "useVbenVxeGrid" "index.vue" "${business_pages[@]}" >/dev/null; then
  fail "存在未使用 useVbenVxeGrid 的业务列表页"
fi

if first_file_without_match "useVbenForm" "form.vue" "${business_pages[@]}" >/dev/null; then
  fail "存在未使用 useVbenForm 的业务表单页"
fi

if first_file_without_match "useVbenModal" "form.vue" "${business_pages[@]}" >/dev/null; then
  fail "存在未使用 useVbenModal 的业务弹窗表单页"
fi

# 认证页使用官方 Authentication* 组件和 Vben* 基础表单组件，属于仓库允许的第二套 vben 风格。
auth_dir="$VIEWS_DIR/_core/authentication"
if first_file_without_match "Authentication(Login|Register|CodeLogin|ForgetPassword|QrCodeLogin)" "*.vue" "$auth_dir" >/dev/null; then
  fail "认证页必须使用官方 Authentication* 组件"
fi

echo "[verify-vben-pages] ok"
