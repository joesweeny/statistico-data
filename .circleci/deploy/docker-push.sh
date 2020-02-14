#!/bin/bash

set -e

docker tag "statisticodata_rest" "216629550457.dkr.ecr.eu-west-2.amazonaws.com/statistico-data/rest:$CIRCLE_SHA1"
docker push "216629550457.dkr.ecr.eu-west-2.amazonaws.com/statistico-data/rest:$CIRCLE_SHA1"
