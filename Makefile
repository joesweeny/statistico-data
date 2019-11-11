build:
	protoc -I/usr/local/include -I. internal/app/proto/competition.proto \
    	--go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/app/proto/fixture.proto \
		--go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/app/proto/result.proto \
    	--go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/app/proto/round.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/app/proto/season.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/app/proto/player_stats.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/app/proto/team_stats.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/app/proto/team.proto \
        --go_out=plugins=grpc:$(GOPATH)/src

	protoc -I/usr/local/include -I. internal/app/proto/venue.proto \
        --go_out=plugins=grpc:$(GOPATH)/src