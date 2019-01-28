package stats

import "github.com/joesweeny/statshub/internal/model"

type TeamRepository interface {
	Insert(m *model.TeamStats) error
	Update(m *model.TeamStats) error
	ByFixtureAndTeam(fixtureId, teamId int) (*model.TeamStats, error)
}
