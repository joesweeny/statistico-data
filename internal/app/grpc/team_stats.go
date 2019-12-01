package grpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/factory"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
)

type TeamStatsService struct {
	fixtureRepository app.FixtureRepository
	xGRepo            app.FixtureTeamXGRepository
	factory           *factory.TeamStatsFactory
	logger            *logrus.Logger
}

func (s TeamStatsService) GetTeamStatsForFixture(c context.Context, r *proto.FixtureRequest) (*proto.TeamStatsResponse, error) {
	fix, err := s.fixtureRepository.ByID(r.FixtureId)

	if err != nil {
		return nil, fmt.Errorf("fixture with ID %d does not exist", r.FixtureId)
	}

	home, err := s.factory.BuildTeamStats(fix, fix.HomeTeamID)

	if err != nil {
		s.logger.Warnf("Error hydrating proto team stats: %s", err.Error())
		return nil, internalServerError
	}

	away, err := s.factory.BuildTeamStats(fix, fix.AwayTeamID)

	if err != nil {
		s.logger.Warnf("Error hydrating proto team stats: %s", err.Error())
		return nil, internalServerError
	}

	xg, err := s.xGRepo.ByFixtureID(fix.ID)

	if err != nil {
		s.logger.Warnf("Error hydrating proto team stats: %s", err.Error())
		return nil, internalServerError
	}

	res := proto.TeamStatsResponse{
		HomeTeam: home,
		AwayTeam: away,
		TeamXg: &proto.TeamXG{
			Home: parseXgRating(xg.Home),
			Away: parseXgRating(xg.Away),
		},
	}

	return &res, nil
}

func NewTeamStatsService(
	r app.FixtureRepository,
	x app.FixtureTeamXGRepository,
	f *factory.TeamStatsFactory,
	log *logrus.Logger,
) *TeamStatsService {
	return &TeamStatsService{fixtureRepository: r, xGRepo: x, factory: f, logger: log}
}

func parseXgRating(xg *float32) *wrappers.FloatValue {
	if xg != nil {
		return &wrappers.FloatValue{
			Value: *xg,
		}
	}

	return nil
}
