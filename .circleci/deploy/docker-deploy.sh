#!/bin/bash

set -e

RELEASE=$CIRCLE_SHA1
API_KEY=${SPORT_MONKS_API_KEY}
echo "Release number: $RELEASE"
echo "API KEY number: $API_KEY"

scp docker-compose.production.yml root@${PRODUCTION_SERVER}:~

ssh -t root@${PRODUCTION_SERVER} 'env RELEASE='"'$RELEASE'"' SPORT_MONK_API_KEY='"'$API_KEY'"' \

docker-compose -f ./docker-compose.production.yml pull api && \

docker-compose -f ./docker-compose.production.yml pull console && \

docker-compose -f ./docker-compose.production.yml up -d'

