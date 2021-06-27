package bootstrap

import (
	"github.com/statistico/statistico-football-data/internal/app/grpc/factory"
	"github.com/statistico/statistico-football-data/internal/app/rest"
)

func (c Container) ProtoFixtureFactory() *factory.FixtureFactory {
	return factory.NewFixtureFactory(
		c.CompetitionRepository(),
		c.RoundRepository(),
		c.SeasonRepository(),
		c.TeamRepository(),
		c.VenueRepository(),
		c.Logger,
	)
}

func (c Container) ProtoPlayerStatsFactory() *factory.PlayerStatsFactory {
	return factory.NewPlayerStatsFactory(c.PlayerStatsRepository(), c.Logger)
}

func (c Container) ProtoResultFactory() *factory.ResultFactory {
	return factory.NewResultFactory(
		c.ResultRepository(),
		c.RoundRepository(),
		c.SeasonRepository(),
		c.TeamRepository(),
		c.TeamStatsRepository(),
		c.VenueRepository(),
		c.Logger,
	)
}

func (c Container) RestFixtureFactory() *rest.FixtureFactory {
	return rest.NewFixtureFactory(c.RoundRepository(), c.TeamRepository(), c.VenueRepository())
}
