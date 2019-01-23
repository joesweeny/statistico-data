package model

import "time"

type Round struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	SeasonID  int       `json:"season_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
