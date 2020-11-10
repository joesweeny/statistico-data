package app

import (
	"time"
)

// CardEvent domain entity
type CardEvent struct {
	ID          uint64    `json:"id"`
	TeamID      uint64    `json:"team_id"`
	Type        string    `json:"type"`
	FixtureID   uint64    `json:"fixture_id"`
	PlayerID    uint64    `json:"player_id"`
	Minute      uint8     `json:"minute"`
	Reason      *string   `json:"reason"`
	CreatedAt   time.Time `json:"created_at"`
}

// GoalEvent domain entity.
type GoalEvent struct {
	ID             uint64    `json:"id"`
	FixtureID      uint64    `json:"fixture_id"`
	TeamID         uint64    `json:"team_id"`
	PlayerID       uint64    `json:"player_id"`
	PlayerAssistID *uint64   `json:"player_assist_id"`
	Minute         int       `json:"minute"`
	Score          string    `json:"score"`
	CreatedAt      time.Time `json:"created_at"`
}

// SubstitutionEvent domain entity.
type SubstitutionEvent struct {
	ID          uint64    `json:"id"`
	FixtureID   uint64    `json:"fixture_id"`
	TeamID      uint64    `json:"team_id"`
	PlayerInID  uint64    `json:"player_in_id"`
	PlayerOutID uint64    `json:"player_out_id"`
	Minute      int       `json:"minute"`
	Injured     *bool     `json:"injured"`
	CreatedAt   time.Time `json:"created_at"`
}

// EventRepository provides an interface to persist event domain struct objects to a storage engine.
type EventRepository interface {
	InsertCardEvent(e *CardEvent) error
	InsertGoalEvent(e *GoalEvent) error
	InsertSubstitutionEvent(e *SubstitutionEvent) error
	CardEventsForFixture(fixtureID uint64) ([]*CardEvent, error)
	GoalEventsForFixture(fixtureID uint64) ([]*GoalEvent, error)
	GoalEventByID(id uint64) (*GoalEvent, error)
	SubstitutionEventByID(id uint64) (*SubstitutionEvent, error)
}

// EventRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type EventRequester interface {
	EventsByFixtureIDs(ids []uint64) (<-chan *GoalEvent, <-chan *SubstitutionEvent, <-chan *CardEvent)
}
