#!/bin/bash

set -e

docker images
docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
docker tag "${DOCKER_IMAGE}:latest" "${DOCKER_HUB_REPOSITORY}/${DOCKER_IMAGE}:$(git rev-parse HEAD)"
docker push "${DOCKER_HUB_REPOSITORY}/${DOCKER_IMAGE}:$(git rev-parse HEAD)"