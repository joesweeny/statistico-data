package grpc

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/errors"
	"github.com/statistico/statistico-data/internal/app/grpc/factory"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TeamService struct {
	teamRepo app.TeamRepository
	logger *logrus.Logger
}

func (t *TeamService) GetTeamByID(ctx context.Context, r *proto.TeamRequest) (*proto.Team, error) {
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

func NewTeamService(r app.TeamRepository, l *logrus.Logger) *TeamService {
	return &TeamService{
		teamRepo: r,
		logger:   l,
	}
}
