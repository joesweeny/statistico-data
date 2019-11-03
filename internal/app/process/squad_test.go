package process_test

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/process"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSquadProcessor_Process(t *testing.T) {
	t.Run("inserts new squad into repository when processing squad command", func(t *testing.T) {
		t.Helper()

		squadRepo := new(mock.SquadRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.SquadRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSquadProcessor(squadRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newSquad(12962, 1)
		ncu := newSquad(12962, 14)

		squads := make([]*app.Squad, 2)
		squads[0] = whu
		squads[1] = ncu

		ch := squadChannel(squads)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("SquadsBySeasonIDs", ids).Return(ch)

		squadRepo.On("BySeasonAndTeam", uint64(1), uint64(12962)).Return(&app.Squad{}, errors.New("not Found"))
		squadRepo.On("BySeasonAndTeam", uint64(14), uint64(12962)).Return(&app.Squad{}, errors.New("not Found"))
		squadRepo.On("Insert", whu).Return(nil)
		squadRepo.On("Insert", ncu).Return(nil)

		processor.Process("squad", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing squad into repository when processing squad command", func(t *testing.T) {
		t.Helper()

		squadRepo := new(mock.SquadRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.SquadRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSquadProcessor(squadRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newSquad(12962, 1)
		ncu := newSquad(12962, 14)

		squads := make([]*app.Squad, 2)
		squads[0] = whu
		squads[1] = ncu

		ch := squadChannel(squads)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("SquadsBySeasonIDs", ids).Return(ch)

		squadRepo.On("BySeasonAndTeam", uint64(1), uint64(12962)).Return(whu, nil)
		squadRepo.On("BySeasonAndTeam", uint64(14), uint64(12962)).Return(ncu, nil)
		squadRepo.On("Update", &whu).Return(nil)
		squadRepo.On("Update", &ncu).Return(nil)

		processor.Process("squad", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert squad into repository when processing squad command", func(t *testing.T) {
		t.Helper()

		squadRepo := new(mock.SquadRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.SquadRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSquadProcessor(squadRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newSquad(12962, 1)
		ncu := newSquad(12962, 14)

		squads := make([]*app.Squad, 2)
		squads[0] = whu
		squads[1] = ncu

		ch := squadChannel(squads)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("SquadsBySeasonIDs", ids).Return(ch)

		squadRepo.On("BySeasonAndTeam", uint64(1), uint64(12962)).Return(&app.Squad{}, errors.New("not Found"))
		squadRepo.On("BySeasonAndTeam", uint64(14), uint64(12962)).Return(&app.Squad{}, errors.New("not Found"))
		squadRepo.On("Insert", whu).Return(errors.New("error occurred"))
		squadRepo.On("Insert", ncu).Return(nil)

		processor.Process("squad", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update squad into repository when processing squad command", func(t *testing.T) {
		t.Helper()

		squadRepo := new(mock.SquadRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.SquadRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSquadProcessor(squadRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newSquad(12962, 1)
		ncu := newSquad(12962, 14)

		squads := make([]*app.Squad, 2)
		squads[0] = whu
		squads[1] = ncu

		ch := squadChannel(squads)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("SquadsBySeasonIDs", ids).Return(ch)

		squadRepo.On("BySeasonAndTeam", uint64(1), uint64(12962)).Return(whu, nil)
		squadRepo.On("BySeasonAndTeam", uint64(14), uint64(12962)).Return(ncu, nil)
		squadRepo.On("Update", &whu).Return(errors.New("error occurred"))
		squadRepo.On("Update", &ncu).Return(nil)

		processor.Process("squad", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new squad into repository when processing squad current season command", func(t *testing.T) {
		t.Helper()

		squadRepo := new(mock.SquadRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.SquadRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSquadProcessor(squadRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newSquad(12962, 1)
		ncu := newSquad(12962, 14)

		squads := make([]*app.Squad, 2)
		squads[0] = whu
		squads[1] = ncu

		ch := squadChannel(squads)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("SquadsBySeasonIDs", ids).Return(ch)

		squadRepo.On("BySeasonAndTeam", uint64(1), uint64(12962)).Return(&app.Squad{}, errors.New("not Found"))
		squadRepo.On("BySeasonAndTeam", uint64(14), uint64(12962)).Return(&app.Squad{}, errors.New("not Found"))
		squadRepo.On("Insert", whu).Return(nil)
		squadRepo.On("Insert", ncu).Return(nil)

		processor.Process("squad:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing squad into repository when processing squad current season command", func(t *testing.T) {
		t.Helper()

		squadRepo := new(mock.SquadRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.SquadRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSquadProcessor(squadRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newSquad(12962, 1)
		ncu := newSquad(12962, 14)

		squads := make([]*app.Squad, 2)
		squads[0] = whu
		squads[1] = ncu

		ch := squadChannel(squads)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("SquadsBySeasonIDs", ids).Return(ch)

		squadRepo.On("BySeasonAndTeam", uint64(1), uint64(12962)).Return(whu, nil)
		squadRepo.On("BySeasonAndTeam", uint64(14), uint64(12962)).Return(ncu, nil)
		squadRepo.On("Update", &whu).Return(nil)
		squadRepo.On("Update", &ncu).Return(nil)

		processor.Process("squad:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert squad into repository when processing squad current season command", func(t *testing.T) {
		t.Helper()

		squadRepo := new(mock.SquadRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.SquadRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSquadProcessor(squadRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newSquad(12962, 1)
		ncu := newSquad(12962, 14)

		squads := make([]*app.Squad, 2)
		squads[0] = whu
		squads[1] = ncu

		ch := squadChannel(squads)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("SquadsBySeasonIDs", ids).Return(ch)

		squadRepo.On("BySeasonAndTeam", uint64(1), uint64(12962)).Return(&app.Squad{}, errors.New("not Found"))
		squadRepo.On("BySeasonAndTeam", uint64(14), uint64(12962)).Return(&app.Squad{}, errors.New("not Found"))
		squadRepo.On("Insert", whu).Return(errors.New("error occurred"))
		squadRepo.On("Insert", ncu).Return(nil)

		processor.Process("squad:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update squad into repository when processing squad current season command", func(t *testing.T) {
		t.Helper()

		squadRepo := new(mock.SquadRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.SquadRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSquadProcessor(squadRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newSquad(12962, 1)
		ncu := newSquad(12962, 14)

		squads := make([]*app.Squad, 2)
		squads[0] = whu
		squads[1] = ncu

		ch := squadChannel(squads)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("SquadsBySeasonIDs", ids).Return(ch)

		squadRepo.On("BySeasonAndTeam", uint64(1), uint64(12962)).Return(whu, nil)
		squadRepo.On("BySeasonAndTeam", uint64(14), uint64(12962)).Return(ncu, nil)
		squadRepo.On("Update", &whu).Return(errors.New("error occurred"))
		squadRepo.On("Update", &ncu).Return(nil)

		processor.Process("squad:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		squadRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
}

func newSquad(season, team uint64) *app.Squad {
	return &app.Squad{
		SeasonID:  season,
		TeamID:    team,
		PlayerIDs: []uint64{34, 57, 89},
		CreatedAt: time.Unix(1547465100, 0),
		UpdatedAt: time.Unix(1547465100, 0),
	}
}

func squadChannel(squads []*app.Squad) chan *app.Squad {
	ch := make(chan *app.Squad, len(squads))

	for _, c := range squads {
		ch <- c
	}

	close(ch)

	return ch
}
