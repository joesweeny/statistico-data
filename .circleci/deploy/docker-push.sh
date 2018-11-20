#!/bin/bash

set -e

docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
docker images

docker tag "statshub_api" "joesweeny/statshub_api:$CIRCLE_SHA1"
docker push "joesweeny/statshub_api:$CIRCLE_SHA1"

docker tag "statshub_console" "joesweeny/statshub_console:$CIRCLE_SHA1"
docker push "joesweeny/statshub_console:$CIRCLE_SHA1"