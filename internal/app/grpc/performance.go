package grpc

import (
	"context"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"github.com/statistico/statistico-data/internal/app/performance"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PerformanceService struct {
	reader performance.StatReader
}

func (s *PerformanceService) GetTeamsMatchingStat(c context.Context, r *proto.TeamStatRequest) (*proto.TeamStatResponse, error) {
	q := performance.StatFilter{
		Action:  r.GetAction(),
		Games:   uint8(r.GetGames()),
		Measure: r.GetMeasure(),
		Metric:  r.GetMetric(),
		Seasons: r.GetSeasons(),
		Stat:    r.GetStat(),
		Value:   r.GetValue(),
		Venue:   r.GetVenue(),
	}

	teams, err := s.reader.TeamsMatchingFilter(&q)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	res := proto.TeamStatResponse{Teams: convertTeams(teams)}

	return &res, nil
}

func convertTeams(t []*performance.Team) []*proto.Team {
	var teams []*proto.Team

	for _, team := range t {
		x := proto.Team{
			Id:   int64(team.ID),
			Name: team.Name,
		}

		teams = append(teams, &x)
	}

	return teams
}

func NewPerformanceService(r performance.StatReader) *PerformanceService {
	return &PerformanceService{reader: r}
}
