package mock

import (
	"github.com/statistico/statistico-proto/data/go"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type CompetitionServer struct {
	mock.Mock
	grpc.ServerStream
}

func (c *CompetitionServer) Send(comp *statisticoproto.Competition) error {
	args := c.Called(comp)
	return args.Error(0)
}

type SeasonServer struct {
	mock.Mock
	grpc.ServerStream
}

func (s *SeasonServer) Send(season *statisticoproto.Season) error {
	args := s.Called(season)
	return args.Error(0)
}

type TeamServer struct {
	mock.Mock
	grpc.ServerStream
}

func (t *TeamServer) Send(team *statisticoproto.Team) error {
	args := t.Called(team)
	return args.Error(0)
}

type TeamStatServer struct {
	mock.Mock
	grpc.ServerStream
}

func (t *TeamStatServer) Send(stat *statisticoproto.TeamStat) error {
	args := t.Called(stat)
	return args.Error(0)
}
