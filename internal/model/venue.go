package model

import "time"

type Venue struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Surface   *string   `json:"surface"`
	Address   *string   `json:"address"`
	City      *string   `json:"city"`
	Capacity  *int      `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
