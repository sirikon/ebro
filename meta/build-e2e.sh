#!/usr/bin/env bash
set -euo pipefail

mkdir -p "$(dirname "$EBRO_BIN")"
cd src
CGO_ENABLED=0 go build -o "$EBRO_BIN" -cover cmd/ebro/main.go
