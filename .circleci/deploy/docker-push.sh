#!/bin/bash

set -e

aws ecr get-login --no-include-email --region $AWS_DEFAULT_REGION | bash

docker tag "statisticodata_console" "216629550457.dkr.ecr.eu-west-2.amazonaws.com/statistico-data/console:$CIRCLE_SHA1"
docker push "216629550457.dkr.ecr.eu-west-2.amazonaws.com/statistico-data/console:$CIRCLE_SHA1"

docker tag "statisticodata_cron" "216629550457.dkr.ecr.eu-west-2.amazonaws.com/statistico-data/cron:$CIRCLE_SHA1"
docker push "216629550457.dkr.ecr.eu-west-2.amazonaws.com/statistico-data/cron:$CIRCLE_SHA1"

docker tag "statisticodata_grpc" "216629550457.dkr.ecr.eu-west-2.amazonaws.com/statistico-data/grpc:$CIRCLE_SHA1"
docker push "216629550457.dkr.ecr.eu-west-2.amazonaws.com/statistico-data/grpc:$CIRCLE_SHA1"

docker tag "statisticodata_rest" "216629550457.dkr.ecr.eu-west-2.amazonaws.com/statistico-data/rest:$CIRCLE_SHA1"
docker push "216629550457.dkr.ecr.eu-west-2.amazonaws.com/statistico-data/rest:$CIRCLE_SHA1"
