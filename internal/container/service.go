package container

import (
	"github.com/joesweeny/statshub/internal/competition"
	"github.com/joesweeny/statshub/internal/country"
	"github.com/joesweeny/statshub/internal/fixture"
	"github.com/joesweeny/statshub/internal/player"
	"github.com/joesweeny/statshub/internal/round"
	"github.com/joesweeny/statshub/internal/season"
	"github.com/joesweeny/statshub/internal/squad"
	"github.com/joesweeny/statshub/internal/team"
	"github.com/joesweeny/statshub/internal/venue"
)

type Service interface {
	Process(command string, done chan bool)
}

func (c Container) CompetitionService() *competition.Service {
	return &competition.Service{
		Repository: &competition.PostgresCompetitionRepository{Connection: c.Database},
		Factory:    competition.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) CountryService() *country.Service {
	return &country.Service{
		Repository: &country.PostgresCountryRepository{Connection: c.Database},
		Factory:    country.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     logger(),
	}
}

func (c Container) FixtureService() *fixture.Service {
	return &fixture.Service{
		Repository: &fixture.PostgresFixtureRepository{Connection: c.Database},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    fixture.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) PlayerProcessor() *player.Processor {
	return &player.Processor{
		Repository: &player.PostgresPlayerRepository{Connection: c.Database},
		SquadRepo:  &squad.PostgresSquadRepository{Connection: c.Database},
		Factory:    player.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) RoundService() *round.Service {
	return &round.Service{
		Repository: &round.PostgresRoundRepository{Connection: c.Database},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    round.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) SeasonService() *season.Service {
	return &season.Service{
		Repository: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    season.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) SquadService() *squad.Service {
	return &squad.Service{
		Repository: &squad.PostgresSquadRepository{Connection: c.Database},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    squad.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) TeamService() *team.Service {
	return &team.Service{
		Repository: &team.PostgresTeamRepository{Connection: c.Database},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    team.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) VenueService() *venue.Service {
	return &venue.Service{
		Repository: &venue.PostgresVenueRepository{Connection: c.Database},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    venue.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}
