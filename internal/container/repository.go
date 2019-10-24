package container

import (
	"github.com/statistico/statistico-data/internal/app/postgres"
)

func (c Container) CountryRepository() *postgres.CountryRepository {
	return postgres.NewCountryRepository(c.Database, c.Clock)
}

func (c Container) VenueRepository() *postgres.VenueRepository {
	return postgres.NewVenueRepository(c.Database, c.Clock)
}
