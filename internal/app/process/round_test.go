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

func TestRoundProcessor_Process(t *testing.T) {
	t.Run("inserts new round into repository when processing round command", func(t *testing.T) {
		t.Helper()

		roundRepo := new(mock.RoundRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.RoundRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewRoundProcessor(roundRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newRound(45)
		two := newRound(51)

		rounds := make([]*app.Round, 2)
		rounds[0] = one
		rounds[1] = two

		ch := roundChannel(rounds)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("RoundsBySeasonIDs", ids).Return(ch)

		roundRepo.On("ByID", uint64(45)).Return(&app.Round{}, errors.New("not Found"))
		roundRepo.On("ByID", uint64(51)).Return(&app.Round{}, errors.New("not Found"))
		roundRepo.On("Insert", one).Return(nil)
		roundRepo.On("Insert", two).Return(nil)

		processor.Process("round", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		roundRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing round into repository when processing round command", func(t *testing.T) {
		t.Helper()

		roundRepo := new(mock.RoundRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.RoundRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewRoundProcessor(roundRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newRound(45)
		two := newRound(51)

		rounds := make([]*app.Round, 2)
		rounds[0] = one
		rounds[1] = two

		ch := roundChannel(rounds)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("RoundsBySeasonIDs", ids).Return(ch)

		roundRepo.On("ByID", uint64(45)).Return(one, nil)
		roundRepo.On("ByID", uint64(51)).Return(two, nil)
		roundRepo.On("Update", &one).Return(nil)
		roundRepo.On("Update", &two).Return(nil)

		processor.Process("round", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		roundRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert round into repository when processing round command", func(t *testing.T) {
		t.Helper()

		roundRepo := new(mock.RoundRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.RoundRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewRoundProcessor(roundRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newRound(45)
		two := newRound(51)

		rounds := make([]*app.Round, 2)
		rounds[0] = one
		rounds[1] = two

		ch := roundChannel(rounds)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("RoundsBySeasonIDs", ids).Return(ch)

		roundRepo.On("ByID", uint64(45)).Return(&app.Round{}, errors.New("not Found"))
		roundRepo.On("ByID", uint64(51)).Return(two, nil)
		roundRepo.On("Insert", one).Return(errors.New("error occurred"))
		roundRepo.On("Update", &two).Return(nil)

		processor.Process("round", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		roundRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update round into repository when processing round command", func(t *testing.T) {
		t.Helper()

		roundRepo := new(mock.RoundRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.RoundRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewRoundProcessor(roundRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newRound(45)
		two := newRound(51)

		rounds := make([]*app.Round, 2)
		rounds[0] = one
		rounds[1] = two

		ch := roundChannel(rounds)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("RoundsBySeasonIDs", ids).Return(ch)

		roundRepo.On("ByID", uint64(45)).Return(&app.Round{}, errors.New("not Found"))
		roundRepo.On("ByID", uint64(51)).Return(two, nil)
		roundRepo.On("Insert", one).Return(nil)
		roundRepo.On("Update", &two).Return(errors.New("error occurred"))

		processor.Process("round", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		roundRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new round into repository when processing round current season command", func(t *testing.T) {
		t.Helper()

		roundRepo := new(mock.RoundRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.RoundRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewRoundProcessor(roundRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newRound(45)
		two := newRound(51)

		rounds := make([]*app.Round, 2)
		rounds[0] = one
		rounds[1] = two

		ch := roundChannel(rounds)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("RoundsBySeasonIDs", ids).Return(ch)

		roundRepo.On("ByID", uint64(45)).Return(&app.Round{}, errors.New("not Found"))
		roundRepo.On("ByID", uint64(51)).Return(&app.Round{}, errors.New("not Found"))
		roundRepo.On("Insert", one).Return(nil)
		roundRepo.On("Insert", two).Return(nil)

		processor.Process("round:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		roundRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing round into repository when processing round current season command", func(t *testing.T) {
		t.Helper()

		roundRepo := new(mock.RoundRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.RoundRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewRoundProcessor(roundRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newRound(45)
		two := newRound(51)

		rounds := make([]*app.Round, 2)
		rounds[0] = one
		rounds[1] = two

		ch := roundChannel(rounds)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("RoundsBySeasonIDs", ids).Return(ch)

		roundRepo.On("ByID", uint64(45)).Return(one, nil)
		roundRepo.On("ByID", uint64(51)).Return(two, nil)
		roundRepo.On("Update", &one).Return(nil)
		roundRepo.On("Update", &two).Return(nil)

		processor.Process("round:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		roundRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert round into repository when processing round current season command", func(t *testing.T) {
		t.Helper()

		roundRepo := new(mock.RoundRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.RoundRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewRoundProcessor(roundRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newRound(45)
		two := newRound(51)

		rounds := make([]*app.Round, 2)
		rounds[0] = one
		rounds[1] = two

		ch := roundChannel(rounds)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("RoundsBySeasonIDs", ids).Return(ch)

		roundRepo.On("ByID", uint64(45)).Return(&app.Round{}, errors.New("not Found"))
		roundRepo.On("ByID", uint64(51)).Return(two, nil)
		roundRepo.On("Insert", one).Return(errors.New("error occurred"))
		roundRepo.On("Update", &two).Return(nil)

		processor.Process("round:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		roundRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update round into repository when processing round current season command", func(t *testing.T) {
		t.Helper()

		roundRepo := new(mock.RoundRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.RoundRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewRoundProcessor(roundRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newRound(45)
		two := newRound(51)

		rounds := make([]*app.Round, 2)
		rounds[0] = one
		rounds[1] = two

		ch := roundChannel(rounds)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("RoundsBySeasonIDs", ids).Return(ch)

		roundRepo.On("ByID", uint64(45)).Return(&app.Round{}, errors.New("not Found"))
		roundRepo.On("ByID", uint64(51)).Return(two, nil)
		roundRepo.On("Insert", one).Return(nil)
		roundRepo.On("Update", &two).Return(errors.New("error occurred"))

		processor.Process("round:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		roundRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
}

func newRound(id uint64) *app.Round {
	return &app.Round{
		ID:        id,
		Name:      "5",
		SeasonID:  uint64(4387),
		StartDate: time.Unix(1548086929, 0),
		EndDate:   time.Unix(1548086929, 0),
	}
}

func roundChannel(rounds []*app.Round) chan *app.Round {
	ch := make(chan *app.Round, len(rounds))

	for _, c := range rounds {
		ch <- c
	}

	close(ch)

	return ch
}
