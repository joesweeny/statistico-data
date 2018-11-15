#!/bin/bash

set -e

sed -i s/\'$RELEASE\'/$(git rev-parse HEAD)/ docker-compose.production.yml
scp docker-compose.production.yml root@${PRODUCTION_SERVER}:~/stats-hub
ssh root@${PRODUCTION_SERVER} "docker-compose -f ~/stats-hub/docker-compose.production.yml pull web && \
docker-compose -f ~/stats-hub/docker-compose.production.yml up -d"
