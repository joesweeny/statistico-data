#!/bin/bash

set -e

docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
docker images
docker tag "stats-hub_web" "joesweeny/stats-hub_web:$(git rev-parse HEAD)"
docker push "joesweeny/stats-hub_web:$(git rev-parse HEAD)"