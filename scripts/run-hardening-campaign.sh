#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/run-hardening-campaign.sh [--start N] [--rounds N] [--publish-every N]

Purpose:
  Execute repeated hardening validation rounds with bounded resource usage.
  Each round prints the round number. Every publish interval runs verifye2e,
  updates the hardening doc checkpoint, and publishes changes if the worktree
  is dirty. When codegen-related files change, the script also runs demo app
  verification in a temporary workspace before the round can be considered done.

Examples:
  ./scripts/run-hardening-campaign.sh --start 2 --rounds 1
  ./scripts/run-hardening-campaign.sh --start 51 --rounds 50 --publish-every 50

Environment overrides:
  HARDENING_SKIP_FRONTEND=0
  HARDENING_SKIP_VERIFYE2E=0
  HARDENING_SKIP_DEMO_VERIFY=0
  HARDENING_DOC_PATH=docs/流程日志/持续优化执行说明.md
  HARDENING_AUDIT_DOC_PATH=docs/流程日志/优化审计记录.md
EOF
}

start_round=1
rounds=1
publish_every=50

while [ "$#" -gt 0 ]; do
  case "$1" in
    --start)
      start_round="${2:-}"
      shift 2
      ;;
    --rounds)
      rounds="${2:-}"
      shift 2
      ;;
    --publish-every)
      publish_every="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown option: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

case "$start_round" in
  ''|*[!0-9]*) echo "--start must be a positive integer" >&2; exit 1 ;;
esac
case "$rounds" in
  ''|*[!0-9]*) echo "--rounds must be a positive integer" >&2; exit 1 ;;
esac
case "$publish_every" in
  ''|*[!0-9]*) echo "--publish-every must be a positive integer" >&2; exit 1 ;;
esac

if [ "$start_round" -le 0 ] || [ "$rounds" -le 0 ] || [ "$publish_every" -le 0 ]; then
  echo "--start, --rounds, and --publish-every must be greater than 0" >&2
  exit 1
fi

: "${HARDENING_SKIP_FRONTEND:=0}"
: "${HARDENING_SKIP_VERIFYE2E:=0}"
: "${HARDENING_SKIP_DEMO_VERIFY:=0}"
: "${HARDENING_DOC_PATH:=docs/流程日志/持续优化执行说明.md}"
: "${HARDENING_AUDIT_DOC_PATH:=docs/流程日志/优化审计记录.md}"

wait_for_safe_load() {
  RESOURCE_LOAD_MAX_PERCENT="${RESOURCE_LOAD_MAX_PERCENT:-80}" \
  RESOURCE_LOAD_RESUME_PERCENT="${RESOURCE_LOAD_RESUME_PERCENT:-50}" \
    bash ./scripts/wait-for-cpu-idle.sh
}

repo_root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
if [ -z "$repo_root" ]; then
  echo "not inside a git repository" >&2
  exit 1
fi

cd "$repo_root"

list_worktree_changes() {
  {
    git diff --name-only -z HEAD --
    git ls-files -z --others --exclude-standard
  } | while IFS= read -r -d '' path; do
    [ -n "$path" ] || continue
    printf '%s\n' "$path"
  done | sort -u
}

worktree_has_path() {
  local target="$1"
  list_worktree_changes | grep -Fx -- "$target" >/dev/null
}

