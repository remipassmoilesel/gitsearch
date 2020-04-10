#!/usr/bin/env bash

IMAGE_NAME="registry.gitlab.com/remipassmoilesel/gitsearch/ci-image:0.2"

docker run -v $(pwd):/build:rw -ti $IMAGE_NAME /bin/bash
