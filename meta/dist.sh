#!/usr/bin/env bash
set -euo pipefail

EBRO_COMMIT="$(git rev-parse --verify HEAD)"
EBRO_VERSION="nightly_$EBRO_COMMIT"
EBRO_TIMESTAMP="$(date +%s)"
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
  sed -i -E "/^ +# gen:EBRO_SUMS/d" "out/dist/${EBRO_VERSION}/ebrow"

  echo "$EBRO_VERSION" >out/dist/VERSION
  if [ -n "$EBRO_RELEASE" ]; then
    echo "$EBRO_RELEASE" >out/dist/RELEASE
    echo "$EBRO_COMMIT" >out/dist/RELEASE_COMMIT
  fi
}

function build {
  variant="$1"
  filename="ebro-${variant}"
  EBRO_BIN="$(pwd)/out/dist/${EBRO_VERSION}/${filename}"

  echo "Building ${GOOS} ${GOARCH}"
  (
    export EBRO_BIN
    export EBRO_VERSION
    export EBRO_COMMIT
    export EBRO_TIMESTAMP
    ./meta/build.sh
  )

  (
    cd "$(dirname "$EBRO_BIN")"
    sha256sum "$filename" >"$filename.sha256"
    hash="$(sed -E 's/^([a-z0-9]+).*$/\1/' <"$filename.sha256")"
    sed -i -E "s/^( +)# gen:EBRO_SUMS/\1[\"$variant\"]=\"${hash}\"\n\1# gen:EBRO_SUMS/" "ebrow"
  )
}

main "$@"
