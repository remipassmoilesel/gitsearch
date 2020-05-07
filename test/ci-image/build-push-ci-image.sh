#!/usr/bin/env bash

IMAGE_NAME="registry.gitlab.com/remipassmoilesel/gitsearch/ci-image:0.6"

docker build . -t $IMAGE_NAME
docker push $IMAGE_NAME
