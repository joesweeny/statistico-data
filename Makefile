build:
	protoc -I/usr/local/include -I. internal/proto/competition/competition.proto \
    	--go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/proto/fixture/fixture.proto \
		--go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/proto/result/result.proto \
    	--go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/proto/round/round.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/proto/season/season.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/proto/stats/player/stats.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/proto/stats/team/stats.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/proto/team/team.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/proto/venue/venue.proto \
        --go_out=plugins=grpc:$(GOPATH)/src