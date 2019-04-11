# Statistico Data

This application is the data warehouse that receives and stores data that powers other Statistico applications

## Getting started
To develop this application locally you will need to following tools and language installed locally:
- Docker
- Docker Compose
- Golang >=1.11

## Dependency Management
[Dep](https://golang.github.io/dep/) is used to handle this applications dependencies. You can install dep locally by executing:

`go get -u github.com/golang/dep/cmd/dep`

Once install you can compile this application's dependencies by executing:

`dep ensure`

## Deployment
This application is auto-deployed through [CircleCI](https://circleci.com/). Deployment scripts can be found in the `.circleci`
directory. To deploy this application to Production simply create a pull request in GitHub to `Master` branch and CircleCI
will do the rest.

## Testing
We have a dedicated docker-compose service to run tests locally which volume mounts our application code into the test container.
To improve testing visibility this application uses [gotestsum](https://github.com/gotestyourself/gotestsum) when executing
tests which can be installed locally by executing:

`go get gotest.tools/gotestsum`

To run the full test suite a handy script is located in the `/bin` directory, to execute:

`bin/docker-dev-test`

To narrow tests down to an individual directory additional flags can be appended:

`bin/docker-dev-test ./internal/....`

To narrow tests down further to individual test cases additional flags can be appended:

`bin/docker-dev-test ./internal/fixture -run=TestById`

Alternatively the test suite can be run by using Golang's inbuilt testing tool and executing the following command:

`bin/docker-dev run --rm test go test -v ./internal/...`

The suite contains integration tests that depend on an external database therefore tests need to be run inside the test
container

## gRPC
Statistico's internal systems communicate via gRPC. This application's gRPC specifications can be found in the 
[/proto](https://github.com/statistico/statistico-data/proto) directory. For more on gRPC view [here](https://grpc.io/docs/guides/)

This application exposes two services:
- FixtureService
- ResultService

The parameters required to access these services are well defined in their respective `.proto` files. 

To access this applications services using a local client we recommend [gRPCurl](https://github.com/fullstorydev/grpcurl). 
Example calls are:

#### To fetch fixtures between a date period
```proto
grpcurl \
    -plaintext \
    -d \
    '{"date_from": "2019-04-03T00:00:00+00:00", "date_to": "2019-04-03T23:59:59+00:00"}' \
    localhost:50051  \
    fixture.FixtureService/ListFixtures
```
#### To fetch a fixture by ID
```proto
grpcurl \
    -plaintext \
    -d \
    '{"fixture_id": 5601}' \
    localhost:50051  \
    fixture.FixtureService/ListFixtures
```
#### To fetch results for a given Team
```proto
grpcurl \
    -plaintext \
    -d \
    '{"team_id": 501, "limit": 75, "date_before": "2019-04-03T23:59:59+00:00"}' \
    localhost:50051  \
    result.ResultService/GetResultsForTeam
```