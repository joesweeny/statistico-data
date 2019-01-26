package model

import "time"

type Player struct {
	ID int `json:"id"`
	CountryId int `json:"country_id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	BirthPlace *string `json:"birth_place"`
	DateOfBirth time.Time `json:"date_of_birth"`
	PositionID int `json:"position_id"`
	Image *string `json:"image_path"`
	CreatedAt time.Time `json:"created_at"`
}
