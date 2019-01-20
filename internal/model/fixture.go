package model

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type Fixture struct {
	ID           int          `json:"id"`
	SeasonID     int          `json:"season_id"`
	RoundID      null.Int     `json:"round_id"`
	VenueID      null.Int     `json:"venue_id"`
	HomeTeamID   int          `json:"home_team_id"`
	AwayTeamID   int          `json:"away_team_id"`
	RefereeID    null.Int     `json:"referee_id"`
	Date         time.Time    `json:"date"`
	CreatedAt  	 time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
