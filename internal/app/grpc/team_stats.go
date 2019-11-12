package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/converter"
	"github.com/statistico/statistico-data/internal/app/proto"
)

type TeamStatsService struct {
	TeamRepository    app.TeamStatsRepository
	FixtureRepository app.FixtureRepository
	Logger *logrus.Logger
}

func (s TeamStatsService) GetTeamStatsForFixture(c context.Context, r *proto.FixtureRequest) (*proto.TeamStatsResponse, error) {
	fix, err := s.FixtureRepository.ByID(uint64(r.FixtureId))

	if err != nil {
		m := fmt.Sprintf("Fixture with ID %d does not exist", r.FixtureId)
		return nil, errors.New(m)
	}

	res := proto.TeamStatsResponse{}

	home, err := s.TeamRepository.ByFixtureAndTeam(uint64(fix.ID), uint64(fix.HomeTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving team stats: FixtureID %d, Home Team ID %d", fix.ID, fix.HomeTeamID)
		s.Logger.Println(e)
		return nil, e
	}

	res.HomeTeam = converter.TeamStatsToProto(home)

	away, err := s.TeamRepository.ByFixtureAndTeam(uint64(fix.ID), uint64(fix.AwayTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving team stats: FixtureID %d, Away Team ID %d", fix.ID, fix.HomeTeamID)
		s.Logger.Println(e)
		return nil, e
	}

	res.AwayTeam = converter.TeamStatsToProto(away)

	return &res, nil
}
