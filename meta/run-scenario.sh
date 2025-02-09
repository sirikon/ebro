#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -eq 0 ]; then
  echo "Usage: ./meta/run-scenario.sh <scenario> <...args>"
  echo "scenarios:"
  (cd "tests/scenarios" && ls)
  exit 1
fi

EBRO_BIN="$(pwd)/out/ebro"
export EBRO_BIN

./meta/build.sh
cd "tests/scenarios/$1"
exec "$EBRO_BIN" "${@:2}"
