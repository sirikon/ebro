#!/usr/bin/env bash
set -euo pipefail

(
    mkdir -p /var/www/git/test.git
    cd /var/www/git/test.git
    cp /content/* .
    git init
    git add .
    git commit -m "first commit"
)
exec httpd-foreground
