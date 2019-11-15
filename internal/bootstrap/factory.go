package bootstrap

import (
	"github.com/statistico/statistico-data/internal/app/grpc/factory"
	"github.com/statistico/statistico-data/internal/app/rest"
)

func (c Container) ProtoFixtureFactory() *factory.FixtureFactory {
	return factory.NewFixtureFactory(c.RoundRepository(), c.TeamRepository(), c.VenueRepository(), c.Logger)
}

func (c Container) ProtoPlayerStatsFactory() *factory.PlayerStatsFactory {
	return factory.NewPlayerStatsFactory(c.PlayerStatsRepository(), c.Logger)
}

func (c Container) ProtoResultFactory() *factory.ResultFactory {
	return factory.NewResultFactory(c.ResultRepository(), c.TeamRepository(), c.VenueRepository(), c.Logger)
}

func (c Container) ProtoTeamStatsFactory() *factory.TeamStatsFactory {
	return factory.NewTeamStatsFactory(c.TeamStatsRepository(), c.Logger)
}

func (c Container) RestFixtureFactory() *rest.FixtureFactory {
	return rest.NewFixtureFactory(c.RoundRepository(), c.TeamRepository(), c.VenueRepository())
}
