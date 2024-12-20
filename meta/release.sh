#!/usr/bin/env bash
set -euo pipefail

EBRO_COMMIT="$(git rev-parse --verify HEAD)"
export EBRO_COMMIT

function main {
    rm -rf dist
    GOOS=linux GOARCH=arm64 build "Linux__aarch64"
    GOOS=linux GOARCH=amd64 build "Linux__x86_64"
    GOOS=darwin GOARCH=arm64 build "Darwin__arm64"
}

function build {
    dest="$(pwd)/dist/${1}/ebro"
    timestamp="$(date +%s)"
    echo "Building ${GOOS} ${GOARCH}"
    mkdir -p "$(dirname "$dest")"
    (
        cd src
        go build \
            -ldflags \
            "-X github.com/sirikon/ebro/internal/constants.version=${EBRO_COMMIT} \
            -X github.com/sirikon/ebro/internal/constants.commit=${EBRO_COMMIT} \
            -X github.com/sirikon/ebro/internal/constants.timestamp=${timestamp}" \
            -o "$dest" \
            cmd/ebro/main.go
    )

}

main "$@"
