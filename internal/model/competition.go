package model

import "time"

type Competition struct {
	ID   		 int       `json:"id"`
	Name  		 string    `json:"name"`
	CountryID    int       `json:"country_id"`
	IsCup        bool      `json:"is_cup"`
	CreatedAt  	 time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}