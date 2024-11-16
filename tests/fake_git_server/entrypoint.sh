#!/usr/bin/env sh
set -e
git config --global init.defaultBranch master
git init --bare /git/fake.git
mkdir -p /temp_repo
(
    cd /temp_repo
    cp -r /content/* .
    git init
    git config user.email "you@example.com"
    git config user.name "example"
    git add .
    git commit -m "first commit"
    git remote add local /git/fake.git
    git push local
)

spawn-fcgi -s /run/fcgi.sock /usr/bin/fcgiwrap
exec nginx -g "daemon off;"
