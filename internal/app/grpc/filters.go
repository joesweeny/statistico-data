package grpc

import (
	"fmt"
	"github.com/statistico/statistico-football-data/internal/app"
	statistico "github.com/statistico/statistico-proto/go"

	"time"
)

func fixtureFilterFromTeamStatRequest(r *statistico.TeamStatRequest) (*app.FixtureFilterQuery, error) {
	var query app.FixtureFilterQuery

	if r.GetDateBefore() != nil {
		date, err := time.Parse(time.RFC3339, r.GetDateBefore().GetValue())

		if err != nil {
			return &query, fmt.Errorf("date provided '%s' is not a valid RFC3339 date", r.GetDateBefore().GetValue())
		}

		query.DateBefore = &date
	}

	if r.GetDateAfter() != nil {
		date, err := time.Parse(time.RFC3339, r.GetDateAfter().GetValue())

		if err != nil {
			return &query, fmt.Errorf("date provided '%s' is not a valid RFC3339 date", r.GetDateAfter().GetValue())
		}

		query.DateAfter = &date
	}

	if r.GetLimit() != nil {
		v := r.GetLimit().GetValue()
		query.Limit = &v
	}

	if len(r.GetSeasonIds()) > 0 {
		query.SeasonIDs = r.GetSeasonIds()
	}

	if r.GetSort() != nil {
		v := r.GetSort().GetValue()
		query.SortBy = &v
	}

	if r.GetVenue() != nil {
		v := r.GetVenue().GetValue()
		query.Venue = &v
	}

	return &query, nil
}
