package player_stats

import (
	"github.com/statistico/statistico-data/internal/fixture"
	pb "github.com/statistico/statistico-data/internal/proto/stats"
	"log"
	"context"
	"errors"
	"fmt"
)

type Service struct {
	PlayerRepository PlayerRepository
	FixtureRepo      fixture.Repository
	Logger 			 *log.Logger
}

func (s Service) GetPlayerStatsForFixture(c context.Context, r *pb.FixtureRequest) (*pb.StatsResponse, error) {
	fix, err := s.FixtureRepo.ById(r.FixtureId)

	if err != nil {
		m := fmt.Sprintf("Fixture with ID %d does not exist", r.FixtureId)
		return nil, errors.New(m)
	}

	res := pb.StatsResponse{}

	home, err := s.PlayerRepository.ByFixtureAndTeam(uint64(fix.ID), uint64(fix.HomeTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving player stats: FixtureID %d, Home Team ID %d", fix.ID, fix.HomeTeamID)
		s.Logger.Println(e)
		return nil, e
	}

	res.HomeTeam = HandlePlayerStats(home)

	away, err := s.PlayerRepository.ByFixtureAndTeam(uint64(fix.ID), uint64(fix.AwayTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving player stats: FixtureID %d, Away Team ID %d", fix.ID, fix.HomeTeamID)
		s.Logger.Println(e)
		return nil, e
	}

	res.AwayTeam = HandlePlayerStats(away)

	return &res, nil
}

func (s Service) GetLineUpForFixture(c context.Context, r *pb.FixtureRequest) (*pb.LineupResponse, error) {
	fix, err := s.FixtureRepo.ById(r.FixtureId)

	if err != nil {
		m := fmt.Sprintf("Fixture with ID %d does not exist", r.FixtureId)
		return nil, errors.New(m)
	}

	res := pb.LineupResponse{}

	home, err := s.PlayerRepository.ByFixtureAndTeam(uint64(fix.ID), uint64(fix.HomeTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving player stats: FixtureID %d, Home Team ID %d", fix.ID, fix.HomeTeamID)
		s.Logger.Println(e)
		return nil, e
	}

	homeLineup := pb.Lineup{
		Start: HandleStartingLineupPlayers(home),
		Bench: HandleSubstituteLineupPlayers(home),
	}

	res.HomeTeam = &homeLineup

	away, err := s.PlayerRepository.ByFixtureAndTeam(uint64(fix.ID), uint64(fix.AwayTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving player stats: FixtureID %d, Away Team ID %d", fix.ID, fix.AwayTeamID)
		s.Logger.Println(e)
		return nil, e
	}

	awayLineup := pb.Lineup{
		Start: HandleStartingLineupPlayers(away),
		Bench: HandleSubstituteLineupPlayers(away),
	}

	res.AwayTeam = &awayLineup

	return &res, nil
}
