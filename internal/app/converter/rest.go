package converter

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/rest"
	"time"
)

// Convert a domain Team struct into a rest Team struct
func TeamToRest(t *app.Team) rest.Team {
	var x rest.Team
	x.ID = t.ID
	x.Name = t.Name

	return x
}

// Convert a domain Round struct into a rest Round struct
func RoundToRest(r *app.Round) rest.Round {
	var x rest.Round
	x.ID =        r.ID
	x.Name =      r.Name
	x.SeasonID =  r.SeasonID
	x.StartDate = rest.Date{
		UTC: r.StartDate.Unix(),
		RFC: r.StartDate.Format(time.RFC3339),
	}
	x.EndDate =   rest.Date{
		UTC: r.EndDate.Unix(),
		RFC: r.EndDate.Format(time.RFC3339),
	}

	return x
}

// Convert a domain Venue struct into a rest Venue struct
func VenueToRest(v *app.Venue) rest.Venue {
	ven := rest.Venue{}
	ven.ID = v.ID
	ven.Name = v.Name

	return ven
}
