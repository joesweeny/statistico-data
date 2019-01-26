package model

import "time"

type Squad struct {
	SeasonID int
	TeamID int
	PlayerIDs []int
	CreatedAt time.Time
	UpdatedAt time.Time
}
