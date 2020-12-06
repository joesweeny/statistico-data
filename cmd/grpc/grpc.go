package main

import (
	"github.com/statistico/statistico-proto/data/go"
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

	statisticoproto.RegisterCompetitionServiceServer(server, app.CompetitionService())
	statisticoproto.RegisterEventServiceServer(server, app.EventService())
	statisticoproto.RegisterFixtureServiceServer(server, app.FixtureService())
	statisticoproto.RegisterPerformanceServiceServer(server, app.PerformanceService())
	statisticoproto.RegisterResultServiceServer(server, app.ResultService())
	statisticoproto.RegisterPlayerStatsServiceServer(server, app.PlayerStatsService())
	statisticoproto.RegisterSeasonServiceServer(server, app.SeasonService())
	statisticoproto.RegisterTeamServiceServer(server, app.TeamService())
	statisticoproto.RegisterTeamStatsServiceServer(server, app.TeamStatsService())

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
