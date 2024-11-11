#!/usr/bin/env bash
set -euo pipefail

function main {
    go test ./...

    ./meta/build.sh

    cd tests

    if [ ! -d ".venv" ]; then
        export POETRY_VIRTUALENVS_IN_PROJECT="true"
        poetry install
    fi

    export PYTHONPATH=src
    ./.venv/bin/python -m unittest discover src "*_test.py"

}

main "$@"
