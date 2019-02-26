build:
	protoc -I/usr/local/include -I. proto/fixture/fixture.proto \
		--go_out=plugins=grpc:$(GOPATH)/src/github.com/joesweeny/statistico-data
