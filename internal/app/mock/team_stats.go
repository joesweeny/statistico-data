package mock

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/mock"
)

type TeamStatsRepository struct {
	mock.Mock
}

func (m *TeamStatsRepository) InsertTeamStats(t *app.TeamStats) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *TeamStatsRepository) UpdateTeamStats(t *app.TeamStats) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *TeamStatsRepository) ByFixtureAndTeam(fixtureID, teamID uint64) (*app.TeamStats, error) {
	args := m.Called(fixtureID, teamID)
	c := args.Get(0).(*app.TeamStats)
	return c, args.Error(1)
}

func (m *TeamStatsRepository) StatByFixtureAndTeam(stat string, fixtureID, teamID uint64) (*app.TeamStat, error) {
	args := m.Called(stat, fixtureID, teamID)
	return args.Get(0).(*app.TeamStat), args.Error(1)
}

func (m *TeamStatsRepository) Get() ([]*app.TeamStats, error) {
	args := m.Called()
	return args.Get(0).([]*app.TeamStats), args.Error(1)
}

type TeamStatsRequester struct {
	mock.Mock
}

func (t *TeamStatsRequester) TeamStatsByFixtureIDs(ids []uint64) <-chan *app.TeamStats {
	args := t.Called(ids)
	return args.Get(0).(chan *app.TeamStats)
}
