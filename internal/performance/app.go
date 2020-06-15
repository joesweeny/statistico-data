package performance

type StatFilter struct {
	Stat    string  `json:"stat"`
	Action  string  `json:"action"`
	Team    string  `json:"team"`
	Metric  string  `json:"metric"`
	Measure string  `json:"measure"`
	Value   float32 `json:"value"`
	Venue   string  `json:"venue"`
	Games   uint8   `json:"games"`
}

type Team struct {
	ID   uint64   `json:"id"`
	Name string   `json:"name"`
}

type StatReader interface {
	GetTeams(s *StatFilter) ([]Team, error)
}
