package grpc

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CompetitionService struct {
	competitionRepo app.CompetitionRepository
	logger *logrus.Logger
}

func (s *CompetitionService) ListCompetitions(r *proto.CompetitionRequest, stream proto.CompetitionService_ListCompetitionsServer) error {
	var query app.CompetitionFilterQuery
	
	if len(r.GetCountryIds()) > 0 {
		query.CountryIds = r.GetCountryIds()
	}
	
	if r.GetIsCup() != nil {
		v := r.GetIsCup().GetValue()
		query.IsCup = &v
	}
	
	if r.GetSort() != nil {
		v := r.GetSort().GetValue()
		query.SortBy = &v
	}
	
	competitions, err := s.competitionRepo.Get(query)

	if err != nil {
		s.logger.Warnf("Error retrieving Competition(s) in Competition Service. Error: %s", err.Error())
		return status.Error(codes.Internal, "Internal server error")
	}
	
	for _, comp := range competitions {
		c := proto.Competition{
			Id:    comp.ID,
			Name:  comp.Name,
			IsCup: comp.IsCup,
		}

		if err := stream.Send(&c); err != nil {
			s.logger.Warnf("Error streaming Competition back to client. Error: %s", err.Error())
			return status.Error(codes.Internal, "Internal server error")
		}
	}

	return nil
}
