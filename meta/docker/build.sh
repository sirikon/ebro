#!/usr/bin/env bash
set -euo pipefail

TARGET="${1:-all}" # dist, website, all
GO_VERSION="$(grep go <.tool-versions | sed -E 's/^go (.*)$/\1/g')"
PYTHON_VERSION="$(grep python <.tool-versions | sed -E 's/^python (.*)$/\1/g')"
POETRY_VERSION="$(grep poetry <.tool-versions | sed -E 's/^poetry (.*)$/\1/g')"
TAG="ebro-build:$(base64 </dev/urandom | tr -d '[A-Z]+/' | head -c 8 || true)"

extra_args=()
if [ -n "${ACTIONS_RUNTIME_TOKEN}" ]; then
    echo "GitHub Actions detected. Enabling GHA Docker cache."
    extra_args+=(--cache-to type=gha --cache-from type=gha)
fi

rm -rf out
docker build \
    -t "$TAG" \
    --file meta/docker/_/Dockerfile \
    --build-arg "GO_VERSION=${GO_VERSION}" \
    --build-arg "PYTHON_VERSION=${PYTHON_VERSION}" \
    --build-arg "POETRY_VERSION=${POETRY_VERSION}" \
    --target "$TARGET" \
    "${extra_args[@]}" \
    .
container_id="$(docker create "$TAG")"
docker cp "$container_id:/out" ./out
docker rm -f "$container_id"
docker rmi "$TAG"
