package performance

type StatFilter struct {
	Action  string  `json:"action"`
	Games   uint8   `json:"games"`
	Measure string  `json:"measure"`
	Metric  string  `json:"metric"`
	Seasons []uint64 `json:"seasons"`
	Stat    string  `json:"stat"`
	Value   float32 `json:"value"`
	Venue   string  `json:"venue"`
}

type Team struct {
	ID   uint64   `json:"id"`
	Name string   `json:"name"`
}

type StatReader interface {
	TeamsMatchingFilter(s *StatFilter) ([]*Team, error)
}
