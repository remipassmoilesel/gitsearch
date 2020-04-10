#!/usr/bin/env bash

set -e

echo "Format, build, tests"
go fmt
go build
go test

echo "Web client build"
cd web-client
yarn install
yarn build

echo "Packaging"
cd ..
tar -cvf web-client.tar web-client/dist
# TODO: package app with web client