#!/usr/bin/env bash
set -euo pipefail

GO_VERSION="$(grep go <.tool-versions | sed -E 's/^go (.*)$/\1/g')"
TAG="ebro-build:$(base64 </dev/urandom | tr -d '[A-Z]+/' | head -c 8 || true)"

rm -rf out
docker build \
    -t "$TAG" \
    --file meta/docker/_/Dockerfile \
    --build-arg "GO_VERSION=${GO_VERSION}" \
    .
container_id="$(docker create "$TAG")"
docker cp "$container_id:/out" ./out
docker rm -f "$container_id"
docker rmi "$TAG"
