package factory

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-proto/data/go"
)

type PlayerStatsFactory struct {
	repo   app.PlayerStatsRepository
	logger *logrus.Logger
}

func (p PlayerStatsFactory) BuildPlayerStats(f *app.Fixture, teamID uint64) ([]*statisticoproto.PlayerStats, error) {
	pl, err := p.repo.ByFixtureAndTeam(f.ID, teamID)

	if err != nil {
		return nil, p.returnLoggedError(f.ID, err)
	}

	return handlePlayerStats(pl), nil
}

func (p PlayerStatsFactory) BuildLineup(f *app.Fixture, teamID uint64) (*statisticoproto.Lineup, error) {
	pl, err := p.repo.ByFixtureAndTeam(f.ID, teamID)

	if err != nil {
		return nil, p.returnLoggedError(f.ID, err)
	}

	lineup := statisticoproto.Lineup{
		Start: handleStartingLineupPlayers(pl),
		Bench: handleSubstituteLineupPlayers(pl),
	}

	return &lineup, nil
}

func (p PlayerStatsFactory) returnLoggedError(id uint64, err error) error {
	p.logger.Warnf("error hydrating proto player stats: fixture %d. error %s", id, err.Error())
	return err
}

func NewPlayerStatsFactory(r app.PlayerStatsRepository, log *logrus.Logger) *PlayerStatsFactory {
	return &PlayerStatsFactory{repo: r, logger: log}
}
