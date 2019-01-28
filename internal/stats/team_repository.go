package stats

import "github.com/joesweeny/statshub/internal/model"

type TeamRepository interface {
	InsertTeamStats(m *model.TeamStats) error
	UpdateTeamStats(m *model.TeamStats) error
	ByFixtureAndTeam(fixtureId, teamId int) (*model.TeamStats, error)
}
