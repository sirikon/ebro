#!/usr/bin/env bash
set -euo pipefail

function main {
    url="https://ebro.srk.bz/$(uname -s)__$(uname -m)/ebro"
    dest="${EBRO_BIN}"
    rm -f "$dest"
    mkdir -p "$(dirname "$dest")"
    curl --fail -L -o "$dest" "$url"
    chmod +x "$dest"
}

main "$@"
