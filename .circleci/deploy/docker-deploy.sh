#!/bin/bash

set -e

RELEASE=$(git rev-parse HEAD)
echo "Release number: $RELEASE"
scp docker-compose.production.yml root@${PRODUCTION_SERVER}:~
ssh -t root@${PRODUCTION_SERVER} 'export RELEASE='"'$RELEASE'"'; bash -s && docker-compose -f ~/docker-compose.production.yml pull web && \
docker-compose -f ~/docker-compose.production.yml up -d'
