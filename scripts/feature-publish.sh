#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/feature-publish.sh "feat(scope): summary"

Behavior:
  1. Stage all current changes
  2. Create one git commit
  3. Push current branch to origin (or GIT_REMOTE)

Optional env:
  GIT_REMOTE=origin
  GIT_BRANCH_OVERRIDE=main
EOF
}

if [ "${1:-}" = "" ]; then
  usage
  exit 1
fi

repo_root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
if [ -z "$repo_root" ]; then
  echo "not inside a git repository" >&2
  exit 1
fi

cd "$repo_root"

if [ -z "$(git status --short)" ]; then
  echo "no changes to publish"
  exit 0
fi

remote="${GIT_REMOTE:-origin}"
branch="${GIT_BRANCH_OVERRIDE:-$(git branch --show-current)}"
message="$1"

if [ -z "$branch" ]; then
  echo "current branch is empty, refusing to push from detached HEAD" >&2
  exit 1
fi

if ! git remote get-url "$remote" >/dev/null 2>&1; then
  echo "git remote '$remote' not found" >&2
  exit 1
fi

git add -A
git commit -m "$message"
git push "$remote" "$branch"

echo "published feature commit to $remote/$branch"
