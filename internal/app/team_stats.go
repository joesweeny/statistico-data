package app

import "time"

// TeamStats domain entity.
type TeamStats struct {
	FixtureID     uint64 `json:"fixture_id"`
	TeamID        uint64 `json:"team_id"`
	TeamShots     `json:"shots"`
	TeamPasses    `json:"passes"`
	TeamAttacks   `json:"attacks"`
	Goals         *int      `json:"goals"`
	Fouls         *int      `json:"fouls"`
	Corners       *int      `json:"corners"`
	Offsides      *int      `json:"offsides"`
	Possession    *int      `json:"possession"`
	YellowCards   *int      `json:"yellow_cards"`
	RedCards      *int      `json:"red_cards"`
	Saves         *int      `json:"saves"`
	Substitutions *int      `json:"substitutions"`
	GoalKicks     *int      `json:"goal_kicks"`
	GoalAttempts  *int      `json:"goal_attempts"`
	FreeKicks     *int      `json:"free_kicks"`
	ThrowIns      *int      `json:"throw_ins"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TeamStat domain struct for a single stat for a fixture
type TeamStat struct {
	FixtureID     uint64    `json:"fixture_id"`
	Stat          string    `json:"stat"`
	Value         *uint32   `json:"value"`
}

// TeamShots domain sub entity.
type TeamShots struct {
	Total      *int `json:"total"`
	OnGoal     *int `json:"on_goal"`
	OffGoal    *int `json:"off_goal"`
	Blocked    *int `json:"blocked"`
	InsideBox  *int `json:"inside_box"`
	OutsideBox *int `json:"outside_box"`
}

// TeamPasses domain sub entity.
type TeamPasses struct {
	Total      *int `json:"total"`
	Accuracy   *int `json:"accuracy"`
	Percentage *float32 `json:"percentage"`
}

// TeamAttacks domain sub entity.
type TeamAttacks struct {
	Total     *int `json:"total"`
	Dangerous *int `json:"dangerous"`
}

// TeamStatsRepository provides an interface to persist TeamStats domain struct objects to a storage engine.
type TeamStatsRepository interface {
	InsertTeamStats(m *TeamStats) error
	UpdateTeamStats(m *TeamStats) error
	ByFixtureAndTeam(fixtureID, teamID uint64) (*TeamStats, error)
	StatByFixtureAndTeam(stat string, fixtureID, teamID uint64) (*TeamStat, error)
	Get() ([]*TeamStats, error)
}

// TeamStatsRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type TeamStatsRequester interface {
	TeamStatsByFixtureIDs(ids []uint64) <-chan *TeamStats
}
