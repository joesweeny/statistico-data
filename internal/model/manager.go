package model

import "time"

type Manager struct {
	ID          int       `json:"id"`
	TeamID      *int      `json:"team_id"`
	CountryID   int       `json:"country_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Nationality string    `json:"nationality"`
	Image       *string   `json:"image_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
