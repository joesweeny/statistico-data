package team_stats

import "github.com/statistico/statistico-data/internal/model"

type TeamRepository interface {
	InsertTeamStats(m *model.TeamStats) error
	UpdateTeamStats(m *model.TeamStats) error
	ByFixtureAndTeam(fixtureId, teamId uint64) (*model.TeamStats, error)
}
