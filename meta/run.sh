#!/usr/bin/env bash
set -euo pipefail

EBRO_BIN="$(pwd)/out/ebro"
export EBRO_BIN

./meta/build.sh
cd playground
exec "$EBRO_BIN" "$@"
