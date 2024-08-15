#!/bin/bash

set -a
source $ENV_DIR/.env.ddns
set +a

docker push "$DOCKER_REPO/$DOCKER_PROJECT/$DOCKER_IMAGE:$DOCKER_TAG"
