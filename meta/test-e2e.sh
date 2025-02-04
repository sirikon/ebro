#!/usr/bin/env bash
set -euo pipefail

function main {
  log "Building ebro for e2e tests"
  ./meta/build-e2e.sh
  log "Running e2e tests"
  ./meta/python/ensure-venv.sh
  export PYTHONPATH=src
  (
    cd tests
    mkdir -p .coverage
    rm -rf .coverage/*
    ./../meta/python/_/.venv/bin/python -m unittest discover src "*_test.py"
    go tool covdata textfmt -i=.coverage -o=.coverage-profile
  )

  if [ "${1:-}" == "coverage" ]; then
    (
      cd src
      go tool cover -html=../tests/.coverage-profile
    )
  fi
}

function log {
  printf "### %s\n" "$1"
}

main "$@"
