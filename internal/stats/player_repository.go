package stats

import "github.com/joesweeny/statshub/internal/model"

type PlayerRepository interface {
	Insert(m *model.PlayerStats) error
	Update(m *model.PlayerStats) error
	ByFixtureAndTeam(fixtureId, teamId int) ([]model.PlayerStats, error)
	ByFixtureAndPlayer(fixtureId, playerId int) (*model.PlayerStats, error)
}