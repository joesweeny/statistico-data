package result

import (
	"errors"
	"fmt"
	"github.com/statistico/statistico-data/internal/fixture"
	"github.com/statistico/statistico-data/internal/model"
	pb "github.com/statistico/statistico-data/internal/proto/result"
	"log"
	"time"
)

const maxLimit = 10000

var ErrTimeParse = errors.New("unable to parse date provided in Request")

type Service struct {
	FixtureRepo fixture.Repository
	ResultRepo  Repository
	Handler
	Logger *log.Logger
}

func (s Service) GetHistoricalResultsForFixture(r *pb.HistoricalResultRequest, stream pb.ResultService_GetHistoricalResultsForFixtureServer) error {
	date, err := time.Parse(time.RFC3339, r.DateBefore)

	if err != nil {
		return ErrTimeParse
	}

	fixtures, err := s.FixtureRepo.ByHomeAndAwayTeam(r.HomeTeamId, r.AwayTeamId, r.Limit, date)

	if err != nil {
		s.Logger.Printf("Error retrieving Fixture(s) in Result Service. Error: %s", err.Error())
		return fmt.Errorf("server error: Unable to fulfil Request")
	}

	return s.sendResults(fixtures, stream)
}

func (s Service) GetResultsForTeam(r *pb.TeamRequest, stream pb.ResultService_GetResultsForTeamServer) error {
	date, err := time.Parse(time.RFC3339, r.DateBefore)

	if err != nil {
		return ErrTimeParse
	}

	limit := r.Limit.GetValue()

	if limit == 0 {
		limit = maxLimit
	}

	fixtures, err := s.FixtureRepo.ByTeamId(r.TeamId, limit, date)

	if err != nil {
		s.Logger.Printf("Error retrieving Fixture(s) in Result Service. Error: %s", err.Error())
		return fmt.Errorf("server error: Unable to fulfil Request")
	}

	return s.sendResults(fixtures, stream)
}

func (s Service) GetResultsForSeason(r *pb.SeasonRequest, stream pb.ResultService_GetResultsForSeasonServer) error {
	date, err := time.Parse(time.RFC3339, r.DateBefore)

	if err != nil {
		return ErrTimeParse
	}

	fixtures, err := s.FixtureRepo.BySeasonId(r.SeasonId, date)

	if err != nil {
		s.Logger.Printf("Error retrieving Fixture(s) in Result Service. Error: %s", err.Error())
		return fmt.Errorf("server error: Unable to fulfil Request")
	}

	return s.sendResults(fixtures, stream)
}

func (s Service) sendResults(f []model.Fixture, stream pb.ResultService_GetResultsForTeamServer) error {
	for _, fix := range f {
		res, err := s.ResultRepo.GetByFixtureId(fix.ID)

		if err != nil {
			return fmt.Errorf("fixture with ID %d does not exist", fix.ID)
		}

		x, err := s.HandleResult(&fix, res)

		if err != nil {
			s.Logger.Printf("Error hydrating Result. Error: %s", err.Error())
			return err
		}

		if err := stream.Send(x); err != nil {
			s.Logger.Printf("Error streaming Result back to client. Error: %s", err.Error())
			continue
		}
	}

	return nil
}
