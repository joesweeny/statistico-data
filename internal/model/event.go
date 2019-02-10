package model

import "time"

type GoalEvent struct {
	ID             int       `json:"id"`
	FixtureID      int       `json:"fixture_id"`
	TeamID         int       `json:"team_id"`
	PlayerID       int       `json:"player_id"`
	PlayerAssistID *int      `json:"player_assist_id"`
	Minute         int       `json:"minute"`
	Score          string    `json:"score"`
	CreatedAt      time.Time `json:"created_at"`
}

type SubstitutionEvent struct {
	ID          int       `json:"id"`
	FixtureID   int       `json:"fixture_id"`
	TeamID      int       `json:"team_id"`
	PlayerInID  int       `json:"player_in_id"`
	PlayerOutID int       `json:"player_out_id"`
	Minute      int       `json:"minute"`
	Injured     *bool     `json:"injured"`
	CreatedAt   time.Time `json:"created_at"`
}
