package grpc

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-football-data/internal/app"
	"github.com/statistico/statistico-football-data/internal/app/errors"
	"github.com/statistico/statistico-football-data/internal/app/grpc/factory"
	statistico "github.com/statistico/statistico-proto/go"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TeamService struct {
	teamRepo app.TeamRepository
	logger *logrus.Logger
	statistico.UnimplementedTeamServiceServer
}

func (t *TeamService) GetTeamByID(ctx context.Context, r *statistico.TeamRequest) (*statistico.Team, error) {
	team, err := t.teamRepo.ByID(r.TeamId)

	if err != nil {
		if err == errors.ErrorNotFound {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("team with ID %d does not exist", r.TeamId))
		}

		t.logger.Errorf("error fetching team in gRPC team service: %s", err.Error())

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return factory.TeamToProto(team), nil
}

func (t *TeamService) GetTeamsBySeasonId(r *statistico.SeasonTeamsRequest, stream statistico.TeamService_GetTeamsBySeasonIdServer) error {
	teams, err := t.teamRepo.BySeasonId(r.GetSeasonId())

	if err != nil {
		t.logger.Errorf("Error retrieving Team(s) in Team Service. Error: %s", err.Error())
		return status.Error(codes.Internal, "Internal server error")
	}

	for _, team := range teams {
		te := factory.TeamToProto(&team)

		if err := stream.Send(te); err != nil {
			t.logger.Errorf("Error streaming Team back to client. Error: %s", err.Error())
			return status.Error(codes.Internal, "Internal server error")
		}
	}

	return nil
}

func NewTeamService(r app.TeamRepository, l *logrus.Logger) *TeamService {
	return &TeamService{
		teamRepo: r,
		logger:   l,
	}
}
