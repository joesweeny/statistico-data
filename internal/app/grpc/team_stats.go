package grpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/factory"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TeamStatsService struct {
	fixtureRepository app.FixtureRepository
	statRepository    app.TeamStatsRepository
	xGRepo            app.FixtureTeamXGRepository
	factory           *factory.TeamStatsFactory
	logger            *logrus.Logger
}

func (s TeamStatsService) GetStatForTeam(r *proto.TeamStatRequest, stream proto.TeamStatsService_GetStatForTeamServer) error {
	query, err := fixtureFilterFromTeamStatRequest(r)

	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	fixtures, err := s.fixtureRepository.ByTeamID(r.GetTeamId(), *query)

	if err != nil {
		s.logger.Warnf("Error retrieving fixture(s) in team stats service. Error: %s", err.Error())
		return status.Error(codes.Internal, "Internal server error")
	}

	for _, fix := range fixtures {
		stat, err := s.parseTeamStat(fix, r.GetStat(), r.GetTeamId(), r.GetOpponent().GetValue())

		if err != nil {
			s.logger.Warnf("Error retrieving team stat in team stats service. Error: %s", err.Error())
			continue
		}

		x := factory.TeamStatToProto(stat)

		if err := stream.Send(x); err != nil {
			s.logger.Warnf("Error streaming team stat back to client. Error: %s", err.Error())
			return status.Error(codes.Internal, "Internal server error")
		}
	}

	return nil
}

func (s TeamStatsService) GetTeamStatsForFixture(c context.Context, r *proto.FixtureRequest) (*proto.TeamStatsResponse, error) {
	fix, err := s.fixtureRepository.ByID(r.FixtureId)

	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("fixture with ID %d does not exist", r.FixtureId))
	}

	home, err := s.factory.BuildTeamStats(fix, fix.HomeTeamID)

	if err != nil {
		s.logger.Warnf("Error hydrating proto team stats: %s", err.Error())
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("home team stats do not exist for fixture %d", r.FixtureId),
		)
	}

	away, err := s.factory.BuildTeamStats(fix, fix.AwayTeamID)

	if err != nil {
		s.logger.Warnf("Error hydrating proto team stats: %s", err.Error())
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("away team stats do not exist for fixture %d", r.FixtureId),
		)
	}

	xg, err := s.xGRepo.ByFixtureID(fix.ID)

	if err != nil {
		s.logger.Warnf("Error hydrating proto team stats: %s", err.Error())
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("xG stats do not exist for fixture %d", r.FixtureId),
		)
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

func (s TeamStatsService) parseTeamStat(f app.Fixture, stat string, teamID uint64, opponent bool) (*app.TeamStat, error) {
	id := parseTeamId(f, teamID, opponent)

	x, err := s.statRepository.StatByFixtureAndTeam(stat, f.ID, id)

	if err != nil {
		return nil, err
	}

	return x, nil
}

func NewTeamStatsService(
	r app.FixtureRepository,
	s app.TeamStatsRepository,
	x app.FixtureTeamXGRepository,
	f *factory.TeamStatsFactory,
	log *logrus.Logger,
) *TeamStatsService {
	return &TeamStatsService{fixtureRepository: r, statRepository: s, xGRepo: x, factory: f, logger: log}
}

func parseXgRating(xg *float32) *wrappers.FloatValue {
	if xg != nil {
		return &wrappers.FloatValue{
			Value: *xg,
		}
	}

	return nil
}

func parseTeamId(f app.Fixture, teamID uint64, opponent bool) uint64 {
	if opponent && (f.HomeTeamID == teamID) {
		return f.AwayTeamID
	}

	if opponent && (f.AwayTeamID == teamID) {
		return f.HomeTeamID
	}

	return teamID
}
