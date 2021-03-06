#!/usr/bin/env bash

set -e

PREFIX=" 👣 👣 👣 "

echo ""
echo "${PREFIX} Web client build"
echo ""

cd web_client
yarn clean
yarn install
yarn build
#yarn test # TODO: test !
cd ..

echo ""
echo "${PREFIX} Packaging"
echo ""

rm web_client/pkged.go || true
pkger -include /web_client/dist -o web_client/
pkger list -include /web_client/dist

./_test.sh

export GOARCH=amd64

export GOOS=linux
go build -o gitsearch-linux

export GOOS=darwin
go build -o gitsearch-macos

export GOOS=windows
go build -o gitsearch-windows
