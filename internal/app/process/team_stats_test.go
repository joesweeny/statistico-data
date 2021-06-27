package process_test

import (
	"errors"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-football-data/internal/app"
	"github.com/statistico/statistico-football-data/internal/app/mock"
	"github.com/statistico/statistico-football-data/internal/app/process"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTeamStatsProcessor_Process(t *testing.T) {
	t.Run("inserts new team stats into repository when processing team stats by season id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		requester.On("TeamStatsBySeasonIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", &home).Return(nil)
		teamStatsRepo.On("InsertTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-season-id", "45", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("log errors if unable to insert team stats into repository when processing team stats by season id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		requester.On("TeamStatsBySeasonIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", &home).Return(errors.New("error occurred"))
		teamStatsRepo.On("InsertTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-season-id", "45", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})

	t.Run("update existing team stats into repository when processing team stats by season id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		requester.On("TeamStatsBySeasonIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("UpdateTeamStats", &home).Return(nil)
		teamStatsRepo.On("UpdateTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-season-id", "45", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("log errors if unable to update team stats into repository when processing team stats by season id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		requester.On("TeamStatsBySeasonIDs", []uint64{45}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("UpdateTeamStats", &home).Return(nil)
		teamStatsRepo.On("UpdateTeamStats", &away).Return(errors.New("error occurred"))

		processor.Process("team-stats:by-season-id", "45", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new team stats into repository when processing team stats by date command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		competitionRepo.On("IDs").Return([]uint64{1, 2, 3}, nil)

		date, _ := time.Parse("2006-01-02", "2021-01-18")

		requester.On("TeamStatsByDate", date, []uint64{1, 2, 3}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", &home).Return(nil)
		teamStatsRepo.On("InsertTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-date", "2021-01-18", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("log errors if unable to insert team stats into repository when processing team stats by date command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		competitionRepo.On("IDs").Return([]uint64{1, 2, 3}, nil)

		date, _ := time.Parse("2006-01-02", "2021-01-18")

		requester.On("TeamStatsByDate", date, []uint64{1, 2, 3}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", &home).Return(errors.New("error occurred"))
		teamStatsRepo.On("InsertTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-date", "2021-01-18", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		competitionRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})

	t.Run("updates existing team stats into repository when processing team stats by date command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		competitionRepo.On("IDs").Return([]uint64{1, 2, 3}, nil)

		date, _ := time.Parse("2006-01-02", "2021-01-18")

		requester.On("TeamStatsByDate", date, []uint64{1, 2, 3}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("UpdateTeamStats", &home).Return(nil)
		teamStatsRepo.On("UpdateTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-date", "2021-01-18", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("log errors if unable to update team stats into repository when processing team stats by date command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		competitionRepo.On("IDs").Return([]uint64{1, 2, 3}, nil)

		date, _ := time.Parse("2006-01-02", "2021-01-18")

		requester.On("TeamStatsByDate", date, []uint64{1, 2, 3}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("UpdateTeamStats", &home).Return(errors.New("error occurred"))
		teamStatsRepo.On("UpdateTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-date", "2021-01-18", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		competitionRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})

	t.Run("inserts new team stats into repository when processing team stats by competition id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		seasonRepo.On("ByCompetitionId", uint64(5), "name_asc").Return(
			[]app.Season{
				*newSeason(1, false),
				*newSeason(2, true),
			},
			nil,
		)

		requester.On("TeamStatsBySeasonIDs", []uint64{1, 2}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", &home).Return(nil)
		teamStatsRepo.On("InsertTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-competition-id", "5", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("log errors if unable to insert team stats into repository when processing team stats by competition id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		seasonRepo.On("ByCompetitionId", uint64(5), "name_asc").Return(
			[]app.Season{
				*newSeason(1, false),
				*newSeason(2, true),
			},
			nil,
		)

		requester.On("TeamStatsBySeasonIDs", []uint64{1, 2}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, errors.New("not found"))
		teamStatsRepo.On("InsertTeamStats", &home).Return(errors.New("error occurred"))
		teamStatsRepo.On("InsertTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-competition-id", "5", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		competitionRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})

	t.Run("updates existing team stats into repository when processing team stats by competition id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		seasonRepo.On("ByCompetitionId", uint64(5), "name_asc").Return(
			[]app.Season{
				*newSeason(1, false),
				*newSeason(2, true),
			},
			nil,
		)

		requester.On("TeamStatsBySeasonIDs", []uint64{1, 2}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("UpdateTeamStats", &home).Return(nil)
		teamStatsRepo.On("UpdateTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-competition-id", "5", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("log errors if unable to update team stats into repository when processing team stats by competition id command", func(t *testing.T) {
		t.Helper()

		teamStatsRepo := new(mock.TeamStatsRepository)
		competitionRepo := new(mock.CompetitionRepository)
		seasonRepo := new(mock.SeasonRepository)
		requester := new(mock.TeamStatsRequester)
		clock := clockwork.NewFakeClock()
		logger, hook := test.NewNullLogger()

		processor := process.NewTeamStatsProcessor(teamStatsRepo, competitionRepo, seasonRepo, requester, clock, logger)

		done := make(chan bool)

		home := newTeamStats(45, 99)
		away := newTeamStats(45, 2)

		stats := make([]app.TeamStats, 2)
		stats[0] = home
		stats[1] = away

		ch := teamStatsChannel(stats)

		seasonRepo.On("ByCompetitionId", uint64(5), "name_asc").Return(
			[]app.Season{
				*newSeason(1, false),
				*newSeason(2, true),
			},
			nil,
		)

		requester.On("TeamStatsBySeasonIDs", []uint64{1, 2}).Return(ch)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(99)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("ByFixtureAndTeam", uint64(45), uint64(2)).Return(&app.TeamStats{}, nil)
		teamStatsRepo.On("UpdateTeamStats", &home).Return(errors.New("error occurred"))
		teamStatsRepo.On("UpdateTeamStats", &away).Return(nil)

		processor.Process("team-stats:by-competition-id", "5", done)

		<-done

		requester.AssertExpectations(t)
		teamStatsRepo.AssertExpectations(t)
		competitionRepo.AssertExpectations(t)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})
}

func newTeamStats(fixtureId, teamId uint64) app.TeamStats {
	return app.TeamStats{
		FixtureID:   fixtureId,
		TeamID:      teamId,
		TeamShots:   app.TeamShots{},
		TeamPasses:  app.TeamPasses{},
		TeamAttacks: app.TeamAttacks{},
		CreatedAt:   time.Unix(1546965200, 0),
		UpdatedAt:   time.Unix(1546965200, 0),
	}
}

func teamStatsChannel(stats []app.TeamStats) chan app.TeamStats {
	ch := make(chan app.TeamStats, len(stats))

	for _, c := range stats {
		ch <- c
	}

	close(ch)

	return ch
}
