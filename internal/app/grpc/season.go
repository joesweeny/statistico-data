package grpc

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/factory"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SeasonService struct {
	seasonRepo app.SeasonRepository
	logger *logrus.Logger
}

func (s *SeasonService) GetSeasonsForCompetition(
	r *proto.SeasonCompetitionRequest,
	stream proto.SeasonService_GetSeasonsForCompetitionServer,
) error {
	seasons, err := s.seasonRepo.ByCompetitionId(r.GetCompetitionId(), r.GetSort().GetValue())

	if err != nil {
		s.logger.Errorf("Error retrieving Season(s) in Season Service. Error: %s", err.Error())
		return status.Error(codes.Internal, "Internal server error")
	}

	for _, season := range seasons {
		se := factory.SeasonToProto(&season)

		if err := stream.Send(se); err != nil {
			s.logger.Errorf("Error streaming Season back to client. Error: %s", err.Error())
			return status.Error(codes.Internal, "Internal server error")
		}
	}

	return nil
}

func NewSeasonService(r app.SeasonRepository, l *logrus.Logger) *SeasonService {
	return &SeasonService{
		seasonRepo: r,
		logger:     l,
	}
}
