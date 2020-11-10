package app

import "time"

type FixtureTeamXG struct {
	ID        uint64    `json:"id"`
	FixtureID uint64    `json:"fixture_id"`
	Home      *float32  `json:"home"`
	Away      *float32  `json:"away"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FixtureTeamXGRepository provides an interface to persist FixtureTeamXG domain struct objects to a storage engine.
type FixtureTeamXGRepository interface {
	Insert(f *FixtureTeamXG) error
	Update(f *FixtureTeamXG) error
	ByID(id uint64) (*FixtureTeamXG, error)
	ByFixtureID(id uint64) (*FixtureTeamXG, error)
}
