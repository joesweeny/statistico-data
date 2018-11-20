#!/bin/bash

set -e

echo "Release number: $RELEASE"
scp docker-compose.production.yml root@${PRODUCTION_SERVER}:~
ssh -t root@${PRODUCTION_SERVER} 'docker-compose -f ./docker-compose.production.yml \
pull api && docker-compose -f ./docker-compose.production.yml && \
docker-compose -f ./docker-compose.production.yml up -d'
