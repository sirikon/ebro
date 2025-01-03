#!/usr/bin/env bash
set -euo pipefail

TARGET="${1:-all}" # dist, website, all
GO_VERSION="$(grep go <.tool-versions | sed -E 's/^go (.*)$/\1/g')"
PYTHON_VERSION="$(grep python <.tool-versions | sed -E 's/^python (.*)$/\1/g')"
POETRY_VERSION="$(grep poetry <.tool-versions | sed -E 's/^poetry (.*)$/\1/g')"
TAG="ebro-build:$(head -c 512 </dev/urandom | base64 | tr -d '[A-Z]+/' | head -c 8)"

rm -rf out
docker build \
    -t "$TAG" \
    --file meta/docker/_/Dockerfile \
    --build-arg "GO_VERSION=${GO_VERSION}" \
    --build-arg "PYTHON_VERSION=${PYTHON_VERSION}" \
    --build-arg "POETRY_VERSION=${POETRY_VERSION}" \
    --target "$TARGET" \
    .
container_id="$(docker create "$TAG")"
docker cp "$container_id:/out" ./out
docker rm -f "$container_id"
docker rmi "$TAG"