is_codegen_related_path() {
  local path="$1"
  case "$path" in
    admin-go/codegen/*)
      return 0
      ;;
    vue-vben-admin/apps/web-antd/src/adapter/component/*)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

collect_codegen_related_changes() {
  local path
  while IFS= read -r path; do
    [ -n "$path" ] || continue
    if is_codegen_related_path "$path"; then
      printf '%s\n' "$path"
    fi
  done < <(list_worktree_changes)
}

has_codegen_related_changes() {
  local first
  first="$(collect_codegen_related_changes | sed -n '1p')"
  [ -n "$first" ]
}

print_codegen_review_scope() {
  local path
  if ! has_codegen_related_changes; then
    return 0
  fi

  echo "[hardening] 本轮需回看的 codegen 相关代码："
  while IFS= read -r path; do
    [ -n "$path" ] || continue
    echo "[hardening]   - $path"
  done < <(collect_codegen_related_changes)
}

require_audit_doc_update_for_codegen_changes() {
  if ! has_codegen_related_changes; then
    return 0
  fi
  if [ ! -f "$HARDENING_AUDIT_DOC_PATH" ]; then
    echo "[hardening] 缺少审计文档: $HARDENING_AUDIT_DOC_PATH" >&2
    return 1
  fi
  if worktree_has_path "$HARDENING_AUDIT_DOC_PATH"; then
    return 0
  fi

  echo "[hardening] 检测到 codegen 相关改动，但未更新审计文档: $HARDENING_AUDIT_DOC_PATH" >&2
  print_codegen_review_scope >&2
  return 1
}

require_process_doc_update() {
  if ! worktree_has_path "scripts/run-hardening-campaign.sh"; then
    return 0
  fi
  if [ ! -f "$HARDENING_DOC_PATH" ]; then
    echo "[hardening] 缺少执行说明文档: $HARDENING_DOC_PATH" >&2
    return 1
  fi
  if worktree_has_path "$HARDENING_DOC_PATH"; then
    return 0
  fi

  echo "[hardening] run-hardening-campaign 脚本已变更，但未同步更新执行说明: $HARDENING_DOC_PATH" >&2
  return 1
}

load_env_file_without_override() {
  local env_file="$1"
  [ -f "$env_file" ] || return 0

  while IFS= read -r assignment; do
    [ -n "$assignment" ] || continue
    export "$assignment"
  done < <(python3 - "$env_file" <<'PY'
import os
import sys

path = sys.argv[1]
with open(path, "r", encoding="utf-8") as fh:
    for raw_line in fh:
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        if line.startswith("export "):
            line = line[len("export "):].strip()
        if "=" not in line:
            continue

        key, value = line.split("=", 1)
        key = key.strip()
        if not key or key in os.environ:
            continue

        value = value.strip()
        if len(value) >= 2 and value[0] == value[-1] and value[0] in "\"'":
            value = value[1:-1]
        else:
            for index, ch in enumerate(value):
                if ch == "#" and (index == 0 or value[index - 1].isspace()):
                    value = value[:index].rstrip()
                    break

        print(f"{key}={value}")
PY
  )
}

run_demo_comprehensive_verify() {
  local round="$1"

  if [ "$HARDENING_SKIP_DEMO_VERIFY" = "1" ]; then
    echo "[hardening] 命中 codegen 相关改动，但按 HARDENING_SKIP_DEMO_VERIFY=1 跳过 demo 应用全面测试"
    return 0
  fi

  if ! command -v mysql >/dev/null 2>&1; then
    echo "[hardening] 缺少 mysql 客户端，无法执行 demo 应用全面测试" >&2
    return 1
  fi
  if ! command -v python3 >/dev/null 2>&1; then
    echo "[hardening] 缺少 python3，无法读取 admin-go/.env 并执行 demo 应用全面测试" >&2
    return 1
  fi
  if ! command -v rsync >/dev/null 2>&1; then
    echo "[hardening] 缺少 rsync，无法准备 demo 应用临时工作区" >&2
    return 1
  fi
  if ! command -v pnpm >/dev/null 2>&1; then
    echo "[hardening] 缺少 pnpm，无法执行 demo 前端类型校验" >&2
    return 1
  fi

  echo "[hardening] 检测到 codegen 相关改动，执行 demo 应用全面测试"

  (
    set -euo pipefail

    local temp_root admin_temp_root frontend_temp_root temp_config env_file existing_tables has_system_menu created_system_menu
    temp_root="$(mktemp -d "${TMPDIR:-/tmp}/baseadmin-demo-verify-XXXXXX")"
    admin_temp_root="$temp_root/admin-go"
    frontend_temp_root="$temp_root/vue-vben-admin"
    temp_config="$admin_temp_root/codegen/codegen.verify.yaml"
    env_file="$repo_root/admin-go/.env"
    created_system_menu=0

    cleanup() {
      if [ -n "${MYSQL_HOST:-}" ] && [ -n "${MYSQL_PORT:-}" ] && [ -n "${MYSQL_USER:-}" ] && [ -n "${MYSQL_DATABASE:-}" ]; then
        MYSQL_PWD="${MYSQL_PASSWORD:-}" mysql \
          --host="$MYSQL_HOST" \
          --port="$MYSQL_PORT" \
          --user="$MYSQL_USER" \
          "$MYSQL_DATABASE" \
          -e "DROP TABLE IF EXISTS \`demo_article\`, \`demo_category\`, \`demo_tag\`;" >/dev/null 2>&1 || true
        if [ "$created_system_menu" = "1" ]; then
          MYSQL_PWD="${MYSQL_PASSWORD:-}" mysql \
            --host="$MYSQL_HOST" \
            --port="$MYSQL_PORT" \
            --user="$MYSQL_USER" \
            "$MYSQL_DATABASE" \
            -e "DROP TABLE IF EXISTS \`system_menu\`;" >/dev/null 2>&1 || true
        fi
      fi
      rm -rf "$temp_root"
    }
    trap cleanup EXIT

    [ -f "$env_file" ] || {
      echo "[hardening] 缺少 admin-go/.env，无法执行 demo 应用全面测试" >&2
      exit 1
    }
    load_env_file_without_override "$env_file"

    : "${MYSQL_HOST:?missing MYSQL_HOST in admin-go/.env}"
    : "${MYSQL_PORT:?missing MYSQL_PORT in admin-go/.env}"
    : "${MYSQL_USER:?missing MYSQL_USER in admin-go/.env}"
    : "${MYSQL_DATABASE:?missing MYSQL_DATABASE in admin-go/.env}"

    existing_tables="$(MYSQL_PWD="${MYSQL_PASSWORD:-}" mysql \
      --batch \
      --skip-column-names \
      --host="$MYSQL_HOST" \
      --port="$MYSQL_PORT" \
      --user="$MYSQL_USER" \
      "$MYSQL_DATABASE" \
      -e "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME IN ('demo_category', 'demo_article', 'demo_tag');")"
    if [ -n "$existing_tables" ]; then
      echo "[hardening] 目标库已存在 demo 验证表，拒绝覆盖：" >&2
      printf '%s\n' "$existing_tables" >&2
      exit 1
    fi

    MYSQL_PWD="${MYSQL_PASSWORD:-}" mysql \
      --host="$MYSQL_HOST" \
      --port="$MYSQL_PORT" \
      --user="$MYSQL_USER" \
      "$MYSQL_DATABASE" < "$repo_root/admin-go/codegen/sql/demo.sql"

    has_system_menu="$(MYSQL_PWD="${MYSQL_PASSWORD:-}" mysql \
      --batch \
      --skip-column-names \
      --host="$MYSQL_HOST" \
      --port="$MYSQL_PORT" \
      --user="$MYSQL_USER" \
      "$MYSQL_DATABASE" \
      -e "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'system_menu';")"
    if [ "$has_system_menu" = "0" ]; then
      MYSQL_PWD="${MYSQL_PASSWORD:-}" mysql \
        --host="$MYSQL_HOST" \
        --port="$MYSQL_PORT" \
        --user="$MYSQL_USER" \
        "$MYSQL_DATABASE" \
        -e "CREATE TABLE \`system_menu\` (
          \`id\` bigint unsigned NOT NULL,
          \`parent_id\` bigint unsigned NOT NULL DEFAULT '0',
          \`title\` varchar(50) NOT NULL,
          \`type\` tinyint NOT NULL DEFAULT '1',
          \`path\` varchar(200) DEFAULT NULL,
          \`component\` varchar(200) DEFAULT NULL,
          \`permission\` varchar(100) DEFAULT NULL,
          \`icon\` varchar(100) DEFAULT NULL,
          \`sort\` int NOT NULL DEFAULT '0',
          \`is_show\` tinyint(1) NOT NULL DEFAULT '1',
          \`is_cache\` tinyint(1) NOT NULL DEFAULT '0',
          \`status\` tinyint(1) NOT NULL DEFAULT '1',
          \`created_by\` bigint unsigned DEFAULT NULL,
          \`dept_id\` bigint unsigned DEFAULT NULL,
          \`created_at\` datetime DEFAULT NULL,
          \`updated_at\` datetime DEFAULT NULL,
          \`deleted_at\` datetime DEFAULT NULL,
          PRIMARY KEY (\`id\`),
          KEY \`idx_parent_id\` (\`parent_id\`),
          KEY \`idx_deleted_at\` (\`deleted_at\`)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"
      created_system_menu=1
    fi

    mkdir -p "$admin_temp_root/codegen"
    cp "$repo_root/admin-go/go.mod" "$admin_temp_root/go.mod"
    cp "$repo_root/admin-go/go.sum" "$admin_temp_root/go.sum"
    cp "$env_file" "$admin_temp_root/.env"
    rsync -a "$repo_root/admin-go/utility/" "$admin_temp_root/utility/"

    rsync -a \
      --exclude node_modules \
      --exclude .pnpm \
      --exclude .git \
      --exclude dist \
      --exclude coverage \
      --exclude .turbo \
      "$repo_root/vue-vben-admin/" "$frontend_temp_root/"

    while IFS= read -r source_node_modules; do
      [ -n "$source_node_modules" ] || continue
      local rel_path
      rel_path="${source_node_modules#"$repo_root/vue-vben-admin/"}"
      mkdir -p "$frontend_temp_root/$(dirname "$rel_path")"
      if [ ! -e "$frontend_temp_root/$rel_path" ]; then
        ln -s "$source_node_modules" "$frontend_temp_root/$rel_path"
      fi
    done < <(find "$repo_root/vue-vben-admin" \( -type d -o -type l \) -name node_modules)

    if [ -e "$repo_root/vue-vben-admin/.pnpm" ]; then
      ln -s "$repo_root/vue-vben-admin/.pnpm" "$frontend_temp_root/.pnpm"
    fi

    if [ ! -e "$frontend_temp_root/apps/web-antd/src/api/system/dict/index.ts" ]; then
      mkdir -p "$frontend_temp_root/apps/web-antd/src/api/system/dict"
      cat > "$frontend_temp_root/apps/web-antd/src/api/system/dict/index.ts" <<'EOF'
export async function getDictByType(_dictType: string): Promise<Array<{ label: string; value: number | string }>> {
  return [];
}
EOF
    fi

    cat > "$temp_config" <<'EOF'
database:
  host: ${MYSQL_HOST}
  port: ${MYSQL_PORT}
  user: ${MYSQL_USER}
  password: ${MYSQL_PASSWORD}
  dbname: ${MYSQL_DATABASE}

backend:
  output: ../app/

frontend:
  output: ../../vue-vben-admin/apps/web-antd/src/

allowed_apps:
  - system
  - upload
  - demo

skip_fields:
  - created_at
  - updated_at
  - deleted_at
  - created_by
  - dept_id

menu_apps:
  system:
    title: 系统管理
    icon: SettingOutlined
    sort: 10
  upload:
    title: 上传管理
    icon: CloudUploadOutlined
    sort: 20
  demo:
    title: 演示管理
    icon: AppstoreOutlined
    sort: 90

allow_missing_dict_module: false
operation_log: false
EOF

    if [ "$HARDENING_SKIP_VERIFYE2E" != "1" ] && [ $((round % publish_every)) -ne 0 ]; then
      (
        cd "$repo_root/admin-go/codegen"
        ../../scripts/run-go-task-with-limits.sh go run ./cmd/verifye2e --stage all
      )
    fi

    (
      cd "$repo_root/admin-go/codegen"
      ../../scripts/run-go-task-with-limits.sh go run . \
        --table demo_category,demo_article,demo_tag \
        --with-dao \
        --force \
        --config "$temp_config"
      ../../scripts/run-go-task-with-limits.sh go run . \
        --table demo_category,demo_article,demo_tag \
        --only menu \
        --dry-run \
        --config "$temp_config" \
        --manifest-out "$temp_root/demo-menu-manifest.json"
    )

    (
      cd "$admin_temp_root"
      "$repo_root/scripts/run-go-task-with-limits.sh" go test ./app/demo/...
    )

    "$repo_root/scripts/run-node-task-with-limits.sh" bash -lc \
      "cd \"$frontend_temp_root/apps/web-antd\" && NODE_OPTIONS=--max-old-space-size=4096 ../../node_modules/.bin/vue-tsc --noEmit --skipLibCheck"
  )

  echo "[hardening] demo 应用全面测试通过"
}

update_doc_checkpoint() {
  local round="$1"
  [ -f "$HARDENING_DOC_PATH" ] || return 0

  if grep -q '^-\s*最近一次 50 轮关口验证：' "$HARDENING_DOC_PATH"; then
    sed -i "s#^-\s*最近一次 50 轮关口验证：.*#- 最近一次 50 轮关口验证：第 ${round} 轮#" "$HARDENING_DOC_PATH"
  fi
}

run_round() {
  local round="$1"

  wait_for_safe_load
  echo "第${round}轮开始"
  if has_codegen_related_changes; then
    echo "[hardening] 检测到本轮存在 codegen 相关改动"
    print_codegen_review_scope
  fi
  bash ./scripts/verify-baseadmin-scope.sh
  bash ./scripts/verify-vben-pages.sh

  (
    cd admin-go
    ../scripts/run-go-task-with-limits.sh go test ./...
  )

  (
    cd admin-go/codegen
    ../../scripts/run-go-task-with-limits.sh go test ./...
    ../../scripts/run-go-task-with-limits.sh go run verify_codegen.go
  )

  if [ "$HARDENING_SKIP_FRONTEND" != "1" ]; then
    bash ./scripts/run-node-task-with-limits.sh pnpm -C vue-vben-admin -F @vben/web-antd typecheck
    bash ./scripts/run-node-task-with-limits.sh pnpm -C vue-vben-admin test:web-antd
  fi

  if has_codegen_related_changes; then
    run_demo_comprehensive_verify "$round"
    require_audit_doc_update_for_codegen_changes
    print_codegen_review_scope
  fi

  require_process_doc_update

  if [ $((round % publish_every)) -ne 0 ]; then
    return 0
  fi

  if [ "$HARDENING_SKIP_VERIFYE2E" != "1" ]; then
    (
      cd admin-go/codegen
      ../../scripts/run-go-task-with-limits.sh go run ./cmd/verifye2e --stage all
    )
  fi

  update_doc_checkpoint "$round"
  require_audit_doc_update_for_codegen_changes
  require_process_doc_update

  if [ -z "$(git status --short)" ]; then
    echo "[hardening] 第${round}轮无文件变更，跳过发布"
    return 0
  fi

  git add -A
  ./scripts/feature-publish.sh --all "chore(hardening): complete round ${round}"
}

for ((offset = 0; offset < rounds; offset++)); do
  run_round "$((start_round + offset))"
done
