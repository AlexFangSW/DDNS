#!/bin/bash

set -a
source .env
set +a

docker run --rm --name ddns --env-file .env \
  "$DOCKER_REPO/$DOCKER_PROJECT/$DOCKER_IMAGE:$DOCKER_TAG"
