#!/usr/bin/env bash
set -euo pipefail

./meta/python/ensure-venv.sh
PYTHONPATH="$(realpath "$(dirname "${BASH_SOURCE[0]}")")/_/src"
export PYTHONPATH
exec ./meta/python/_/.venv/bin/python -m flask \
    --app ebro_website.app \
    run --port 8000 --debug
