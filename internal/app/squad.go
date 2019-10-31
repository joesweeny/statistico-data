package app

import (
	"time"
)

// Squad domain entity.
type Squad struct {
	SeasonID  int64
	TeamID    int64
	PlayerIDs []int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SquadRepository provides an interface to persist Squad domain struct objects to a storage engine.
type SquadRepository interface {
	Insert(m *Squad) error
	Update(m *Squad) error
	BySeasonAndTeam(seasonId, teamId int64) (*Squad, error)
	All() ([]Squad, error)
	CurrentSeason() ([]Squad, error)
}
