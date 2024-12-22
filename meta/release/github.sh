#!/usr/bin/env bash
set -euo pipefail

if [ ! -f "out/dist/RELEASE" ]; then
    echo "Nothing to release"
fi
echo "Done"
