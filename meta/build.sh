#!/usr/bin/env bash
set -euo pipefail

EBRO_BIN="${EBRO_BIN:-"$(pwd)/out/ebro"}"
EBRO_VERSION="${EBRO_VERSION:-"dev"}"
EBRO_COMMIT="${EBRO_COMMIT:-"HEAD"}"
EBRO_TIMESTAMP="${EBRO_TIMESTAMP:-"0"}"

cd src
mkdir -p "$(dirname "$EBRO_BIN")"
CGO_ENABLED=0 go build \
  -ldflags "-X github.com/sirikon/ebro/internal/constants.version=${EBRO_VERSION} \
            -X github.com/sirikon/ebro/internal/constants.commit=${EBRO_COMMIT} \
            -X github.com/sirikon/ebro/internal/constants.timestamp=${EBRO_TIMESTAMP}" \
  -o "$EBRO_BIN" \
  "$@" \
  cmd/ebro/main.go
