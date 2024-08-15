#!/bin/bash

set -a
source $ENV_DIR/.env.ddns
set +a

go run main.go
