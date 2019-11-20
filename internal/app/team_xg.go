package app

import "time"

type FixtureTeamXG struct {
	ID        uint64    `json:"id"`
	FixtureID uint64    `json:"fixture_id"`
	Home      float32   `json:"home"`
	Away      float32   `json:"away"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
