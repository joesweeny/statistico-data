package model

import "time"

type TeamStats struct {
	FixtureID     int `json:"fixture_id"`
	TeamID        int `json:"team_id"`
	TeamShots     `json:"shots"`
	TeamPasses    `json:"passes"`
	TeamAttacks   `json:"attacks"`
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

type TeamShots struct {
	Total      *int `json:"total"`
	OnGoal     *int `json:"on_goal"`
	OffGoal    *int `json:"off_goal"`
	Blocked    *int `json:"blocked"`
	InsideBox  *int `json:"inside_box"`
	OutsideBox *int `json:"outside_box"`
}

type TeamPasses struct {
	Total      *int `json:"total"`
	Accuracy   *int `json:"accuracy"`
	Percentage *int `json:"percentage"`
}

type TeamAttacks struct {
	Total     *int `json:"total"`
	Dangerous *int `json:"dangerous"`
}
