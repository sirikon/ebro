#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/_"
export POETRY_VIRTUALENVS_IN_PROJECT="true"
poetry "$@"
