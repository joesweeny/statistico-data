package main

import (
	"github.com/statistico/statistico-data/internal/app/proto"
	"github.com/statistico/statistico-data/internal/config"
	"github.com/statistico/statistico-data/internal/container"
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

	proto.RegisterFixtureServiceServer(server, app.FixtureService())
	proto.RegisterResultServiceServer(server, app.ResultService())
	proto.RegisterPlayerStatsServiceServer(server, app.PlayerStatsService())
	proto.RegisterTeamStatsServiceServer(server, app.TeamStatsService())

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
