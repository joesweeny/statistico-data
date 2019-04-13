package stats

import "github.com/statistico/statistico-data/internal/model"

type PlayerRepository interface {
	InsertPlayerStats(m *model.PlayerStats) error
	UpdatePlayerStats(m *model.PlayerStats) error
	ByFixtureAndPlayer(fixtureId, playerId int) (*model.PlayerStats, error)
	ByFixtureAndTeam(fixtureId, teamId uint64) ([]*model.PlayerStats, error)
}
