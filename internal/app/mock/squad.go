package mock

import (
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/mock"
)

type SquadRepository struct {
	mock.Mock
}

func (m SquadRepository) Insert(c *model.Squad) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m SquadRepository) Update(c *model.Squad) error {
	args := m.Called(&c)
	return args.Error(0)
}

func (m SquadRepository) BySeasonAndTeam(seasonId, teamId int) (*model.Squad, error) {
	args := m.Called(seasonId, teamId)
	c := args.Get(0).(*model.Squad)
	return c, args.Error(1)
}

func (m SquadRepository) All() ([]model.Squad, error) {
	args := m.Called()
	return args.Get(0).([]model.Squad), args.Error(1)
}

func (m SquadRepository) CurrentSeason() ([]model.Squad, error) {
	args := m.Called()
	return args.Get(0).([]model.Squad), args.Error(1)
}
