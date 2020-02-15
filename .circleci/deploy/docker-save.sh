#!/bin/bash

set -e

mkdir -p /tmp/workspace/docker-cache

docker save -o /tmp/workspace/docker-cache/statisticodata_migrate.tar statisticodata_migrate:latest
docker save -o /tmp/workspace/docker-cache/statisticodata_rest.tar statisticodata_rest:latest
