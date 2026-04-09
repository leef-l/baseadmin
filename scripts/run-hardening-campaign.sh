#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/run-hardening-campaign.sh [--start N] [--rounds N] [--publish-every N]

Purpose:
  Execute repeated hardening validation rounds with bounded resource usage.
  Each round prints the round number. Every publish interval runs codegen demo
  verification, updates the hardening doc checkpoint, and publishes changes if
  the worktree is dirty.

Examples:
  ./scripts/run-hardening-campaign.sh --start 2 --rounds 1
  ./scripts/run-hardening-campaign.sh --start 51 --rounds 50 --publish-every 50

Environment overrides:
  HARDENING_SKIP_FRONTEND=0
  HARDENING_SKIP_VERIFYE2E=0
  HARDENING_DOC_PATH=docs/持续优化执行说明.md
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
: "${HARDENING_DOC_PATH:=docs/持续优化执行说明.md}"

repo_root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
if [ -z "$repo_root" ]; then
  echo "not inside a git repository" >&2
  exit 1
fi

cd "$repo_root"

update_doc_checkpoint() {
  local round="$1"
  [ -f "$HARDENING_DOC_PATH" ] || return 0

  if grep -q '^-\s*最近一次 50 轮关口验证：' "$HARDENING_DOC_PATH"; then
    sed -i "s#^-\s*最近一次 50 轮关口验证：.*#- 最近一次 50 轮关口验证：第 ${round} 轮#" "$HARDENING_DOC_PATH"
  fi
}

run_round() {
  local round="$1"

  echo "第${round}轮开始"
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
