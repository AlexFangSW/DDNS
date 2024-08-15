#!/bin/bash

set -a
source $ENV_DIR/.env.ddns
set +a

docker run -d --name ddns \
  --env-file $ENV_DIR/.env.ddns \
  --network infra  \
	--restart=unless-stopped \
  "$DOCKER_REPO/$DOCKER_PROJECT/$DOCKER_IMAGE:$DOCKER_TAG"
