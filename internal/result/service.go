package result

import (
	"errors"
	"github.com/statistico/statistico-data/internal/fixture"
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

	s.Logger.Printf("Number of fixtures retrieved %d", len(fixtures))

	if err != nil {
		return err
	}

	for _, fix := range fixtures {
		s.Logger.Printf("Fixture ID: %d", fix.ID)

		res, err := s.ResultRepo.GetByFixtureId(fix.ID)

		if err != nil {
			return err
		}

		x, err := s.HandleResult(&fix, res)

		if err != nil {
			return err
		}

		if err := stream.Send(x); err != nil {
			s.Logger.Printf("Error hydrating Result. ID: %d. Error: %s", res.FixtureID, err.Error())
			return err
		}
	}

	return nil
}
