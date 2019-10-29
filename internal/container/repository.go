package container

import (
	"github.com/statistico/statistico-data/internal/app/postgres"
)

func (c Container) CompetitionRepository() *postgres.CompetitionRepository {
	return postgres.NewCompetitionRepository(c.Database, c.Clock)
}

func (c Container) CountryRepository() *postgres.CountryRepository {
	return postgres.NewCountryRepository(c.Database, c.Clock)
}

func (c Container) RoundRepository() *postgres.RoundRepository {
	return postgres.NewRoundRepository(c.Database, c.Clock)
}

func (c Container) SeasonRepository() *postgres.SeasonRepository {
	return postgres.NewSeasonRepository(c.Database, c.Clock)
}

func (c Container) TeamRepository() *postgres.TeamRepository {
	return postgres.NewTeamRepository(c.Database, c.Clock)
}

func (c Container) VenueRepository() *postgres.VenueRepository {
	return postgres.NewVenueRepository(c.Database, c.Clock)
}
