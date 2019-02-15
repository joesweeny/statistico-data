package main

import (
	pb "github.com/joesweeny/statshub/proto/fixture"
	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()

	pb.RegisterFixtureServiceServer(s)
}