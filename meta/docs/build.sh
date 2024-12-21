#!/usr/bin/env bash
set -euo pipefail

PYTHONPATH="$(realpath "$(dirname "${BASH_SOURCE[0]}")")/_/src"
export PYTHONPATH
exec ./meta/python/_/.venv/bin/python -m ebro_docs.freeze
