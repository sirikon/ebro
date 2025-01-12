#!/usr/bin/env bash
set -euo pipefail

function main {
  ./meta/test-unit.sh
  log "Building ebro for e2e tests"
  ./meta/build.sh
  ./meta/test-e2e.sh
}

function log {
  printf "### %s\n" "$1"
}

main "$@"
