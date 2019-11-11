package team_stats

import (
	"context"
	"errors"
	"fmt"
	"github.com/statistico/statistico-data/internal/app"
	proto2 "github.com/statistico/statistico-data/internal/app/proto"
	"github.com/statistico/statistico-data/internal/proto"
	"log"
)

type Service struct {
	TeamRepository    app.TeamStatsRepository
	FixtureRepository app.FixtureRepository
	Logger            *log.Logger
}

func (s Service) GetTeamStatsForFixture(c context.Context, r *proto2.FixtureRequest) (*proto2.StatsResponse, error) {
	fix, err := s.FixtureRepository.ByID(r.FixtureId)

	if err != nil {
		m := fmt.Sprintf("Fixture with ID %d does not exist", r.FixtureId)
		return nil, errors.New(m)
	}

	res := proto2.StatsResponse{}

	home, err := s.TeamRepository.ByFixtureAndTeam(uint64(fix.ID), uint64(fix.HomeTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving team stats: FixtureID %d, Home Team ID %d", fix.ID, fix.HomeTeamID)
		s.Logger.Println(e)
		return nil, e
	}

	res.HomeTeam = proto.TeamStatsToProto(home)

	away, err := s.TeamRepository.ByFixtureAndTeam(uint64(fix.ID), uint64(fix.AwayTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving team stats: FixtureID %d, Away Team ID %d", fix.ID, fix.HomeTeamID)
		s.Logger.Println(e)
		return nil, e
	}

	res.AwayTeam = proto.TeamStatsToProto(away)

	return &res, nil
}
