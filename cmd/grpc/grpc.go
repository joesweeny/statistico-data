package main

import (
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"github.com/statistico/statistico-data/internal/bootstrap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	app := bootstrap.BuildContainer(bootstrap.BuildConfig())

	opts := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle:5*time.Minute})
	server := grpc.NewServer(opts)

	proto.RegisterCompetitionServiceServer(server, app.CompetitionService())
	proto.RegisterFixtureServiceServer(server, app.FixtureService())
	proto.RegisterPerformanceServiceServer(server, app.PerformanceService())
	proto.RegisterResultServiceServer(server, app.ResultService())
	proto.RegisterPlayerStatsServiceServer(server, app.PlayerStatsService())
	proto.RegisterSeasonServiceServer(server, app.SeasonService())
	proto.RegisterTeamServiceServer(server, app.TeamService())
	proto.RegisterTeamStatsServiceServer(server, app.TeamStatsService())

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
