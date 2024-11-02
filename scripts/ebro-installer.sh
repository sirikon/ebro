#!/usr/bin/env bash
set -euo pipefail

function main {
    os="$(detect-os)"
    arch="$(detect-arch)"
    url="https://ebro.srk.bz/${os}_${arch}/ebro"
    dest="${EBRO_BIN}"
    rm -rf "$(dirname "$dest")"
    mkdir -p "$(dirname "$dest")"
    curl -L -o "$dest" "$url"
    chmod +x "$dest"
}

function detect-os {
    uname -s | tr '[:upper:]' '[:lower:]'
}

function detect-arch {
    value="$(arch)"
    if [ "$value" == "x86_64" ]; then
        echo "amd64"
        return
    fi
    echo "$value"
}

main "$@"
