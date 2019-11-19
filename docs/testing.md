# Testing
We have a dedicated docker-compose service to run tests locally which volume mounts our application code into the test container.
To improve testing visibility this application uses [gotestsum](https://github.com/gotestyourself/gotestsum) when executing
tests which can be installed locally by executing:

`go get gotest.tools/gotestsum`

To run the full test suite a handy script is located in the `/bin` directory, to execute:

`bin/docker-dev-test`

To narrow tests down to an individual directory additional flags can be appended:

`bin/docker-dev-test ./internal/....`

To narrow tests down further to individual test cases additional flags can be appended:

`bin/docker-dev-test ./internal/app/postgres -run=TestById`

Alternatively the test suite can be run by using Golang's inbuilt testing tool and executing the following command:

`bin/docker-dev run --rm test go test -v ./internal/...`

The suite contains integration tests that depend on an external database therefore tests need to be run inside the test
container