package mock

import (
	understat "github.com/statistico/statistico-understat-parser"
	"github.com/stretchr/testify/mock"
)

type Parser struct {
	mock.Mock
}

func (m Parser) GetLeagueFixtures(league, season string) <-chan *understat.Fixture {
	args := m.Called(league, season)
	return args.Get(0).(chan *understat.Fixture)
}
