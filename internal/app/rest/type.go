package rest

type Date struct {
	UTC int64 `json:"utc"`
	RFC string `json:"rfc"`
}

type Fixture struct {
	ID         uint64 `json:"id"`
	HomeTeam   Team   `json:"home_team"`
	AwayTeam   Team   `json:"away_team"`
	Round	   Round  `json:"round"`
	Venue	   Venue  `json:"venue"`
	Date       Date   `json:"date"`
}

type Round struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	SeasonID uint64 `json:"season_id"`
	StartDate Date `json:"start_date"`
	EndDate Date `json:"end_date"`
}

type Team struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Venue struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
