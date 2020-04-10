#!/usr/bin/env bash

set -e

echo ""
echo "Go format, build, tests"
echo ""

go fmt
go build
go test

echo ""
echo "Web client build"
echo ""

cd web-client
yarn install
yarn build

echo ""
echo "Packaging"
echo ""

cd ..
tar -cvf web-client.tar web-client/dist
# TODO: package app with web client
