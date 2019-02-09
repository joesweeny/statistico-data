package model

import "time"

type PlayerStats struct {
	FixtureID         int    `json:"fixture_id"`
	PlayerID          int    `json:"player_id"`
	TeamID            int    `json:"team_id"`
	Position          *string `json:"position"`
	FormationPosition *int   `json:"formation_position"`
	IsSubstitute      bool   `json:"is_substitute"`
	PlayerShots       `json:"shots"`
	PlayerGoals       `json:"goals"`
	PlayerFouls       `json:"fouls"`
	YellowCards       *int `json:"yellow_cards"`
	RedCard           *int `json:"red_card"`
	PlayerPenalties   `json:"penalties"`
	PlayerCrosses     `json:"crosses"`
	PlayerPasses      `json:"passes"`
	Assists           *int `json:"assists"`
	Offsides          *int `json:"offsides"`
	Saves             *int `json:"saves"`
	HitWoodwork       *int      `json:"hit_woodwork"`
	Tackles           *int      `json:"tackles"`
	Blocks            *int      `json:"blocks"`
	Interceptions     *int      `json:"interceptions"`
	Clearances        *int      `json:"clearances"`
	MinutesPlayed     *int      `json:"minutes_played"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type PlayerShots struct {
	Total  *int `json:"total"`
	OnGoal *int `json:"on_goal"`
}

type PlayerGoals struct {
	Scored   *int `json:"scored"`
	Conceded *int `json:"conceded"`
}

type PlayerFouls struct {
	Drawn     *int `json:"drawn"`
	Committed *int `json:"committed"`
}

type PlayerCrosses struct {
	Total    *int `json:"total"`
	Accuracy *int `json:"accuracy"`
}

type PlayerPasses struct {
	Total    *int `json:"total"`
	Accuracy *int `json:"accuracy"`
}

type PlayerPenalties struct {
	Scored    *int `json:"scored"`
	Missed    *int `json:"missed"`
	Saved     *int `json:"saved"`
	Committed *int `json:"committed"`
	Won       *int `json:"won"`
}
