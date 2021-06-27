package main

import (
	"github.com/statistico/statistico-football-data/internal/bootstrap"
	"github.com/statistico/statistico-proto/go"
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

	statistico.RegisterCompetitionServiceServer(server, app.CompetitionService())
	statistico.RegisterEventServiceServer(server, app.EventService())
	statistico.RegisterFixtureServiceServer(server, app.FixtureService())
	statistico.RegisterResultServiceServer(server, app.ResultService())
	statistico.RegisterPlayerStatsServiceServer(server, app.PlayerStatsService())
	statistico.RegisterSeasonServiceServer(server, app.SeasonService())
	statistico.RegisterTeamServiceServer(server, app.TeamService())
	statistico.RegisterTeamStatsServiceServer(server, app.TeamStatsService())

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
