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

func TestSeasonProcessor_Process(t *testing.T) {
	t.Run("inserts new season", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		requester := new(mock.SeasonRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSeasonProcessor(repo, requester, logger)

		done := make(chan bool)

		current := newSeason(8, true)
		old := newSeason(2, false)

		seasons := make([]*app.Season, 2)
		seasons[0] = current
		seasons[1] = old

		ch :=seasonChannel(seasons)

		requester.On("Seasons").Return(ch)

		repo.On("ByID", int64(8)).Return(&app.Season{}, errors.New("not found"))
		repo.On("ByID", int64(2)).Return(&app.Season{}, errors.New("not found"))
		repo.On("Insert", current).Return(nil)
		repo.On("Insert", old).Return(nil)

		processor.Process("season", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing season", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		requester := new(mock.SeasonRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSeasonProcessor(repo, requester, logger)

		done := make(chan bool)

		current := newSeason(8, true)
		old := newSeason(2, false)

		seasons := make([]*app.Season, 2)
		seasons[0] = current
		seasons[1] = old

		ch :=seasonChannel(seasons)

		requester.On("Seasons").Return(ch)

		repo.On("ByID", int64(8)).Return(current, nil)
		repo.On("ByID", int64(2)).Return(old, nil)
		repo.On("Update", &current).Return(nil)
		repo.On("Update", &old).Return(nil)

		processor.Process("season", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert season", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		requester := new(mock.SeasonRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSeasonProcessor(repo, requester, logger)

		done := make(chan bool)

		current := newSeason(8, true)
		old := newSeason(2, false)

		seasons := make([]*app.Season, 2)
		seasons[0] = current
		seasons[1] = old

		ch :=seasonChannel(seasons)

		requester.On("Seasons").Return(ch)

		repo.On("ByID", int64(8)).Return(&app.Season{}, errors.New("not found"))
		repo.On("ByID", int64(2)).Return(&app.Season{}, errors.New("not found"))
		repo.On("Insert", current).Return(errors.New("error occurred"))
		repo.On("Insert", old).Return(nil)

		processor.Process("season", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update season", func(t *testing.T) {
		t.Helper()

		repo := new(mock.SeasonRepository)
		requester := new(mock.SeasonRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewSeasonProcessor(repo, requester, logger)

		done := make(chan bool)

		current := newSeason(8, true)
		old := newSeason(2, false)

		seasons := make([]*app.Season, 2)
		seasons[0] = current
		seasons[1] = old

		ch :=seasonChannel(seasons)

		requester.On("Seasons").Return(ch)

		repo.On("ByID", int64(8)).Return(current, nil)
		repo.On("ByID", int64(2)).Return(old, nil)
		repo.On("Update", &current).Return(errors.New("error occurred"))
		repo.On("Update", &old).Return(nil)

		processor.Process("season", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
}

func newSeason(id int64, current bool) *app.Season {
	return &app.Season{
		ID:        id,
		Name:      "2018-2019",
		CompetitionID:  int64(560),
		IsCurrent: current,
		CreatedAt: time.Unix(1546965200, 0),
		UpdatedAt: time.Unix(1546965200, 0),
	}
}

func seasonChannel(seasons []*app.Season) chan *app.Season {
	ch := make(chan *app.Season, len(seasons))

	for _, c := range seasons {
		ch <- c
	}

	close(ch)

	return ch
}
