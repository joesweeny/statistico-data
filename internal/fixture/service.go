package fixture

import (
	pb "github.com/joesweeny/statshub/proto/fixture"
	"time"
	"github.com/pkg/errors"
	"log"
)

const dateFormat = "2006-01-02"
var ErrTimeParse = errors.New("unable to parse date provided in Request")

type Service struct {
	Repository
	Handler
	Logger *log.Logger
}

func (s *Service) ListFixtures(r *pb.Request, stream pb.FixtureService_ListFixturesServer) error {
	from, err := time.Parse(dateFormat, r.DateFrom)
	y, m, d := from.Date()
	from = time.Date(y, m, d, 0, 0, 0, 0, from.Location())

	if err != nil {
		return ErrTimeParse
	}

	to, err := time.Parse(dateFormat, r.DateTo)
	y, m, d = to.Date()
	to = time.Date(y, m, d, 23, 59, 59, 59, to.Location())

	if err != nil {
		return ErrTimeParse
	}

	fixtures, err := s.Repository.Between(from, to)

	if err != nil {
		return err
	}

	for _, fix := range fixtures {
		f, err := s.HandleFixture(&fix)

		if err != nil {
			return err
		}

		if err := stream.Send(f); err != nil {
			s.Logger.Printf("Error hydrating Fixture. ID: %d. Error: %s", fix.ID, err.Error())
			return err
		}
	}

	return nil
}
