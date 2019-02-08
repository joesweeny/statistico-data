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

type Processor interface {
	Process(command string, done chan bool)
}

func (c Container) CompetitionProcessor() *competition.Processor {
	return &competition.Processor{
		Repository: &competition.PostgresCompetitionRepository{Connection: c.Database},
		Factory:    competition.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) CountryProcessor() *country.Processor {
	return &country.Processor{
		Repository: &country.PostgresCountryRepository{Connection: c.Database},
		Factory:    country.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     logger(),
	}
}

func (c Container) FixtureProcessor() *fixture.Processor {
	return &fixture.Processor{
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

func (c Container) RoundProcessor() *round.Processor {
	return &round.Processor{
		Repository: &round.PostgresRoundRepository{Connection: c.Database},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    round.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) SeasonProcessor() *season.Processor {
	return &season.Processor{
		Repository: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    season.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) SquadProcessor() *squad.Processor {
	return &squad.Processor{
		Repository: &squad.PostgresSquadRepository{Connection: c.Database},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    squad.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) TeamProcessor() *team.Processor {
	return &team.Processor{
		Repository: &team.PostgresTeamRepository{Connection: c.Database},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    team.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) VenueProcessor() *venue.Processor {
	return &venue.Processor{
		Repository: &venue.PostgresVenueRepository{Connection: c.Database},
		SeasonRepo: &season.PostgresSeasonRepository{Connection: c.Database},
		Factory:    venue.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}
