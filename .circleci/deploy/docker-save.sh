#!/bin/bash

set -e

mkdir -p docker-cache

docker save -o docker-cache/statisticodata_rest.tar statisticodata_rest:latest