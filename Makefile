build:
	protoc -I/usr/local/include -I. proto/competition/competition.proto \
    	--go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. proto/fixture/fixture.proto \
		--go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. proto/result/result.proto \
    	--go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. proto/round/round.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. proto/season/season.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. proto/team/team.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. proto/venue/venue.proto \
        --go_out=plugins=grpc:$(GOPATH)/src