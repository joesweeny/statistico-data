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

// SquadRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type SquadRequester interface {
	SquadsBySeasonIDs(seasonIDs []int64) <-chan *Squad
}
