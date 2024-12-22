#!/usr/bin/env bash
set -euo pipefail

EBRO_COMMIT="$(git rev-parse --verify HEAD)"
export EBRO_COMMIT
EBRO_VERSION="nightly/$EBRO_COMMIT"
EBRO_RELEASE=""
if [ "$(git tag --points-at "$EBRO_COMMIT" | wc -l)" == "1" ]; then
    EBRO_VERSION="$(git tag --points-at "$EBRO_COMMIT")"
    EBRO_RELEASE="$EBRO_VERSION"
fi

function main {
    rm -rf out/dist
    mkdir -p "out/dist/${EBRO_VERSION}"

    sed -E "s/^.*# gen:EBRO_VERSION$/EBRO_VERSION=\"${EBRO_VERSION}\"/" <scripts/ebrow >"out/dist/${EBRO_VERSION}/ebrow"

    GOOS=linux GOARCH=arm64 build "Linux__aarch64"
    GOOS=linux GOARCH=amd64 build "Linux__x86_64"
    GOOS=darwin GOARCH=arm64 build "Darwin__arm64"
    if [ -n "$EBRO_RELEASE" ]; then
        echo "$EBRO_RELEASE" >out/dist/RELEASE
    fi
}

function build {
    variant="$1"
    dest="$(pwd)/out/dist/${EBRO_VERSION}/ebro-${variant}"
    timestamp="$(date +%s)"
    echo "Building ${GOOS} ${GOARCH}"
    mkdir -p "$(dirname "$dest")"
    (
        cd src
        CGO_ENABLED=0 go build \
            -ldflags \
            "-X github.com/sirikon/ebro/internal/constants.version=${EBRO_VERSION} \
            -X github.com/sirikon/ebro/internal/constants.commit=${EBRO_COMMIT} \
            -X github.com/sirikon/ebro/internal/constants.timestamp=${timestamp}" \
            -o "$dest" \
            cmd/ebro/main.go
    )
    sha256sum "$dest" | sed -E 's/^([a-z0-9]+).*$/\1/' >"$dest.sha256"
}

main "$@"
