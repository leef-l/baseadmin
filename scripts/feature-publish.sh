#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./scripts/feature-publish.sh [--all] "feat(scope): summary"

Behavior:
  1. By default, publish already-staged changes only
  2. Create one git commit
  3. Push current branch to origin (or GIT_REMOTE)

Optional env:
  GIT_REMOTE=origin
  GIT_BRANCH_OVERRIDE=main
EOF
}

stage_all=0
while [ "$#" -gt 0 ]; do
  case "$1" in
    --all)
      stage_all=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    --)
      shift
      break
      ;;
    -*)
      echo "unknown option: $1" >&2
      usage
      exit 1
      ;;
    *)
      break
      ;;
  esac
done

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
message="$(printf '%s' "$1" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')"

if [ -z "$branch" ]; then
  echo "current branch is empty, refusing to push from detached HEAD" >&2
  exit 1
fi

if [ -z "$message" ]; then
  echo "commit message is empty after trimming whitespace" >&2
  exit 1
fi

if [ -n "$(git diff --name-only --diff-filter=U)" ]; then
  echo "unresolved merge conflicts detected, refusing to publish" >&2
  exit 1
fi

if ! git remote get-url "$remote" >/dev/null 2>&1; then
  echo "git remote '$remote' not found" >&2
  exit 1
fi

if [ "$stage_all" = "1" ]; then
  git add -A
else
  if [ -z "$(git diff --cached --name-only)" ]; then
    echo "no staged changes to publish; stage exact paths first or rerun with --all" >&2
    exit 1
  fi

  if [ -n "$(git diff --name-only)" ] || [ -n "$(git ls-files --others --exclude-standard)" ]; then
    echo "unstaged or untracked changes detected; refusing to publish a mixed worktree" >&2
    echo "stage exact paths first, or rerun with --all if every change belongs to this delivery" >&2
    exit 1
  fi
fi

if [ -z "$(git diff --cached --name-only)" ]; then
  echo "no staged changes to publish"
  exit 0
fi

git commit -m "$message"
git push "$remote" "$branch"

echo "published feature commit to $remote/$branch"
