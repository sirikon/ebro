#!/usr/bin/env bash
set -euo pipefail

function main {
  log "Building ebro for e2e tests"
  ./meta/build.sh
  while true; do
    ./meta/test-e2e.sh
  done
}

function log {
  printf "### %s\n" "$1"
}

main "$@"
