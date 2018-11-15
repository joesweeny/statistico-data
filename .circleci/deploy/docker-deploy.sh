#!/bin/bash

set -e

scp docker-compose.production.yml root@${PRODUCTION_SERVER}:~/stats-hub
ssh root@${PRODUCTION_SERVER} "docker-compose -f ~/stats-hub/docker-compose.production.yml pull web && \
docker-compose -f ~/stats-hub/docker-compose.production.yml up -d && \
docker restart nginxproxy_nginx-proxy_1"