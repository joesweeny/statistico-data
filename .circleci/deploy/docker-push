#!/bin/bash

set -e

docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
docker tag "${DOCKER_IMAGE}:latest" "${DOCKER_HUB_REPOSITORY}/${DOCKER_IMAGE}:latest"
docker push "${DOCKER_HUB_REPOSITORY}/${DOCKER_IMAGE}:latest"