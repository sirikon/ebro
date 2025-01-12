#!/usr/bin/env bash
set -euo pipefail

TAG="ebro-run:$(head -c 512 </dev/urandom | base64 | tr -d '[A-Z]+/' | head -c 8)"

docker build \
  -t "$TAG" \
  --file meta/docker/_/Dockerfile \
  .
set +e
docker run -it --rm \
  -v ./:/w \
  -v /var/run/docker.sock:/var/run/docker.sock \
  "$TAG" "$@"
exitcode="$?"
set -e
docker rmi "$TAG"
exit "$exitcode"
