#!/usr/bin/env bash
set -euo pipefail

root="$(pwd)"
mkdir -p out
cd src
CGO_ENABLED=0 go build -o "$root/out/ebro-e2e" -cover cmd/ebro/main.go
