package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/handler"
	"github.com/statistico/statistico-data/internal/app/proto"
	"time"
)

type FixtureService struct {
	FixtureRepo app.FixtureRepository
	Handler handler.FixtureHandler
	Logger *logrus.Logger
}

func (s *FixtureService) ListFixtures(r *proto.DateRangeRequest, stream proto.FixtureService_ListFixturesServer) error {
	from, err := time.Parse(time.RFC3339, r.DateFrom)

	if err != nil {
		return ErrTimeParse
	}

	to, err := time.Parse(time.RFC3339, r.DateTo)

	if err != nil {
		return ErrTimeParse
	}

	fixtures, err := s.FixtureRepo.Between(from, to)

	if err != nil {
		s.Logger.Printf("Error retrieving Fixture(s). Error: %s", err.Error())
		m := fmt.Sprint("Server Error: Unable to fulfil Request")
		return errors.New(m)
	}

	for _, fix := range fixtures {
		f, err := s.Handler.HandleFixture(&fix)

		if err != nil {
			s.Logger.Printf("Error hydrating Fixture. Error: %s", err.Error())
			continue
		}

		if err := stream.Send(f); err != nil {
			s.Logger.Printf("Error streaming Fixture back to client. Error: %s", err.Error())
			continue
		}
	}

	return nil
}

func (s *FixtureService) FixtureByID(c context.Context, r *proto.FixtureRequest) (*proto.Fixture, error) {
	fix, err := s.FixtureRepo.ByID(uint64(r.FixtureId))

	if err != nil {
		m := fmt.Sprintf("Fixture with ID %d does not exist", r.FixtureId)
		return nil, errors.New(m)
	}

	f, err := s.Handler.HandleFixture(fix)

	if err != nil {
		s.Logger.Printf("Error hydrating Fixture. Error: %s", err.Error())
		return nil, err
	}

	return f, nil
}
