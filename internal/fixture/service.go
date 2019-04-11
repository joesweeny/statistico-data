package fixture

import (
	"errors"
	pb "github.com/statistico/statistico-data/internal/proto/fixture"
	"golang.org/x/net/context"
	"log"
	"time"
	"fmt"
)

var ErrTimeParse = errors.New("unable to parse date provided in Request")

type Service struct {
	Repository
	Handler
	Logger *log.Logger
}

func (s *Service) ListFixtures(r *pb.DateRangeRequest, stream pb.FixtureService_ListFixturesServer) error {
	from, err := time.Parse(time.RFC3339, r.DateFrom)

	if err != nil {
		return ErrTimeParse
	}

	to, err := time.Parse(time.RFC3339, r.DateTo)

	if err != nil {
		return ErrTimeParse
	}

	fixtures, err := s.Repository.Between(from, to)

	if err != nil {
		s.Logger.Printf("Error retrieving Fixture(s). Error: %s", err.Error())
		m := fmt.Sprint("Server Error: Unable to fulfil Request")
		return errors.New(m)
	}

	for _, fix := range fixtures {
		f, err := s.HandleFixture(&fix)

		if err != nil {
			s.Logger.Printf("Error hydrating Fixture. Error: %s", err.Error())
			continue
		}

		if err := stream.Send(f); err != nil {
			s.Logger.Printf("Error hydrating streaming Fixture back to client. Error: %s", err.Error())
			continue
		}
	}

	return nil
}

func (s *Service) FixtureByID(c context.Context, r *pb.FixtureRequest) (*pb.Fixture, error) {
	fix, err := s.Repository.ById(int(r.FixtureId))

	if err != nil {
		m := fmt.Sprintf("Fixture with ID %d does not exist", r.FixtureId)
		return nil, errors.New(m)
	}

	f, err := s.HandleFixture(fix)

	if err != nil {
		s.Logger.Printf("Error hydrating Fixture. Error: %s", err.Error())
		return nil, err
	}

	return f, nil
}
