#!/bin/bash

set -e

mkdir -p docker-cache

docker save -o /tmp/workspace/docker-cache/statisticodata_rest.tar statisticodata_rest:latest