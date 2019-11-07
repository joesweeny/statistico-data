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

func (c Container) EventRepository() *postgres.EventRepository {
	return postgres.NewEventRepository(c.Database, c.Clock)
}

func (c Container) FixtureRepository() *postgres.FixtureRepository {
	return postgres.NewFixtureRepository(c.Database, c.Clock)
}

func (c Container) PlayerRepository() *postgres.PlayerRepository {
	return postgres.NewPlayerRepository(c.Database, c.Clock)
}

func (c Container) RoundRepository() *postgres.RoundRepository {
	return postgres.NewRoundRepository(c.Database, c.Clock)
}

func (c Container) ResultRepository() *postgres.ResultRepository {
	return postgres.NewResultRepository(c.Database, c.Clock)
}

func (c Container) SeasonRepository() *postgres.SeasonRepository {
	return postgres.NewSeasonRepository(c.Database, c.Clock)
}

func (c Container) SquadRepository() *postgres.SquadRepository {
	return postgres.NewSquadRepository(c.Database, c.Clock)
}

func (c Container) TeamRepository() *postgres.TeamRepository {
	return postgres.NewTeamRepository(c.Database, c.Clock)
}

func (c Container) TeamStatsRepository() *postgres.TeamStatsRepository {
	return postgres.NewTeamStatsRepository(c.Database, c.Clock)
}

func (c Container) VenueRepository() *postgres.VenueRepository {
	return postgres.NewVenueRepository(c.Database, c.Clock)
}
