#!/bin/bash

set -e

docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
docker images
docker tag "statshub_api" "joesweeny/statshub_api:$(git rev-parse HEAD)"
docker push "joesweeny/statshub_api:$(git rev-parse HEAD)"