package grpc

import (
	"context"
	"github.com/statistico/statistico-football-data/internal/app/performance"
	statistico "github.com/statistico/statistico-proto/go"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PerformanceService struct {
	reader performance.StatReader
	statistico.UnimplementedPerformanceServiceServer
}

func (s *PerformanceService) GetTeamsMatchingStat(c context.Context, r *statistico.TeamStatPerformanceRequest) (*statistico.TeamStatResponse, error) {
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

	res := statistico.TeamStatResponse{Teams: convertTeams(teams)}

	return &res, nil
}

func convertTeams(t []*performance.Team) []*statistico.Team {
	var teams []*statistico.Team

	for _, team := range t {
		x := statistico.Team{
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
