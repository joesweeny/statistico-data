#!/bin/bash

set -e

scp docker-compose.production.yml root@${PRODUCTION_SERVER}:~/production
ssh root@${PRODUCTION_SERVER} "docker-compose -f ./production/docker-compose.production.yml pull web
docker-compose -f ./production/docker-compose.production.yml up -d