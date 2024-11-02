#!/usr/bin/env bash
set -euo pipefail

cd playground
exec go run ../cmd/ebro/main.go "$@"
