#!/bin/bash

set -a
source .env
set +a

docker push "$DOCKER_REPO/$DOCKER_PROJECT/$DOCKER_IMAGE:$DOCKER_TAG"
