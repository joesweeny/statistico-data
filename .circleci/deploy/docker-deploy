#!/bin/bash

set -e

scp docker-compose.production.yml root@${PRODUCTION_SERVER}:~/staging
ssh root@${PRODUCTION_SERVER} "docker-compose -f ./production/docker-compose.production.yml pull web && \
docker-compose -f ./production/docker-compose.production.yml up -d