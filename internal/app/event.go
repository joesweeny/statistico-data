package app

import (
	"time"
)

// GoalEvent domain entity.
type GoalEvent struct {
	ID             int64       `json:"id"`
	FixtureID      int64       `json:"fixture_id"`
	TeamID         int64       `json:"team_id"`
	PlayerID       int64       `json:"player_id"`
	PlayerAssistID *int64      `json:"player_assist_id"`
	Minute         int       `json:"minute"`
	Score          string    `json:"score"`
	CreatedAt      time.Time `json:"created_at"`
}

// SubstitutionEvent domain entity.
type SubstitutionEvent struct {
	ID          int64       `json:"id"`
	FixtureID   int64       `json:"fixture_id"`
	TeamID      int64       `json:"team_id"`
	PlayerInID  int64       `json:"player_in_id"`
	PlayerOutID int64       `json:"player_out_id"`
	Minute      int       `json:"minute"`
	Injured     *bool     `json:"injured"`
	CreatedAt   time.Time `json:"created_at"`
}

// EventRepository provides an interface to persist event domain struct objects to a storage engine.
type EventRepository interface {
	InsertGoalEvent(e *GoalEvent) error
	InsertSubstitutionEvent(e *SubstitutionEvent) error
	GoalEventByID(id int64) (*GoalEvent, error)
	SubstitutionEventByID(id int64) (*SubstitutionEvent, error)
}

// EventRequester provides an interface allowing this application to request data from an external
// data provider. The requester implementation is responsible for creating the channel, filtering struct data into
// the channel before closing the channel once successful execution is complete.
type EventRequester interface {
	EventsByFixtureIDs(ids []int64) (<-chan *GoalEvent, <-chan *SubstitutionEvent)
}
