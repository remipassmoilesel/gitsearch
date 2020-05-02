#!/usr/bin/env bash

#
# This is script must be run like this:
#
#   $ cd gitsearch
#   $ ./test/test.sh
#

set -e

PREFIX=" 👣 👣 👣 "

echo ""
echo "${PREFIX} Go format, test, build"
echo ""

go fmt ./...
go test ./... -race -cover -coverprofile=coverage.out
go build
