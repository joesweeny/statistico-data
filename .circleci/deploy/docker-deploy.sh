#!/bin/bash

set -e

RELEASE=$CIRCLE_SHA1
API_KEY=$SPORT_MONKS_API_KEY
echo "Release number: $RELEASE"
echo "API KEY number: $APY_KEY"

scp docker-compose.production.yml root@${PRODUCTION_SERVER}:~

ssh -t root@${PRODUCTION_SERVER} 'export RELEASE='"'$RELEASE'"'; \

docker-compose -f ./docker-compose.production.yml \

pull api && docker-compose -f ./docker-compose.production.yml && \

docker-compose -f ./docker-compose.production.yml up -d'

