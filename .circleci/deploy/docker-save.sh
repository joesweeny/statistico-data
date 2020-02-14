#!/bin/bash

set -e

mkdir -p /tmp/workspace

docker save -o /tmp/workspace/docker-cache/statisticodata_rest.tar statisticodata_rest:latest