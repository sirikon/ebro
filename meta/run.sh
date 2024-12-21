#!/usr/bin/env bash
set -euo pipefail

./meta/build.sh
cd playground
exec ../out/ebro "$@"
