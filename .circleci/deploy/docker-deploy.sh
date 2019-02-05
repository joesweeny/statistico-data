#!/bin/bash

set -e

RELEASE=$CIRCLE_SHA1
SPORTMONKS_API_KEY=${SPORTMONKS_API_KEY}
DB_DRIVER=${DB_DRIVER}
DB_HOST=${DB_HOST}
DB_NAME=${DB_NAME}
DB_PASSWORD=${DB_PASSWORD}
DB_PORT=${DB_PORT}
DB_USER=${DB_USER}

scp docker-compose.production.yml root@${PRODUCTION_SERVER}:~

ssh -t root@${PRODUCTION_SERVER} 'export RELEASE='"'$RELEASE'"' export SPORTMONKS_API_KEY='"'$SPORTMONKS_API_KEY'"' export DB_DRIVER='"'$DB_DRIVER'"' export DB_HOST='"'$DB_HOST'"' export DB_NAME='"'$DB_NAME'"' export DB_NAME='"'$DB_NAME'"' export DB_PASSWORD='"'$DB_PASSWORD'"' export DB_PORT='"'$DB_PORT'"' export DB_USER='"'$DB_USER'"' \

docker-compose -f ./docker-compose.production.yml pull api && \

docker-compose -f ./docker-compose.production.yml pull console && \

docker-compose -f ./docker-compose.production.yml pull migrate && \

docker-compose -f ./docker-compose.production.yml pull cron && \

docker-compose -f ./docker-compose.production.yml up -d'

