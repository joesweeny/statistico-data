package main

import (
	"github.com/statistico/statistico-data/internal/config"
	"github.com/statistico/statistico-data/internal/container"
	fix "github.com/statistico/statistico-data/internal/proto/fixture"
	st "github.com/statistico/statistico-data/internal/proto/stats/player"
	res "github.com/statistico/statistico-data/internal/proto/result"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", "grpc:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	app := container.Bootstrap(config.GetConfig())

	server := grpc.NewServer()

	fix.RegisterFixtureServiceServer(server, app.FixtureService())
	res.RegisterResultServiceServer(server, app.ResultService())
	st.RegisterPlayerStatsServiceServer(server, app.PlayerStatsService())

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
