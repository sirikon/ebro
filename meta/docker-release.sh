#!/usr/bin/env bash
set -euo pipefail

rm -rf dist
tag="ebro_builder_$(date +'%Y%m%d_%H%M%S')"
docker build -t "${tag}" .
container_id="$(docker create "${tag}")"
docker cp "$container_id:/dist" ./dist
docker rm -f "$container_id"
