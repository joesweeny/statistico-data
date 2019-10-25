package container

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/sportmonks"
)

func (c Container) CompetitionRequester() app.CompetitionRequester {
	return sportmonks.NewCompetitionRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) CountryRequester() app.CountryRequester {
	return sportmonks.NewCountryRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) VenueRequester() app.VenueRequester {
	return sportmonks.NewVenueRequester(c.NewSportMonksClient, c.NewLogger)
}
