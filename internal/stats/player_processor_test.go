package stats

import (
	"github.com/stretchr/testify/mock"
	"github.com/joesweeny/statshub/internal/model"
	"testing"
	"github.com/jonboulle/clockwork"
	"log"
	"io/ioutil"
)

func TestProcessPlayerStats(t *testing.T) {
	repo := new(mockPlayerRepository)
	processor := PlayerProcessor{
		PlayerRepository: repo,
		PlayerFactory:    PlayerFactory{clockwork.NewFakeClock()},
		Logger:     log.New(ioutil.Discard, "", 0),
	}

	player := newClientLineupPlayer()

	t.Run("creates new stats struct and inserts if not present in database", func(t *testing.T) {
		repo.On("ByFixtureAndPlayer", 1203, 20918).Return(&model.PlayerStats{}, ErrNotFound)
		repo.On("InsertPlayerStats", mock.Anything).Return(nil)
		processor.ProcessPlayerStats(player, true)
	})

	t.Run("stats struct if not inserted if already present in database", func(t *testing.T) {
		repo.On("ByFixtureAndPlayer", 1203, 20918).Return(&model.PlayerStats{}, nil)
		repo.AssertNotCalled(t,"InsertPlayerStats", mock.Anything)
		processor.ProcessPlayerStats(player, true)
	})
}

type mockPlayerRepository struct {
	mock.Mock
}

func (m mockPlayerRepository) InsertPlayerStats(s *model.PlayerStats) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m mockPlayerRepository) UpdatePlayerStats(s *model.PlayerStats) error {
	args := m.Called(s)
	return args.Error(0)
}

func (m mockPlayerRepository) ByFixtureAndPlayer(fixtureId, playerId int) (*model.PlayerStats, error) {
	args := m.Called(fixtureId, playerId)
	return args.Get(0).(*model.PlayerStats), args.Error(1)
}
