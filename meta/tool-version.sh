#!/usr/bin/env bash
set -euo pipefail

grep "$1" <.tool-versions | sed -E "s/^$1 (.*)$/\1/g"
