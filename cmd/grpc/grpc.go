package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	"github.com/joesweeny/statistico-data/internal/container"
	"github.com/joesweeny/statistico-data/internal/config"
	fix "github.com/joesweeny/statistico-data/proto/fixture"
	res "github.com/joesweeny/statistico-data/proto/result"
	"google.golang.org/grpc/reflection"
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

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
