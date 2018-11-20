#!/bin/bash

set -e

docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
docker images

docker tag "statshub_api" "joesweeny/statshub_api:$RELEASE"
docker push "joesweeny/statshub_api:$RELEASE"

docker tag "statshub_console" "joesweeny/statshub_console:$RELEASE"
docker push "joesweeny/statshub_console:$RELEASE"