#!/usr/bin/env bash

set -e

PREFIX=" 👣 👣 👣 "

echo ""
echo "${PREFIX} Web client build"
echo ""

cd web_client
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

echo ""
echo "${PREFIX} Go format, test, build"
echo ""

export GS_TEST_DATA="$(pwd)/test/data"
go fmt ./...
go test ./... -race
go build
