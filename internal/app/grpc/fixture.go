package grpc

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/factory"
	"github.com/statistico/statistico-proto/data/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type FixtureService struct {
	fixtureRepo app.FixtureRepository
	factory     *factory.FixtureFactory
	logger      *logrus.Logger
}

func (s *FixtureService) ListSeasonFixtures(r *statisticoproto.SeasonFixtureRequest, stream statisticoproto.FixtureService_ListSeasonFixturesServer) error {
	from, err := time.Parse(time.RFC3339, r.DateFrom)

	if err != nil {
		return status.Error(codes.InvalidArgument, "Date provided is not a valid RFC3339 date")
	}

	to, err := time.Parse(time.RFC3339, r.DateTo)

	if err != nil {
		return status.Error(codes.InvalidArgument, "Date provided is not a valid RFC3339 date")
	}

	query := app.FixtureRepositoryQuery{
		DateTo:   &to,
		DateFrom: &from,
	}

	fixtures, err := s.fixtureRepo.Get(query)

	if err != nil {
		s.logger.Warnf("Error retrieving Fixture(s). Error: %s", err.Error())
		return status.Error(codes.Internal, "Internal server error")
	}

	for _, fix := range fixtures {
		f, err := s.factory.BuildFixture(&fix)

		if err != nil {
			s.logger.Warnf("Error hydrating Fixture. Error: %s", err.Error())
			continue
		}

		if err := stream.Send(f); err != nil {
			s.logger.Warnf("Error streaming Fixture back to client. Error: %s", err.Error())
			continue
		}
	}

	return nil
}

func (s *FixtureService) FixtureByID(c context.Context, r *statisticoproto.FixtureRequest) (*statisticoproto.Fixture, error) {
	fix, err := s.fixtureRepo.ByID(r.FixtureId)

	if err != nil {
		s.logger.Warnf("Error fetching fixture in gRPC fixture service: %s", err.Error())
		return nil, status.Error(codes.NotFound, fmt.Sprintf("fixture with ID %d does not exist", r.FixtureId))
	}

	f, err := s.factory.BuildFixture(fix)

	if err != nil {
		s.logger.Warnf("Error hydrating Fixture: %s", err.Error())
		return nil, status.Error(codes.DataLoss, "Internal server error")
	}

	return f, nil
}

func (s *FixtureService) Search(r *statisticoproto.FixtureSearchRequest, stream statisticoproto.FixtureService_SearchServer) error {
	query, err := buildFixtureRepositoryQuery(r)

	if err != nil {
		return err
	}

	fixtures, err := s.fixtureRepo.Get(query)

	if err != nil {
		s.logger.Warnf("Error retrieving Fixture(s). Error: %s", err.Error())
		return status.Error(codes.Internal, "Internal server error")
	}

	for _, fix := range fixtures {
		f, err := s.factory.BuildFixture(&fix)

		if err != nil {
			s.logger.Errorf("Error hydrating Fixture. Error: %s", err.Error())
			return status.Error(codes.Internal, "Internal server error")
		}

		if err := stream.Send(f); err != nil {
			s.logger.Errorf("Error streaming Fixture back to client. Error: %s", err.Error())
			return status.Error(codes.Internal, "Internal server error")
		}
	}

	return nil
}

func NewFixtureService(r app.FixtureRepository, f *factory.FixtureFactory, log *logrus.Logger) *FixtureService {
	return &FixtureService{fixtureRepo: r, factory: f, logger: log}
}

func buildFixtureRepositoryQuery(r *statisticoproto.FixtureSearchRequest) (app.FixtureRepositoryQuery, error) {
	var query app.FixtureRepositoryQuery

	if r.GetDateBefore() != nil {
		date, err := time.Parse(time.RFC3339, r.GetDateBefore().GetValue())

		if err != nil {
			return query, status.Error(
				codes.InvalidArgument,
				fmt.Sprintf("Date provided '%s' is not a valid RFC3339 date", r.GetDateBefore().GetValue()),
			)
		}

		query.DateTo = &date
	}

	if r.GetDateAfter() != nil {
		date, err := time.Parse(time.RFC3339, r.GetDateAfter().GetValue())

		if err != nil {
			return query, status.Error(
				codes.InvalidArgument,
				fmt.Sprintf("Date provided '%s' is not a valid RFC3339 date", r.GetDateAfter().GetValue()),
			)
		}

		query.DateFrom = &date
	}

	if r.GetLimit() != nil {
		v := r.GetLimit().GetValue()
		query.Limit = &v
	}

	if len(r.GetSeasonIds()) > 0 {
		query.SeasonIDs = r.GetSeasonIds()
	}

	if r.GetSort() != nil {
		v := r.GetSort().GetValue()
		query.SortBy = &v
	}

	return query, nil
}
