#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")/../api"
go generate -n -v ./...
go generate -x ./...
echo "done go generate ./..."
cat -n /app/api/graphql/models/generated.go
if [[ "$(git status -s 2>/dev/null | head -1)" != "" ]]; then
  # shellcheck disable=SC2016
  echo '--- FAIL: The generated API code is out of sync with the recent changes. Please run `go generate ./...` under `./api` to regenerate it and commit it to this branch.'
  echo 'These are the changes:'
  git status -s
  exit 1
fi

echo '--- PASS: All generated code is in sync with the project.'
