#!/usr/bin/env bash
set -euo pipefail

function main {
  log "Running unit tests"
  (
    cd src
    go test ./...
  )
}

function log {
  printf "### %s\n" "$1"
}

main "$@"
