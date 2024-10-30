#!/usr/bin/env bash
set -euo pipefail

cd playground
go run ../cmd/ebro/main.go "$@"
