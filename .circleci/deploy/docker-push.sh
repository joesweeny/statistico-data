#!/bin/bash

set -e

docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
docker images

docker tag "statisticodata_console" "joesweeny/statisticodata_console:$CIRCLE_SHA1"
docker push "joesweeny/statisticodata_console:$CIRCLE_SHA1"

docker tag "statisticodata_migrate" "joesweeny/statisticodata_migrate:$CIRCLE_SHA1"
docker push "joesweeny/statisticodata_migrate:$CIRCLE_SHA1"

docker tag "statisticodata_cron" "joesweeny/statisticodata_cron:$CIRCLE_SHA1"
docker push "joesweeny/statisticodata_cron:$CIRCLE_SHA1"

docker tag "statisticodata_grpc" "joesweeny/statisticodata_grpc:$CIRCLE_SHA1"
docker push "joesweeny/statisticodata_grpc:$CIRCLE_SHA1"

docker tag "statisticodata_rest" "joesweeny/statisticodata_rest:$CIRCLE_SHA1"
docker push "joesweeny/statisticodata_rest:$CIRCLE_SHA1"