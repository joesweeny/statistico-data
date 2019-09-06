package container

import (
	"github.com/statistico/statistico-data/internal/app/postgres"
)

func (c Container) CountryRepository() *postgres.CountryRepository {
	return postgres.NewCountryRepository(c.Database, c.Clock)
}
