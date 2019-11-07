package app

import "time"

// PlayerStats domain entity.
type PlayerStats struct {
	FixtureID         uint64     `json:"fixture_id"`
	PlayerID          uint64     `json:"player_id"`
	TeamID            uint64     `json:"team_id"`
	Position          *string `json:"position"`
	FormationPosition *int    `json:"formation_position"`
	IsSubstitute      bool    `json:"is_substitute"`
	PlayerShots       `json:"shots"`
	PlayerGoals       `json:"goals"`
	PlayerFouls       `json:"fouls"`
	YellowCards       *int `json:"yellow_cards"`
	RedCard           *int `json:"red_card"`
	PlayerPenalties   `json:"penalties"`
	PlayerCrosses     `json:"crosses"`
	PlayerPasses      `json:"passes"`
	Assists           *int      `json:"assists"`
	Offsides          *int      `json:"offsides"`
	Saves             *int      `json:"saves"`
	HitWoodwork       *int      `json:"hit_woodwork"`
	Tackles           *int      `json:"tackles"`
	Blocks            *int      `json:"blocks"`
	Interceptions     *int      `json:"interceptions"`
	Clearances        *int      `json:"clearances"`
	MinutesPlayed     *int      `json:"minutes_played"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// PlayerShots domain sub entity.
type PlayerShots struct {
	Total  *int `json:"total"`
	OnGoal *int `json:"on_goal"`
}

// PlayerGoals domain sub entity.
type PlayerGoals struct {
	Scored   *int `json:"scored"`
	Conceded *int `json:"conceded"`
}

// PlayerFouls domain sub entity.
type PlayerFouls struct {
	Drawn     *int `json:"drawn"`
	Committed *int `json:"committed"`
}

// PlayerCrosses domain sub entity.
type PlayerCrosses struct {
	Total    *int `json:"total"`
	Accuracy *int `json:"accuracy"`
}

// PlayerPasses domain sub entity.
type PlayerPasses struct {
	Total    *int `json:"total"`
	Accuracy *int `json:"accuracy"`
}

// PlayerPenalties domain sub entity.
type PlayerPenalties struct {
	Scored    *int `json:"scored"`
	Missed    *int `json:"missed"`
	Saved     *int `json:"saved"`
	Committed *int `json:"committed"`
	Won       *int `json:"won"`
}
