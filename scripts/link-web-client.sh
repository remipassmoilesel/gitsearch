#!/usr/bin/env bash

# This script link web-client/dist to gitsearch data folder, for development purposes

mkdir -p ~/.gitsearch
rm ~/.gitsearch/web-client
ln -s $(pwd)/web-client/dist ~/.gitsearch/web-client