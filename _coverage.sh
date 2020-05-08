#!/usr/bin/env bash

#
# This is script must be run like this:
#
#   $ cd gitsearch
#   $ ./test/_coverage.sh
#

set -e

PREFIX=" 👣 👣 👣 "

./_test.sh

echo ""
echo "${PREFIX} Go cover"
echo ""

go tool cover -html=coverage.out