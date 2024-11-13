#!/usr/bin/env bash
set -euo pipefail

docker run --rm \
    -v ./entrypoint.sh:/entrypoint.sh \
    -v ./content:/content:ro \
    -v ./httpd.conf:/usr/local/apache2/conf/httpd.conf:ro \
    httpd:2.4 /entrypoint.sh
