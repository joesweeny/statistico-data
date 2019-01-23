package model

import "time"

type Team struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	ShortCode     *string     `json:"short_code"`
	CountryID     *int        `json:"country_id"`
	VenueID       int         `json:"venue_id"`
	NationalTeam  bool        `json:"national_team"`
	Founded       *int	      `json:"founded"`
	Logo          *string     `json:"logo"`
	CreatedAt  	  time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}
