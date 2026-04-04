#!/bin/sh
set -eu

cd /app

app_name="${1:-}"

if [ -z "$app_name" ]; then
  echo "usage: start-go-app.sh <system|upload>"
  exit 1
fi

case "$app_name" in
  system)
    exec ./system
    ;;
  upload)
    exec ./upload
    ;;
  *)
    echo "unsupported app: $app_name"
    exit 1
    ;;
esac
