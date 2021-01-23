package mock

import (
	"github.com/statistico/statistico-proto/go"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type CompetitionServer struct {
	mock.Mock
	grpc.ServerStream
}

func (c *CompetitionServer) Send(comp *statistico.Competition) error {
	args := c.Called(comp)
	return args.Error(0)
}

type SeasonServer struct {
	mock.Mock
	grpc.ServerStream
}

func (s *SeasonServer) Send(season *statistico.Season) error {
	args := s.Called(season)
	return args.Error(0)
}

type TeamServer struct {
	mock.Mock
	grpc.ServerStream
}

func (t *TeamServer) Send(team *statistico.Team) error {
	args := t.Called(team)
	return args.Error(0)
}

type TeamStatServer struct {
	mock.Mock
	grpc.ServerStream
}

func (t *TeamStatServer) Send(stat *statistico.TeamStat) error {
	args := t.Called(stat)
	return args.Error(0)
}
