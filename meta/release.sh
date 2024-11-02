#!/usr/bin/env bash
set -euo pipefail

EBRO_VERSION="$(git rev-parse --verify HEAD)"
export EBRO_VERSION

function main {
    rm -rf dist
    GOOS=linux GOARCH=arm64 build
    GOOS=linux GOARCH=amd64 build
    GOOS=darwin GOARCH=arm64 build
    GOOS=darwin GOARCH=amd64 build
}

function build {
    echo "Building ${GOOS} ${GOARCH}"
    mkdir -p "dist/${GOOS}_${GOARCH}"
    go build \
        -ldflags "-X github.com/sirikon/ebro/cmd/ebro/cli.version=${EBRO_VERSION}" \
        -o "dist/${GOOS}_${GOARCH}/ebro" \
        cmd/ebro/main.go
}

main "$@"
