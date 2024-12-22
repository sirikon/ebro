#!/usr/bin/env bash
set -euo pipefail

if [ ! -f "out/dist/RELEASE" ]; then
    echo "Nothing to release"
    exit 0
fi

EBRO_VERSION="$(cat out/dist/RELEASE)"
EBRO_COMMIT="$(cat out/dist/RELEASE_COMMIT)"

gh release create "$EBRO_VERSION" "out/dist/${EBRO_VERSION}/"* \
    --target "${EBRO_COMMIT}" \
    --verify-tag \
    --draft
