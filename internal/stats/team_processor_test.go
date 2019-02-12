package stats

import (
	"github.com/stretchr/testify/mock"
	"github.com/jonboulle/clockwork"
	"log"
	"io/ioutil"
	"github.com/joesweeny/statshub/internal/model"
	"testing"
)

func TestProcessTeamStats(t *testing.T) {
	repo := new(mockTeamRepository)
	processor := TeamProcessor{
		TeamRepository: repo,
		TeamFactory:    TeamFactory{clockwork.NewFakeClock()},
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	team := newClientTeamStats()

	t.Run("creates new stats struct and inserts if not present in database", func(t *testing.T) {
		repo.On("ByFixtureAndTeam", 34019, 960).Return(&model.TeamStats{}, ErrNotFound)
		repo.On("InsertTeamStats", mock.Anything).Return(nil)
		processor.ProcessTeamStats(team)
	})

	t.Run("stats struct if not inserted if already present in database", func(t *testing.T) {
		repo.On("ByFixtureAndTeam", 1203, 20918).Return(&model.PlayerStats{}, nil)
		repo.AssertNotCalled(t,"InsertTeamStats", mock.Anything)
		processor.ProcessTeamStats(team)
	})
}

type mockTeamRepository struct {
	mock.Mock
}

func (m mockTeamRepository) InsertTeamStats(s *model.TeamStats) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m mockTeamRepository) UpdateTeamStats(s *model.TeamStats) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m mockTeamRepository) ByFixtureAndTeam(fixtureId, teamId int) (*model.TeamStats, error) {
	args := m.Called(fixtureId, teamId)
	return args.Get(0).(*model.TeamStats), args.Error(1)
}
