# Statistico Data

This application is the data warehouse that receives and stores data that powers other Statistico applications

### Getting started
To develop this application locally you will need to following tools and language installed locally:
- Docker
- Docker Compose
- Golang >=1.11

### Dependency Management
[Dep](https://golang.github.io/dep/) is used to handle this applications dependencies. You can install dep locally by executing:

`go get -u github.com/golang/dep/cmd/dep`

Once install you can compile this application's dependencies by executing:

`dep ensure`

### Deployment
This application is auto-deployed through [CircleCI](https://circleci.com/). Deployment scripts can be found in the `.circleci`
directory. To deploy this application to Production simply create a pull request in GitHub to `Master` branch and CircleCI
will do the rest.

### Testing
We have a dedicated docker-compose service to run tests locally which volume mounts our application code into the test container.
To improve testing visibility this application uses [gotestsum](https://github.com/gotestyourself/gotestsum) when executing
tests which can be installed locally by executing:

`go get gotest.tools/gotestsum`

To run the full test suite a handy script is located in the `/bin` directory, to execute:

`bin/docker-dev-test`

To narrow tests down to individual packages additional flags can be appended:

`bin/docker-dev-test ./internal/fixture`

To narrow down tests further to individual test cases additional flags can be appended:

`bin/docker-dev-test ./internalfixture -run=TestById`
