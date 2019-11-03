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

func TestTeamProcessor_Process(t *testing.T) {
	t.Run("inserts new team into repository when processing team command", func(t *testing.T) {
		t.Helper()

		teamRepo := new(mock.TeamRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamProcessor(teamRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newTeam(1, "West Ham United")
		ncu := newTeam(14, "Newcastle United")

		teams := make([]*app.Team, 2)
		teams[0] = whu
		teams[1] = ncu

		ch := teamChannel(teams)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("TeamsBySeasonIDs", ids).Return(ch)

		teamRepo.On("ByID", uint64(1)).Return(&app.Team{}, errors.New("not Found"))
		teamRepo.On("ByID", uint64(14)).Return(&app.Team{}, errors.New("not Found"))
		teamRepo.On("Insert", whu).Return(nil)
		teamRepo.On("Insert", ncu).Return(nil)

		processor.Process("team", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		teamRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing team into repository when processing team command", func(t *testing.T) {
		t.Helper()

		teamRepo := new(mock.TeamRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamProcessor(teamRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newTeam(1, "West Ham United")
		ncu := newTeam(14, "Newcastle United")

		teams := make([]*app.Team, 2)
		teams[0] = whu
		teams[1] = ncu

		ch := teamChannel(teams)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("TeamsBySeasonIDs", ids).Return(ch)

		teamRepo.On("ByID", uint64(1)).Return(whu, nil)
		teamRepo.On("ByID", uint64(14)).Return(ncu, nil)
		teamRepo.On("Update", &whu).Return(nil)
		teamRepo.On("Update", &ncu).Return(nil)

		processor.Process("team", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		teamRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert team into repository when processing team command", func(t *testing.T) {
		t.Helper()

		teamRepo := new(mock.TeamRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamProcessor(teamRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newTeam(1, "West Ham United")
		ncu := newTeam(14, "Newcastle United")

		teams := make([]*app.Team, 2)
		teams[0] = whu
		teams[1] = ncu

		ch := teamChannel(teams)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("TeamsBySeasonIDs", ids).Return(ch)

		teamRepo.On("ByID", uint64(1)).Return(&app.Team{}, errors.New("not Found"))
		teamRepo.On("ByID", uint64(14)).Return(&app.Team{}, errors.New("not Found"))
		teamRepo.On("Insert", whu).Return(errors.New("error occurred"))
		teamRepo.On("Insert", ncu).Return(nil)

		processor.Process("team", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		teamRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update team into repository when processing team command", func(t *testing.T) {
		t.Helper()

		t.Helper()

		teamRepo := new(mock.TeamRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamProcessor(teamRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newTeam(1, "West Ham United")
		ncu := newTeam(14, "Newcastle United")

		teams := make([]*app.Team, 2)
		teams[0] = whu
		teams[1] = ncu

		ch := teamChannel(teams)

		ids := []uint64{45, 51}

		seasonRepo.On("IDs").Return(ids, nil)

		requester.On("TeamsBySeasonIDs", ids).Return(ch)

		teamRepo.On("ByID", uint64(1)).Return(whu, nil)
		teamRepo.On("ByID", uint64(14)).Return(ncu, nil)
		teamRepo.On("Update", &whu).Return(errors.New("error occurred"))
		teamRepo.On("Update", &ncu).Return(nil)

		processor.Process("team", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		teamRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new team into repository when processing team current season command", func(t *testing.T) {
		t.Helper()

		teamRepo := new(mock.TeamRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamProcessor(teamRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newTeam(1, "West Ham United")
		ncu := newTeam(14, "Newcastle United")

		teams := make([]*app.Team, 2)
		teams[0] = whu
		teams[1] = ncu

		ch := teamChannel(teams)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("TeamsBySeasonIDs", ids).Return(ch)

		teamRepo.On("ByID", uint64(1)).Return(&app.Team{}, errors.New("not Found"))
		teamRepo.On("ByID", uint64(14)).Return(&app.Team{}, errors.New("not Found"))
		teamRepo.On("Insert", whu).Return(nil)
		teamRepo.On("Insert", ncu).Return(nil)

		processor.Process("team:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		teamRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing team into repository when processing team current season command", func(t *testing.T) {
		t.Helper()

		teamRepo := new(mock.TeamRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamProcessor(teamRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newTeam(1, "West Ham United")
		ncu := newTeam(14, "Newcastle United")

		teams := make([]*app.Team, 2)
		teams[0] = whu
		teams[1] = ncu

		ch := teamChannel(teams)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("TeamsBySeasonIDs", ids).Return(ch)

		teamRepo.On("ByID", uint64(1)).Return(whu, nil)
		teamRepo.On("ByID", uint64(14)).Return(ncu, nil)
		teamRepo.On("Update", &whu).Return(nil)
		teamRepo.On("Update", &ncu).Return(nil)

		processor.Process("team:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		teamRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert team into repository when processing team current season command", func(t *testing.T) {
		t.Helper()

		t.Helper()

		teamRepo := new(mock.TeamRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamProcessor(teamRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newTeam(1, "West Ham United")
		ncu := newTeam(14, "Newcastle United")

		teams := make([]*app.Team, 2)
		teams[0] = whu
		teams[1] = ncu

		ch := teamChannel(teams)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("TeamsBySeasonIDs", ids).Return(ch)

		teamRepo.On("ByID", uint64(1)).Return(&app.Team{}, errors.New("not Found"))
		teamRepo.On("ByID", uint64(14)).Return(&app.Team{}, errors.New("not Found"))
		teamRepo.On("Insert", whu).Return(errors.New("error occurred"))
		teamRepo.On("Insert", ncu).Return(nil)

		processor.Process("team:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		teamRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update team into repository when processing team current season command", func(t *testing.T) {
		t.Helper()

		t.Helper()

		teamRepo := new(mock.TeamRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamProcessor(teamRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		whu := newTeam(1, "West Ham United")
		ncu := newTeam(14, "Newcastle United")

		teams := make([]*app.Team, 2)
		teams[0] = whu
		teams[1] = ncu

		ch := teamChannel(teams)

		ids := []uint64{45, 51}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("TeamsBySeasonIDs", ids).Return(ch)

		teamRepo.On("ByID", uint64(1)).Return(whu, nil)
		teamRepo.On("ByID", uint64(14)).Return(ncu, nil)
		teamRepo.On("Update", &whu).Return(errors.New("error occurred"))
		teamRepo.On("Update", &ncu).Return(nil)

		processor.Process("team:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		teamRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
}

func newTeam(id uint64, name string) *app.Team {
	return &app.Team{
		ID:           id,
		Name:         name,
		VenueID:      560,
		CountryID:    uint64(462),
		NationalTeam: false,
		CreatedAt:    time.Unix(1546965200, 0),
		UpdatedAt:    time.Unix(1546965200, 0),
	}
}

func teamChannel(teams []*app.Team) chan *app.Team {
	ch := make(chan *app.Team, len(teams))

	for _, c := range teams {
		ch <- c
	}

	close(ch)

	return ch
}
