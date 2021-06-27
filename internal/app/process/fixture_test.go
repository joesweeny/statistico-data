package process_test

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-football-data/internal/app"
	"github.com/statistico/statistico-football-data/internal/app/mock"
	"github.com/statistico/statistico-football-data/internal/app/process"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFixtureProcessor_Process(t *testing.T) {
	t.Run("inserts new fixture into repository when processing fixtures by competition id command", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.FixtureRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewFixtureProcessor(fixtureRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newFixture(34)
		two := newFixture(400)

		fixtures := make([]app.Fixture, 2)
		fixtures[0] = one
		fixtures[1] = two

		ch := fixtureChannel(fixtures)

		seasonRepo.On("ByCompetitionId", uint64(5), "name_asc").Return(
			[]app.Season{
				*newSeason(1, false),
				*newSeason(2, true),
			},
			nil,
		)

		requester.On("FixturesBySeasonIDs", []uint64{1, 2}).Return(ch)

		fixtureRepo.On("ByID", uint64(34)).Return(&app.Fixture{}, errors.New("not Found"))
		fixtureRepo.On("ByID", uint64(400)).Return(&app.Fixture{}, errors.New("not Found"))
		fixtureRepo.On("Insert", &one).Return(nil)
		fixtureRepo.On("Insert", &two).Return(nil)

		processor.Process("fixtures:by-competition-id", "5", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing fixture into repository when processing fixtures by competition id command", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.FixtureRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewFixtureProcessor(fixtureRepo, seasonRepo, requester, logger)

		done := make(chan bool)
	
		one := newFixture(34)
		two := newFixture(400)

		fixtures := make([]app.Fixture, 2)
		fixtures[0] = one
		fixtures[1] = two

		ch := fixtureChannel(fixtures)

		seasonRepo.On("ByCompetitionId", uint64(5), "name_asc").Return(
			[]app.Season{
				*newSeason(1, false),
				*newSeason(2, true),
			},
			nil,
		)

		requester.On("FixturesBySeasonIDs", []uint64{1, 2}).Return(ch)

		fixtureRepo.On("ByID", uint64(34)).Return(&app.Fixture{}, nil)
		fixtureRepo.On("ByID", uint64(400)).Return(&app.Fixture{}, nil)
		fixtureRepo.On("Update", &one).Return(nil)
		fixtureRepo.On("Update", &two).Return(nil)

		processor.Process("fixtures:by-competition-id", "5", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert fixture into repository when processing fixtures by competition id command", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.FixtureRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewFixtureProcessor(fixtureRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newFixture(34)
		two := newFixture(400)

		fixtures := make([]app.Fixture, 2)
		fixtures[0] = one
		fixtures[1] = two

		ch := fixtureChannel(fixtures)

		seasonRepo.On("ByCompetitionId", uint64(5), "name_asc").Return(
			[]app.Season{
				*newSeason(1, false),
				*newSeason(2, true),
			},
			nil,
		)

		requester.On("FixturesBySeasonIDs", []uint64{1, 2}).Return(ch)

		fixtureRepo.On("ByID", uint64(34)).Return(&app.Fixture{}, errors.New("not Found"))
		fixtureRepo.On("ByID", uint64(400)).Return(&app.Fixture{}, errors.New("not Found"))
		fixtureRepo.On("Insert", &one).Return(errors.New("error occurred"))
		fixtureRepo.On("Insert", &two).Return(nil)

		processor.Process("fixtures:by-competition-id", "5", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update fixture into repository when processing fixtures by competition id command", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.FixtureRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewFixtureProcessor(fixtureRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newFixture(34)
		two := newFixture(400)

		fixtures := make([]app.Fixture, 2)
		fixtures[0] = one
		fixtures[1] = two

		ch := fixtureChannel(fixtures)

		seasonRepo.On("ByCompetitionId", uint64(5), "name_asc").Return(
			[]app.Season{
				*newSeason(1, false),
				*newSeason(2, true),
			},
			nil,
		)

		requester.On("FixturesBySeasonIDs", []uint64{1, 2}).Return(ch)

		fixtureRepo.On("ByID", uint64(34)).Return(&app.Fixture{}, nil)
		fixtureRepo.On("ByID", uint64(400)).Return(&app.Fixture{}, nil)
		fixtureRepo.On("Update", &one).Return(errors.New("error occurred"))
		fixtureRepo.On("Update", &two).Return(nil)

		processor.Process("fixtures:by-competition-id", "5", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new fixture into repository when processing fixture current season command", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.FixtureRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewFixtureProcessor(fixtureRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newFixture(34)
		two := newFixture(400)

		fixtures := make([]app.Fixture, 2)
		fixtures[0] = one
		fixtures[1] = two

		ch := fixtureChannel(fixtures)

		ids := []uint64{34, 400}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("FixturesBySeasonIDs", ids).Return(ch)

		fixtureRepo.On("ByID", uint64(34)).Return(&app.Fixture{}, errors.New("not Found"))
		fixtureRepo.On("ByID", uint64(400)).Return(&app.Fixture{}, errors.New("not Found"))
		fixtureRepo.On("Insert", &one).Return(nil)
		fixtureRepo.On("Insert", &two).Return(nil)

		processor.Process("fixtures:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing fixture into repository when processing fixture current season command", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.FixtureRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewFixtureProcessor(fixtureRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newFixture(34)
		two := newFixture(400)

		fixtures := make([]app.Fixture, 2)
		fixtures[0] = one
		fixtures[1] = two

		ch := fixtureChannel(fixtures)

		ids := []uint64{34, 400}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("FixturesBySeasonIDs", ids).Return(ch)

		fixtureRepo.On("ByID", uint64(34)).Return(&app.Fixture{}, nil)
		fixtureRepo.On("ByID", uint64(400)).Return(&app.Fixture{}, nil)
		fixtureRepo.On("Update", &one).Return(nil)
		fixtureRepo.On("Update", &two).Return(nil)

		processor.Process("fixtures:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("logs error when unable to insert fixture into repository when processing fixture current season command", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.FixtureRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewFixtureProcessor(fixtureRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newFixture(34)
		two := newFixture(400)

		fixtures := make([]app.Fixture, 2)
		fixtures[0] = one
		fixtures[1] = two

		ch := fixtureChannel(fixtures)

		ids := []uint64{34, 400}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("FixturesBySeasonIDs", ids).Return(ch)

		fixtureRepo.On("ByID", uint64(34)).Return(&app.Fixture{}, errors.New("not Found"))
		fixtureRepo.On("ByID", uint64(400)).Return(&app.Fixture{}, errors.New("not Found"))
		fixtureRepo.On("Insert", &one).Return(errors.New("error occurred"))
		fixtureRepo.On("Insert", &two).Return(nil)

		processor.Process("fixtures:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("logs error when unable to update fixture into repository when processing fixture current season command", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.FixtureRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewFixtureProcessor(fixtureRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newFixture(34)
		two := newFixture(400)

		fixtures := make([]app.Fixture, 2)
		fixtures[0] = one
		fixtures[1] = two

		ch := fixtureChannel(fixtures)

		ids := []uint64{34, 400}

		seasonRepo.On("CurrentSeasonIDs").Return(ids, nil)

		requester.On("FixturesBySeasonIDs", ids).Return(ch)

		fixtureRepo.On("ByID", uint64(34)).Return(&app.Fixture{}, nil)
		fixtureRepo.On("ByID", uint64(400)).Return(&app.Fixture{}, nil)
		fixtureRepo.On("Update", &one).Return(errors.New("error occurred"))
		fixtureRepo.On("Update", &two).Return(nil)

		processor.Process("fixtures:current-season", "", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("delete fixture from repository if fixture status is Deleted", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.FixtureRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewFixtureProcessor(fixtureRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newFixture(34)
		two := newFixture(400)

		status := "Deleted"

		two.Status = &status

		fixtures := make([]app.Fixture, 2)
		fixtures[0] = one
		fixtures[1] = two

		ch := fixtureChannel(fixtures)

		seasonRepo.On("ByCompetitionId", uint64(5), "name_asc").Return(
			[]app.Season{
				*newSeason(1, false),
				*newSeason(2, true),
			},
			nil,
		)

		requester.On("FixturesBySeasonIDs", []uint64{1, 2}).Return(ch)

		fixtureRepo.On("ByID", uint64(34)).Return(&app.Fixture{}, errors.New("not Found"))
		fixtureRepo.On("Delete", uint64(400)).Return(nil)
		fixtureRepo.On("Insert", &one).Return(nil)

		processor.Process("fixtures:by-competition-id", "5", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("delete fixture from repository if fixture status is POSTP", func(t *testing.T) {
		t.Helper()

		fixtureRepo := new(mock.FixtureRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.FixtureRequester)
		logger, hook := test.NewNullLogger()

		processor := process.NewFixtureProcessor(fixtureRepo, seasonRepo, requester, logger)

		done := make(chan bool)

		one := newFixture(34)
		two := newFixture(400)

		status := "POSTP"

		two.Status = &status

		fixtures := make([]app.Fixture, 2)
		fixtures[0] = one
		fixtures[1] = two

		ch := fixtureChannel(fixtures)

		seasonRepo.On("ByCompetitionId", uint64(5), "name_asc").Return(
			[]app.Season{
				*newSeason(1, false),
				*newSeason(2, true),
			},
			nil,
		)

		requester.On("FixturesBySeasonIDs", []uint64{1, 2}).Return(ch)

		fixtureRepo.On("ByID", uint64(34)).Return(&app.Fixture{}, errors.New("not Found"))
		fixtureRepo.On("Delete", uint64(400)).Return(nil)
		fixtureRepo.On("Insert", &one).Return(nil)

		processor.Process("fixtures:by-competition-id", "5", done)

		<-done

		requester.AssertExpectations(t)
		seasonRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})
}

func newFixture(id uint64) app.Fixture {
	var roundId = uint64(165789)

	return app.Fixture{
		ID:         id,
		SeasonID:   uint64(14567),
		RoundID:    &roundId,
		HomeTeamID: 451,
		AwayTeamID: 924,
		Date:       time.Unix(1548086929, 0),
		CreatedAt:  time.Unix(1546965200, 0),
		UpdatedAt:  time.Unix(1546965200, 0),
	}
}

func fixtureChannel(fixtures []app.Fixture) chan app.Fixture {
	ch := make(chan app.Fixture, len(fixtures))

	for _, c := range fixtures {
		ch <- c
	}

	close(ch)

	return ch
}
