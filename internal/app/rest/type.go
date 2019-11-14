package rest

type Date struct {
	UTC uint64 `json:"utc"`
	RFC string `json:"rfc"`
}

type Fixture struct {
	ID         uint64 `json:"id"`
	HomeTeamID Team   `json:"home_team"`
	AwayTeamID Team   `json:"away_team"`
	RoundID    Round  `json:"round"`
	VenueID    Venue  `json:"venue"`
	Date       Date   `json:"date"`
}

type Round struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Team struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Venue struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
