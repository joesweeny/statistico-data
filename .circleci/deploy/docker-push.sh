#!/bin/bash

set -e

docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
docker tag "statshub_web" "joesweeny/statshub_web:$(git rev-parse HEAD)"
docker push "joesweeny/statshub_web:$(git rev-parse HEAD)"