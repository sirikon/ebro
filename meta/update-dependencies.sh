#!/usr/bin/env bash
set -euo pipefail

cd src
go get -u ./...
go mod tidy
