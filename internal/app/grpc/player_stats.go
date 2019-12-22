package grpc

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/factory"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlayerStatsService struct {
	fixtureRepo app.FixtureRepository
	factory     *factory.PlayerStatsFactory
	logger      *logrus.Logger
}

func (s PlayerStatsService) GetPlayerStatsForFixture(c context.Context, r *proto.FixtureRequest) (*proto.PlayerStatsResponse, error) {
	fix, err := s.fixtureRepo.ByID(r.FixtureId)

	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("fixture with ID %d does not exist", r.FixtureId))
	}

	home, err := s.factory.BuildPlayerStats(fix, fix.HomeTeamID)

	if err != nil {
		s.logger.Warnf("Error hydrating proto player stats: %s", err.Error())
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("home player stats do not exist for fixture %d", r.FixtureId),
		)
	}

	away, err := s.factory.BuildPlayerStats(fix, fix.AwayTeamID)

	if err != nil {
		s.logger.Warnf("Error hydrating proto player stats: %s", err.Error())
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("away player stats do not exist for fixture %d", r.FixtureId),
		)
	}

	res := proto.PlayerStatsResponse{
		HomeTeam: home,
		AwayTeam: away,
	}

	return &res, nil
}

func (s PlayerStatsService) GetLineUpForFixture(c context.Context, r *proto.FixtureRequest) (*proto.LineupResponse, error) {
	fix, err := s.fixtureRepo.ByID(r.FixtureId)

	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("fixture with ID %d does not exist", r.FixtureId))
	}

	home, err := s.factory.BuildLineup(fix, fix.HomeTeamID)

	if err != nil {
		s.logger.Warnf("Error hydrating proto lineup: %s", err.Error())
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("home lineup do not exist for fixture %d", r.FixtureId),
		)
	}

	away, err := s.factory.BuildLineup(fix, fix.AwayTeamID)

	if err != nil {
		s.logger.Warnf("Error hydrating proto lineup: %s", err.Error())
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("away lineup do not exist for fixture %d", r.FixtureId),
		)
	}

	res := proto.LineupResponse{
		HomeTeam: home,
		AwayTeam: away,
	}

	return &res, nil
}

func NewPlayerStatsService(r app.FixtureRepository, f *factory.PlayerStatsFactory, log *logrus.Logger) *PlayerStatsService {
	return &PlayerStatsService{fixtureRepo: r, factory: f, logger: log}
}
