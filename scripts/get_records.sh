#!/bin/bash

# load env
set -a
source $ENV_DIR/.env.ddns
set +a

# list records
curl -XGET "https://api.cloudflare.com/client/v4/zones/${ZONE_ID}/dns_records" \
  --header 'Content-Type: application/json' \
  --header "Authorization: Bearer $TOKEN"
