#!/usr/bin/env bash
set -euo pipefail

if [ ! -f "out/dist/RELEASE" ]; then
    echo "Nothing to release"
    exit 0
fi
echo "Done"
