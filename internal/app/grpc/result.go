package grpc

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/factory"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type ResultService struct {
	fixtureRepo app.FixtureRepository
	factory     *factory.ResultFactory
	logger      *logrus.Logger
}

func (s ResultService) GetHistoricalResultsForFixture(r *proto.HistoricalResultRequest, stream proto.ResultService_GetHistoricalResultsForFixtureServer) error {
	date, err := time.Parse(time.RFC3339, r.DateBefore)

	if err != nil {
		return status.Error(codes.InvalidArgument, "Date provided is not a valid RFC3339 date")
	}

	limit := uint64(r.Limit)

	query := app.FixtureRepositoryQuery{
		HomeTeamID: &r.HomeTeamId,
		AwayTeamID: &r.AwayTeamId,
		DateTo:     &date,
		Limit:      &limit,
	}

	fixtures, err := s.fixtureRepo.Get(query)

	if err != nil {
		s.logger.Warnf("Error retrieving Fixture(s) in Result Service. Error: %s", err.Error())
		return status.Error(codes.Internal, "Internal server error")
	}

	return s.sendResults(fixtures, stream)
}

func (s ResultService) GetResultsForTeam(r *proto.TeamResultRequest, stream proto.ResultService_GetResultsForTeamServer) error {
	var query app.FixtureFilterQuery

	if r.GetDateBefore() != nil {
		date, err := time.Parse(time.RFC3339, r.GetDateBefore().GetValue())

		if err != nil {
			return status.Error(codes.InvalidArgument, "Date provided is not a valid RFC3339 date")
		}

		query.DateBefore = &date
	}

	if r.GetLimit() != nil {
		v := r.GetLimit().GetValue()
		query.Limit = &v
	}

	if r.GetSort() != nil {
		v := r.GetSort().GetValue()
		query.SortBy = &v
	}

	if r.GetVenue() != nil {
		v := r.GetVenue().GetValue()
		query.Venue = &v
	}

	fixtures, err := s.fixtureRepo.ByTeamID(r.GetTeamId(), query)

	if err != nil {
		s.logger.Warnf("Error retrieving Fixture(s) in Result Service. Error: %s", err.Error())
		return status.Error(codes.Internal, "Internal server error")
	}

	return s.sendResults(fixtures, stream)
}

func (s ResultService) GetResultsForSeason(r *proto.SeasonRequest, stream proto.ResultService_GetResultsForSeasonServer) error {
	date, err := time.Parse(time.RFC3339, r.DateBefore)

	if err != nil {
		return status.Error(codes.InvalidArgument, "Date provided is not a valid RFC3339 date")
	}

	id := uint64(r.SeasonId)

	query := app.FixtureRepositoryQuery{
		SeasonID: &id,
		DateTo:   &date,
	}

	fixtures, err := s.fixtureRepo.Get(query)

	if err != nil {
		s.logger.Warnf("Error retrieving Fixture(s) in Result Service. Error: %s", err.Error())
		return status.Error(codes.Internal, "Internal server error")
	}

	return s.sendResults(fixtures, stream)
}

func (s ResultService) sendResults(f []app.Fixture, stream proto.ResultService_GetResultsForTeamServer) error {
	for _, fix := range f {
		x, err := s.factory.BuildResult(&fix)

		if err != nil {
			s.logger.Warnf("Error hydrating Result. Error: %s", err.Error())
			return status.Error(codes.Internal, "Internal server error")
		}

		if err := stream.Send(x); err != nil {
			s.logger.Warnf("Error streaming Result back to client. Error: %s", err.Error())
			return status.Error(codes.Internal, "Internal server error")
		}
	}

	return nil
}

func NewResultService(r app.FixtureRepository, f *factory.ResultFactory, log *logrus.Logger) *ResultService {
	return &ResultService{fixtureRepo: r, factory: f, logger: log}
}
