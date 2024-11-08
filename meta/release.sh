#!/usr/bin/env bash
set -euo pipefail

EBRO_VERSION="$(git rev-parse --verify HEAD)"
export EBRO_VERSION

function main {
    rm -rf dist
    GOOS=linux GOARCH=arm64 build "Linux__aarch64"
    GOOS=linux GOARCH=amd64 build "Linux__x86_64"
    GOOS=darwin GOARCH=arm64 build "Darwin__arm64"
}

function build {
    dest="dist/${1}/ebro"
    echo "Building ${GOOS} ${GOARCH}"
    mkdir -p "$(dirname "$dest")"
    go build \
        -ldflags "-X github.com/sirikon/ebro/cmd/ebro/cli.version=${EBRO_VERSION}" \
        -o "$dest" \
        cmd/ebro/main.go
}

main "$@"
