package mock

import (
	"github.com/statistico/statistico-football-data/internal/app/performance"
	"github.com/stretchr/testify/mock"
)

type StatReader struct {
	mock.Mock
}

func (s *StatReader) TeamsMatchingFilter(f *performance.StatFilter) ([]*performance.Team, error) {
	args := s.Called(f)
	return args.Get(0).([]*performance.Team), args.Error(1)
}
