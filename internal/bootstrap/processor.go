package bootstrap

import (
	"github.com/statistico/statistico-data/internal/app/process"
)

type Processor interface {
	Process(command string, option string, done chan bool)
}

func (c Container) CompetitionProcessor() *process.CompetitionProcessor {
	return process.NewCompetitionProcessor(
		c.CompetitionRepository(),
		c.CompetitionRequester(),
		c.Logger,
	)
}

func (c Container) CountryProcessor() *process.CountryProcessor {
	return process.NewCountryProcessor(
		c.CountryRepository(),
		c.CountryRequester(),
		c.Logger,
	)
}

func (c Container) EventProcessor() *process.EventProcessor {
	return process.NewEventProcessor(
		c.EventRepository(),
		c.SeasonRepository(),
		c.EventRequester(),
		c.Clock,
		c.Logger,
	)
}

func (c Container) FixtureProcessor() *process.FixtureProcessor {
	return process.NewFixtureProcessor(
		c.FixtureRepository(),
		c.SeasonRepository(),
		c.FixtureRequester(),
		c.Logger,
	)
}

func (c Container) FixtureTeamXGProcessor() *process.FixtureTeamXGProcessor {
	return process.NewFixtureTeamXGProcessor(
		c.FixtureTeamXGRepository(),
		c.FixtureRepository(),
		c.UnderstatParser,
		c.Logger,
	)
}

func (c Container) PlayerProcessor() *process.PlayerProcessor {
	return process.NewPlayerProcessor(
		c.PlayerRepository(),
		c.SquadRepository(),
		c.PlayerRequester(),
		c.Logger,
	)
}

func (c Container) PlayerStatsProcessor() *process.PlayerStatsProcessor {
	return process.NewPlayerStatsProcessor(
		c.PlayerStatsRepository(),
		c.SeasonRepository(),
		c.PlayerStatsRequester(),
		c.Clock,
		c.Logger,
	)
}

func (c Container) ResultProcessor() *process.ResultProcessor {
	return process.NewResultProcessor(
		c.ResultRepository(),
		c.SeasonRepository(),
		c.ResultRequester(),
		c.Clock,
		c.Logger,
	)
}

func (c Container) RoundProcessor() *process.RoundProcessor {
	return process.NewRoundProcessor(
		c.RoundRepository(),
		c.SeasonRepository(),
		c.RoundRequester(),
		c.Logger,
	)
}

func (c Container) SeasonProcessor() *process.SeasonProcessor {
	return process.NewSeasonProcessor(
		c.SeasonRepository(),
		c.SeasonRequester(),
		c.Logger,
	)
}

func (c Container) SquadProcessor() *process.SquadProcessor {
	return process.NewSquadProcessor(
		c.SquadRepository(),
		c.SeasonRepository(),
		c.SquadRequester(),
		c.Logger,
	)
}

func (c Container) TeamProcessor() *process.TeamProcessor {
	return process.NewTeamProcessor(
		c.TeamRepository(),
		c.SeasonRepository(),
		c.TeamRequester(),
		c.Logger,
	)
}

func (c Container) TeamStatsProcessor() *process.TeamStatsProcessor {
	return process.NewTeamStatsProcessor(
		c.TeamStatsRepository(),
		c.SeasonRepository(),
		c.TeamStatsRequester(),
		c.Clock,
		c.Logger,
	)
}

func (c Container) VenueProcessor() *process.VenueProcessor {
	return process.NewVenueProcessor(
		c.VenueRepository(),
		c.SeasonRepository(),
		c.VenueRequester(),
		c.Logger,
	)
}
