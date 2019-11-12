package grpc

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/handler"
	"github.com/statistico/statistico-data/internal/app/proto"
	"time"
)

const maxLimit = 10000

var ErrTimeParse = errors.New("unable to parse date provided in Request")

type ResultService struct {
	FixtureRepo app.FixtureRepository
	ResultRepo  app.ResultRepository
	Handler handler.ResultHandler
	Logger *logrus.Logger
}

func (s ResultService) GetHistoricalResultsForFixture(r *proto.HistoricalResultRequest, stream proto.ResultService_GetHistoricalResultsForFixtureServer) error {
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

func (s ResultService) GetResultsForTeam(r *proto.TeamRequest, stream proto.ResultService_GetResultsForTeamServer) error {
	date, err := time.Parse(time.RFC3339, r.DateBefore)

	if err != nil {
		return ErrTimeParse
	}

	limit := r.Limit.GetValue()

	if limit == 0 {
		limit = maxLimit
	}

	fixtures, err := s.FixtureRepo.ByTeamID(uint64(r.TeamId), limit, date)

	if err != nil {
		s.Logger.Printf("Error retrieving Fixture(s) in Result Service. Error: %s", err.Error())
		return fmt.Errorf("server error: Unable to fulfil Request")
	}

	return s.sendResults(fixtures, stream)
}

func (s ResultService) GetResultsForSeason(r *proto.SeasonRequest, stream proto.ResultService_GetResultsForSeasonServer) error {
	date, err := time.Parse(time.RFC3339, r.DateBefore)

	if err != nil {
		return ErrTimeParse
	}

	fixtures, err := s.FixtureRepo.BySeasonIDBefore(uint64(r.SeasonId), date)

	if err != nil {
		s.Logger.Printf("Error retrieving Fixture(s) in Result Service. Error: %s", err.Error())
		return fmt.Errorf("server error: Unable to fulfil Request")
	}

	return s.sendResults(fixtures, stream)
}

func (s ResultService) sendResults(f []app.Fixture, stream proto.ResultService_GetResultsForTeamServer) error {
	for _, fix := range f {
		res, err := s.ResultRepo.ByFixtureID(uint64(fix.ID))

		if err != nil {
			return fmt.Errorf("fixture with ID %d does not exist", fix.ID)
		}

		x, err := s.Handler.HandleResult(&fix, res)

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
