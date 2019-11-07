package process_test

import (
	"errors"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/mock"
	"github.com/statistico/statistico-data/internal/app/process"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTeamStatsProcessor_Process(t *testing.T) {
	t.Run("inserts new team stats into repository when processing team stats by result id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		fixtureRepo.On("ByID", uint64(45)).Return(newFixture(45), nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", home).Return(nil)
		teamStatsRepo.On("InsertTeamStats", away).Return(nil)

		processor.Process("team-stats:by-result-id", "45", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing team stats into repository when processing team stats by result id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		fixtureRepo.On("ByID", uint64(45)).Return(newFixture(45), nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(home, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(away, nil)
		teamStatsRepo.On("UpdateTeamStats", home).Return(nil)
		teamStatsRepo.On("UpdateTeamStats", away).Return(nil)

		processor.Process("team-stats:by-result-id", "45", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("log errors if unable to insert team stats into repository when processing team stats by result id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		fixtureRepo.On("ByID", uint64(45)).Return(newFixture(45), nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", home).Return(errors.New("error occurred"))
		teamStatsRepo.On("InsertTeamStats", away).Return(nil)

		processor.Process("team-stats:by-result-id", "45", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("log errors if unable to update team stats into repository when processing team stats by result id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		fixtureRepo.On("ByID", uint64(45)).Return(newFixture(45), nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(home, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(away, nil)
		teamStatsRepo.On("UpdateTeamStats", home).Return(nil)
		teamStatsRepo.On("UpdateTeamStats", away).Return(errors.New("error occurred"))

		processor.Process("team-stats:by-result-id", "45", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new team stats into repository when processing team stats by season id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		fix := []app.Fixture{*newFixture(45)}

		ch := teamStatsChannel(stats)

		fixtureRepo.On("BySeasonID", uint64(20)).Return(fix, nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", home).Return(nil)
		teamStatsRepo.On("InsertTeamStats", away).Return(nil)

		processor.Process("team-stats:by-season-id", "20", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing team stats into repository when processing team stats by season id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		fix := []app.Fixture{*newFixture(45)}

		ch := teamStatsChannel(stats)

		fixtureRepo.On("BySeasonID", uint64(20)).Return(fix, nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(home, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(away, nil)
		teamStatsRepo.On("UpdateTeamStats", home).Return(nil)
		teamStatsRepo.On("UpdateTeamStats", away).Return(nil)

		processor.Process("team-stats:by-season-id", "20", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("log errors if unable to insert team stats into repository when processing team stats by season id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		fix := []app.Fixture{*newFixture(45)}

		ch := teamStatsChannel(stats)

		fixtureRepo.On("BySeasonID", uint64(20)).Return(fix, nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", home).Return(errors.New("error occurred"))
		teamStatsRepo.On("InsertTeamStats", away).Return(nil)

		processor.Process("team-stats:by-season-id", "20", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("log errors if unable to update team stats into repository when processing team stats by season id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		fix := []app.Fixture{*newFixture(45)}

		ch := teamStatsChannel(stats)

		fixtureRepo.On("BySeasonID", uint64(20)).Return(fix, nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(home, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(away, nil)
		teamStatsRepo.On("UpdateTeamStats", home).Return(nil)
		teamStatsRepo.On("UpdateTeamStats", away).Return(errors.New("error occurred"))

		processor.Process("team-stats:by-season-id", "20", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new team stats into repository when processing team stats today command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		fixtureRepo.On("IDsBetween", from, to).Return([]uint64{34}, nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{34}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", home).Return(nil)
		teamStatsRepo.On("InsertTeamStats", away).Return(nil)

		processor.Process("team-stats:today", "", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("updates existing team stats into repository when processing team stats today command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		fixtureRepo.On("IDsBetween", from, to).Return([]uint64{34}, nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{34}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(home, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(away, nil)
		teamStatsRepo.On("UpdateTeamStats", home).Return(nil)
		teamStatsRepo.On("UpdateTeamStats", away).Return(nil)

		processor.Process("team-stats:today", "", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("log errors if unable to insert team stats into repository when processing team stats today command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		fixtureRepo.On("IDsBetween", from, to).Return([]uint64{34}, nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{34}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", home).Return(errors.New("error occurred"))
		teamStatsRepo.On("InsertTeamStats", away).Return(nil)

		processor.Process("team-stats:today", "", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})

	t.Run("log errors if unable to update team stats into repository when processing team stats today command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		fixtureRepo := new(mock.FixtureRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, fixtureRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]*app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		now := clock.Now()
		y, m, d := now.Date()
		from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
		to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

		fixtureRepo.On("IDsBetween", from, to).Return([]uint64{34}, nil)
		requester.On("TeamStatsByFixtureIDs", []uint64{34}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(99), uint64(45)).Return(home, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(2), uint64(45)).Return(away, nil)
		teamStatsRepo.On("UpdateTeamStats", home).Return(nil)
		teamStatsRepo.On("UpdateTeamStats", away).Return(errors.New("error occurred"))

		processor.Process("team-stats:today", "", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		fixtureRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	})
}

func newTeamStats(fixtureId, teamId uint64) *app.TeamStats {
	return &app.TeamStats{
		FixtureID:   fixtureId,
		TeamID:      teamId,
		TeamShots:   app.TeamShots{},
		TeamPasses:  app.TeamPasses{},
		TeamAttacks: app.TeamAttacks{},
		CreatedAt:   time.Unix(1546965200, 0),
		UpdatedAt:   time.Unix(1546965200, 0),
	}
}

func teamStatsChannel(stats []*app.TeamStats) chan *app.TeamStats {
	ch := make(chan *app.TeamStats, len(stats))

	for _, c := range stats {
		ch <- c
	}

	close(ch)

	return ch
}