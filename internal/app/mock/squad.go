package mock

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type SquadRepository struct {
	mock.Mock
}

func (m SquadRepository) Insert(c *app.Squad) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m SquadRepository) Update(c *app.Squad) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m SquadRepository) BySeasonAndTeam(seasonId, teamId int64) (*app.Squad, error) {
	args := m.Called(seasonId, teamId)
	c := args.Get(0).(*app.Squad)
	return c, args.Error(1)
}

func (m SquadRepository) All() ([]app.Squad, error) {
	args := m.Called()
	return args.Get(0).([]app.Squad), args.Error(1)
}

func (m SquadRepository) CurrentSeason() ([]app.Squad, error) {
	args := m.Called()
	return args.Get(0).([]app.Squad), args.Error(1)
}

type SquadRequester struct {
	mock.Mock
}

func (m SquadRequester) SquadsBySeasonIDs(seasonIDs []int64) <-chan *app.Squad {
	args := m.Called(seasonIDs)
	return args.Get(0).(chan *app.Squad)
}
