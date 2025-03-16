#!/usr/bin/env bash
set -euo pipefail

function main {
  enable_coverage="false"
  if [ "${1:-}" == "coverage" ]; then
    enable_coverage="true"
  fi

  if [ "${EBRO_BIN:-}" == "" ]; then
    EBRO_BIN="$(pwd)/out/ebro-e2e"
    export EBRO_BIN
    log "Building ebro for e2e tests"
    if [ "$enable_coverage" == "true" ]; then
      ./meta/build.sh -cover
    else
      ./meta/build.sh
    fi
  fi

  log "Running e2e tests"
  ./meta/python/ensure-venv.sh
  export PYTHONPATH=src
  (
    cd tests
    if [ "$enable_coverage" == "true" ]; then
      mkdir -p .coverage
      rm -rf .coverage/*
    fi
    ./../meta/python/_/.venv/bin/python -m unittest discover src "*_test.py"
    if [ "$enable_coverage" == "true" ]; then
      log "Collecting coverage data"
      go tool covdata textfmt -i=.coverage -o=.coverage-profile
    fi
  )

  if [ "$enable_coverage" == "true" ]; then
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
