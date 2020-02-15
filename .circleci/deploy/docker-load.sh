#!/bin/bash

set -e

docker load -i /tmp/workspace/docker-cache/statisticodata_grpc.tar
docker load -i /tmp/workspace/docker-cache/statisticodata_rest.tar