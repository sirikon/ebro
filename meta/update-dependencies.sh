#!/usr/bin/env bash
set -euo pipefail

echo "Updating Go dependencies"
(
  cd src
  go get -u ./...
  go mod tidy
)

echo "Updating Python dependencies"
./meta/python/poetry.sh update
