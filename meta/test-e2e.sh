#!/usr/bin/env bash
set -euo pipefail

function main {
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
