package model

import "time"

type Season struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	LeagueID     int       `json:"league_id"`
	Current      bool      `json:"current"`
	CreatedAt  	 time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
