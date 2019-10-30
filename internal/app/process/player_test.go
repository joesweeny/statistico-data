package process_test

import (
	"errors"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/process"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPlayerProcessor_Process(t *testing.T) {
	t.Run("inserts player if does not already exist", func(t *testing.T) {
		t.Helper()

		playerRepo := new(mock.PlayerRepository)
		squadRepo := new(mock.SquadRepository)
		requester := new(mock.PlayerRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewPlayerProcessor(playerRepo, squadRepo, requester, logger)

		done := make(chan bool)

		def := newPlayer(1)
		mid := newPlayer(2)
		str := newPlayer(3)

		squad := newSquad()

		squadRepo.On("All").Return(squad, nil)

		playerRepo.On("ByID", int64(1)).Return(&app.Player{}, errors.New("not found"))
		playerRepo.On("ByID", int64(2)).Return(&app.Player{}, errors.New("not found"))
		playerRepo.On("ByID", int64(3)).Return(&app.Player{}, errors.New("not found"))

		requester.On("PlayerByID", int64(1)).Return(def)
		requester.On("PlayerByID", int64(2)).Return(mid)
		requester.On("PlayerByID", int64(3)).Return(str)

		playerRepo.On("Insert", def).Return(nil)
		playerRepo.On("Insert", mid).Return(nil)
		playerRepo.On("Insert", str).Return(nil)

		processor.Process("player", "", done)

		<-done

		requester.AssertExpectations(t)
		playerRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("does not insert player if already exists", func(t *testing.T) {
		t.Helper()

		playerRepo := new(mock.PlayerRepository)
		squadRepo := new(mock.SquadRepository)
		requester := new(mock.PlayerRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewPlayerProcessor(playerRepo, squadRepo, requester, logger)

		done := make(chan bool)

		def := newPlayer(1)
		mid := newPlayer(2)
		str := newPlayer(3)

		squad := newSquad()

		squadRepo.On("All").Return(squad, nil)

		playerRepo.On("ByID", int64(1)).Return(&app.Player{}, errors.New("not found"))
		playerRepo.On("ByID", int64(2)).Return(&app.Player{}, nil)
		playerRepo.On("ByID", int64(3)).Return(&app.Player{}, errors.New("not found"))

		requester.On("PlayerByID", int64(1)).Return(def)
		requester.On("PlayerByID", int64(3)).Return(str)

		requester.AssertNotCalled(t, "PlayerByID", int64(2))

		playerRepo.On("Insert", def).Return(nil)
		playerRepo.On("Insert", str).Return(nil)

		playerRepo.AssertNotCalled(t, "Insert", mid)

		processor.Process("player", "", done)

		<-done

		requester.AssertExpectations(t)
		playerRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error if cannot insert player into repository", func(t *testing.T) {
		t.Helper()

		playerRepo := new(mock.PlayerRepository)
		squadRepo := new(mock.SquadRepository)
		requester := new(mock.PlayerRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewPlayerProcessor(playerRepo, squadRepo, requester, logger)

		done := make(chan bool)

		def := newPlayer(1)
		mid := newPlayer(2)
		str := newPlayer(3)

		squad := newSquad()

		squadRepo.On("All").Return(squad, nil)

		playerRepo.On("ByID", int64(1)).Return(&app.Player{}, errors.New("not found"))
		playerRepo.On("ByID", int64(2)).Return(&app.Player{}, nil)
		playerRepo.On("ByID", int64(3)).Return(&app.Player{}, errors.New("not found"))

		requester.On("PlayerByID", int64(1)).Return(def)
		requester.On("PlayerByID", int64(3)).Return(str)

		requester.AssertNotCalled(t, "PlayerByID", int64(2))

		playerRepo.On("Insert", def).Return(nil)
		playerRepo.On("Insert", str).Return(errors.New("cannot insert"))

		playerRepo.AssertNotCalled(t, "Insert", mid)

		processor.Process("player", "", done)

		<-done
		
		requester.AssertExpectations(t)
		playerRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.NotNil(t, hook.LastEntry().Message)
	})
}

func newPlayer(id int64) *app.Player {
	return &app.Player{
		ID:          id,
		CountryId:   int64(154),
		FirstName:   "Manuel",
		LastName:    "Lanzini",
	}
}

func newSquad() []model.Squad {
	var squads []model.Squad

	s := model.Squad{
		SeasonID:  45,
		TeamID:    98,
		PlayerIDs: []int{1, 2, 3},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	squads = append(squads, s)

	return squads
}
