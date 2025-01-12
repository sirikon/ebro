#!/usr/bin/env bash
set -euo pipefail

function main {
  log "Running unit tests"
  (
    cd src
    go test ./...
  )

  log "Building ebro for e2e tests"
  ./meta/build.sh

  log "Running e2e tests"
  ./meta/python/ensure-venv.sh
  export PYTHONPATH=src
  cd tests
  ./../meta/python/_/.venv/bin/python -m unittest discover src "*_test.py"
}

function log {
  printf "### %s\n" "$1"
}

main "$@"
