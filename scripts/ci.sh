#!/usr/bin/env bash

set -e

PREFIX=" 👣 👣 👣 "

echo ""
echo "${PREFIX} Go format, build, tests"
echo ""

go fmt
go build
go test -race

echo ""
echo "${PREFIX} Web client build"
echo ""

cd web-client
yarn install
yarn build
#yarn test

# TODO: package app and web client
