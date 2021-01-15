package mock

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type PlayerStatsRepository struct {
	mock.Mock
}

func (m PlayerStatsRepository) Insert(p *app.PlayerStats) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m PlayerStatsRepository) Update(p *app.PlayerStats) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m PlayerStatsRepository) ByFixtureAndPlayer(fixtureID, playerID uint64) (*app.PlayerStats, error) {
	args := m.Called(fixtureID, playerID)
	c := args.Get(0).(*app.PlayerStats)
	return c, args.Error(1)
}

func (m PlayerStatsRepository) ByFixtureAndTeam(fixtureID, teamID uint64) ([]*app.PlayerStats, error) {
	args := m.Called(fixtureID, teamID)
	c := args.Get(0).([]*app.PlayerStats)
	return c, args.Error(1)
}

type PlayerStatsRequester struct {
	mock.Mock
}

func (m PlayerStatsRequester) PlayerStatsByFixtureIDs(ids []uint64) <-chan *app.PlayerStats {
	args := m.Called(ids)
	return args.Get(0).(chan *app.PlayerStats)
}

func (m PlayerStatsRequester) PlayerStatsBySeasonIDs(ids []uint64) <-chan *app.PlayerStats {
	args := m.Called(ids)
	return args.Get(0).(chan *app.PlayerStats)
}
