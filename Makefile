build:
	protoc -I. --go_out=plugins=grpc:$(GOPATH)/src/github.com/joesweeny/statshub \
		proto/fixture/fixture.proto