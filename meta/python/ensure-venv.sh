#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/_"
if [ ! -d ".venv" ]; then
  ./../poetry.sh install
fi
