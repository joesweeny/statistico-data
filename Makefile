build:
	protoc -I/usr/local/include -I. proto/fixture/fixture.proto \
		--go_out=plugins=grpc:$(GOPATH)/src/github.com/joesweeny/statistico-data

	protoc -I/usr/local/include -I. proto/result/result.proto \
    		--go_out=plugins=grpc:$(GOPATH)/src/github.com/joesweeny/statistico-data