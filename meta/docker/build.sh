#!/usr/bin/env bash
set -euo pipefail

TARGET="${1:-all}" # dist, website, all
GO_VERSION="$(./meta/tool-version.sh go)"
PYTHON_VERSION="$(./meta/tool-version.sh python)"
POETRY_VERSION="$(./meta/tool-version.sh poetry)"
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
