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
)

func TestCompetitionProcessor_Process(t *testing.T) {
	t.Run("inserts new competition", func(t *testing.T) {
		t.Helper()

		repo := new(mock.CompetitionRepository)
		requester := new(mock.CompetitionRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewCompetitionProcessor(repo, requester, logger)

		done := make(chan bool)

		prem := newCompetition(8, "Premier League")
		cham := newCompetition(16, "Championship")

		competitions := make([]*app.Competition, 2)
		competitions[0] = prem
		competitions[1] = cham

		ch := competitionChannel(competitions)

		requester.On("Competitions").Return(ch)

		repo.On("ByID", uint64(8)).Return(&app.Competition{}, errors.New("not found"))
		repo.On("ByID", uint64(16)).Return(&app.Competition{}, errors.New("not found"))
		repo.On("Insert", prem).Return(nil)
		repo.On("Insert", cham).Return(nil)

		processor.Process("competition", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing competition", func(t *testing.T) {
		t.Helper()

		repo := new(mock.CompetitionRepository)
		requester := new(mock.CompetitionRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewCompetitionProcessor(repo, requester, logger)

		done := make(chan bool)

		prem := newCompetition(8, "Premier League")
		cham := newCompetition(16, "Championship")

		competitions := make([]*app.Competition, 2)
		competitions[0] = prem
		competitions[1] = cham

		ch := competitionChannel(competitions)

		requester.On("Competitions").Return(ch)

		repo.On("ByID", uint64(8)).Return(prem, nil)
		repo.On("ByID", uint64(16)).Return(&app.Competition{}, errors.New("not found"))
		repo.On("Update", &prem).Return(nil)
		repo.On("Insert", cham).Return(nil)

		processor.Process("competition", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert competition", func(t *testing.T) {
		t.Helper()

		repo := new(mock.CompetitionRepository)
		requester := new(mock.CompetitionRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewCompetitionProcessor(repo, requester, logger)

		done := make(chan bool)

		prem := newCompetition(8, "Premier League")
		cham := newCompetition(16, "Championship")

		competitions := make([]*app.Competition, 2)
		competitions[0] = prem
		competitions[1] = cham

		ch := competitionChannel(competitions)

		requester.On("Competitions").Return(ch)

		repo.On("ByID", uint64(8)).Return(&app.Competition{}, errors.New("not found"))
		repo.On("ByID", uint64(16)).Return(&app.Competition{}, errors.New("not found"))
		repo.On("Insert", prem).Return(errors.New("error occurred"))
		repo.On("Insert", cham).Return(nil)

		processor.Process("competition", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update competition", func(t *testing.T) {
		t.Helper()

		repo := new(mock.CompetitionRepository)
		requester := new(mock.CompetitionRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewCompetitionProcessor(repo, requester, logger)

		done := make(chan bool)

		prem := newCompetition(8, "Premier League")
		cham := newCompetition(16, "Championship")

		competitions := make([]*app.Competition, 2)
		competitions[0] = prem
		competitions[1] = cham

		ch := competitionChannel(competitions)

		requester.On("Competitions").Return(ch)

		repo.On("ByID", uint64(8)).Return(prem, nil)
		repo.On("ByID", uint64(16)).Return(cham, nil)
		repo.On("Update", &prem).Return(errors.New("error occurred"))
		repo.On("Update", &cham).Return(nil)

		processor.Process("competition", "", done)

		<-done

		repo.AssertExpectations(t)
		requester.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
}

func newCompetition(id uint64, name string) *app.Competition {
	return &app.Competition{
		ID:        id,
		Name:      name,
		CountryID: uint64(462),
		IsCup:     false,
	}
}

func competitionChannel(competitions []*app.Competition) chan *app.Competition {
	ch := make(chan *app.Competition, len(competitions))

	for _, c := range competitions {
		ch <- c
	}

	close(ch)

	return ch
}
