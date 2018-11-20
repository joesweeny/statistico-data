#!/bin/bash

set -e

RELEASE=$CIRCLE_SHA1
echo "Release number: $RELEASE"
scp docker-compose.production.yml root@${PRODUCTION_SERVER}:~

ssh -t root@${PRODUCTION_SERVER} \
'export RELEASE='"'$RELEASE'"'; \
'export SPORT_MONKS_API_KEY='"'${SPORT_MONKS_API_KEY}'"'; \

docker-compose -f ./docker-compose.production.yml \
pull api && docker-compose -f ./docker-compose.production.yml pull console && \
docker-compose -f ./docker-compose.production.yml up -d'
