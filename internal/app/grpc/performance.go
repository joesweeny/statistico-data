package grpc

import (
	"context"
	"github.com/statistico/statistico-data/internal/app/performance"
	"github.com/statistico/statistico-proto/data/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PerformanceService struct {
	reader performance.StatReader
}

func (s *PerformanceService) GetTeamsMatchingStat(c context.Context, r *statisticoproto.TeamStatPerformanceRequest) (*statisticoproto.TeamStatResponse, error) {
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

	res := statisticoproto.TeamStatResponse{Teams: convertTeams(teams)}

	return &res, nil
}

func convertTeams(t []*performance.Team) []*statisticoproto.Team {
	var teams []*statisticoproto.Team

	for _, team := range t {
		x := statisticoproto.Team{
			Id:   team.ID,
			Name: team.Name,
		}

		teams = append(teams, &x)
	}

	return teams
}

func NewPerformanceService(r performance.StatReader) *PerformanceService {
	return &PerformanceService{reader: r}
}
