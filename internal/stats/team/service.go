package team_stats

import (
	"github.com/statistico/statistico-data/internal/fixture"
	"log"
	pb "github.com/statistico/statistico-data/internal/proto/stats/team"
	"context"
	"errors"
	"fmt"
	"github.com/statistico/statistico-data/internal/proto"
)

type Service struct {
	TeamRepository TeamRepository
	FixtureRepository fixture.Repository
	Logger *log.Logger
}

func (s Service) GetTeamStatsForFixture(c context.Context, r *pb.FixtureRequest) (*pb.StatsResponse, error) {
	fix, err := s.FixtureRepository.ById(r.FixtureId)

	if err != nil {
		m := fmt.Sprintf("Fixture with ID %d does not exist", r.FixtureId)
		return nil, errors.New(m)
	}

	res := pb.StatsResponse{}

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
