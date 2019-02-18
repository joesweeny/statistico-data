package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	"github.com/joesweeny/statshub/internal/container"
	"github.com/joesweeny/statshub/internal/config"
	pb "github.com/joesweeny/statshub/proto/fixture"
)

func main() {
	lis, err := net.Listen("tcp", "grpc:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	app := container.Bootstrap(config.GetConfig())

	server := grpc.NewServer()

	pb.RegisterFixtureServiceServer(server, app.FixtureService())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
