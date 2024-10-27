#!/usr/bin/env bash
set -euo pipefail

mkdir -p out
go build -o out/ebro cmd/ebro/main.go
