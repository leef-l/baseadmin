#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VIEWS_DIR="$ROOT_DIR/vue-vben-admin/apps/web-antd/src/views"

fail() {
  echo "[verify-vben-pages] $1" >&2
  exit 1
}

if ! command -v rg >/dev/null 2>&1; then
  fail "缺少 rg，无法执行页面规则校验"
fi

# 业务页必须走当前仓库的 adapter/form + vben modal/grid 方案。
business_pages=(
  "$VIEWS_DIR/system"
  "$VIEWS_DIR/upload"
)

for dir in "${business_pages[@]}"; do
  if rg -n "rules:\s*\[" "$dir" -g '*.vue' >/dev/null; then
    fail "业务页存在数组 rules 写法，请统一使用当前 vben form 规则约定: $dir"
  fi
  if rg -n "document\.createElement|createElement\('input'\)|createElement\('a'\)" "$dir" -g '*.vue' >/dev/null; then
    fail "业务页存在裸 DOM 交互写法，请改用当前 vben 组件/工具链: $dir"
  fi
done

if rg --files-without-match "useVbenVxeGrid" "$VIEWS_DIR/system" "$VIEWS_DIR/upload" -g 'index.vue' >/dev/null; then
  fail "存在未使用 useVbenVxeGrid 的业务列表页"
fi

if rg --files-without-match "useVbenForm" "$VIEWS_DIR/system" "$VIEWS_DIR/upload" -g 'form.vue' >/dev/null; then
  fail "存在未使用 useVbenForm 的业务表单页"
fi

if rg --files-without-match "useVbenModal" "$VIEWS_DIR/system" "$VIEWS_DIR/upload" -g 'form.vue' >/dev/null; then
  fail "存在未使用 useVbenModal 的业务弹窗表单页"
fi

# 认证页使用官方 Authentication* 组件和 Vben* 基础表单组件，属于仓库允许的第二套 vben 风格。
auth_dir="$VIEWS_DIR/_core/authentication"
if rg --files-without-match "Authentication(Login|Register|CodeLogin|ForgetPassword|QrCodeLogin)" "$auth_dir" -g '*.vue' >/dev/null; then
  fail "认证页必须使用官方 Authentication* 组件"
fi

echo "[verify-vben-pages] ok"
