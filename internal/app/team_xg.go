package app

import "time"

type FixtureTeamXG struct {
	FixtureID uint64    `json:"fixture_id"`
	Home      float32   `json:"home"`
	Away      float32   `json:"away"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
