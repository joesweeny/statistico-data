package container

import (
	"github.com/statistico/statistico-data/internal/app/process"
	"github.com/statistico/statistico-data/internal/competition"
	"github.com/statistico/statistico-data/internal/event"
	"github.com/statistico/statistico-data/internal/fixture"
	"github.com/statistico/statistico-data/internal/player"
	"github.com/statistico/statistico-data/internal/result"
	"github.com/statistico/statistico-data/internal/round"
	"github.com/statistico/statistico-data/internal/season"
	"github.com/statistico/statistico-data/internal/squad"
	"github.com/statistico/statistico-data/internal/stats/player"
	"github.com/statistico/statistico-data/internal/stats/team"
	"github.com/statistico/statistico-data/internal/team"
	"github.com/statistico/statistico-data/internal/venue"
)

type Processor interface {
	Process(command string, option string, done chan bool)
}

func (c Container) CompetitionProcessor() *competition.Processor {
	return &competition.Processor{
		Repository: &competition.PostgresCompetitionRepository{Connection: c.Database},
		Factory:    competition.Factory{Clock: clock()},
		Client:     c.SportMonksClient,
		Logger:     c.Logger,
	}
}

func (c Container) CountryProcessor() *process.CountryProcessor {
	return process.NewCountryProcessor(
		c.CountryRepository(),
		c.CountryRequester(),
		c.NewLogger,
	)
}

func (c Container) EventProcessor() event.Processor {
	return event.Processor{
		Repository: &event.PostgresEventRepository{Connection: c.Database},
		Factory:    event.Factory{Clock: clock()},
		Logger:     c.Logger,
		FixtureRepo:     &fixture.PostgresFixtureRepository{Connection: c.Database},
		Client:     c.SportMonksClient,
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

func (c Container) PlayerStatsProcessor() player_stats.Processor {
	return player_stats.Processor{
		PlayerRepository: &player_stats.PostgresPlayerStatsRepository{Connection: c.Database},
		PlayerFactory:    player_stats.PlayerFactory{Clock: clock()},
		Logger:           c.Logger,
		FixtureRepo:     &fixture.PostgresFixtureRepository{Connection: c.Database},
		Client:     c.SportMonksClient,
	}
}

func (c Container) ResultProcessor() *result.Processor {
	return &result.Processor{
		Repository:      &result.PostgresResultRepository{Connection: c.Database},
		FixtureRepo:     &fixture.PostgresFixtureRepository{Connection: c.Database},
		Factory:         result.Factory{Clock: c.Clock},
		Client:          c.SportMonksClient,
		Logger:          c.Logger,
		Clock:	 		 c.Clock,
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

func (c Container) TeamStatsProcessor() team_stats.Processor {
	return team_stats.Processor{
		TeamRepository: &team_stats.PostgresTeamStatsRepository{Connection: c.Database},
		TeamFactory:    team_stats.TeamFactory{Clock: clock()},
		Logger:         c.Logger,
		FixtureRepo:     &fixture.PostgresFixtureRepository{Connection: c.Database},
		Client:     c.SportMonksClient,
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
