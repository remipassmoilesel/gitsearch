#!/usr/bin/env bash

set -e

PREFIX=" 👣 👣 👣 "

echo ""
echo "${PREFIX} Go format, test, build"
echo ""

go fmt ./...
go test ./... -race
go build
