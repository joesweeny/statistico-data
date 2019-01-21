package model

import (
	"time"
)

type Fixture struct {
	ID           int          `json:"id"`
	SeasonID     int          `json:"season_id"`
	RoundID      *int         `json:"round_id"`
	VenueID      *int         `json:"venue_id"`
	HomeTeamID   int          `json:"home_team_id"`
	AwayTeamID   int          `json:"away_team_id"`
	RefereeID    *int         `json:"referee_id"`
	Date         time.Time    `json:"date"`
	CreatedAt  	 time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
