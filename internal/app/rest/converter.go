package rest


import (
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

// Convert a domain Team struct into a rest Team struct
func convertAppTeam(t *app.Team) Team {
	var x Team
	x.ID = t.ID
	x.Name = t.Name

	return x
}

// Convert a domain Round struct into a rest Round struct
func convertAppRound(r *app.Round) Round {
	var x Round
	x.ID =        r.ID
	x.Name =      r.Name
	x.SeasonID =  r.SeasonID
	x.StartDate = Date{
		UTC: uint64(r.StartDate.Unix()),
		RFC: r.StartDate.Format(time.RFC3339),
	}
	x.EndDate =   Date{
		UTC: uint64(r.EndDate.Unix()),
		RFC: r.EndDate.Format(time.RFC3339),
	}

	return x
}

// Convert a domain Venue struct into a rest Venue struct
func convertAppVenue(v *app.Venue) Venue {
	ven := Venue{}
	ven.ID = v.ID
	ven.Name = v.Name

	return ven
}
