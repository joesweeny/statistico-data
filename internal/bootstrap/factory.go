package bootstrap

import "github.com/statistico/statistico-data/internal/app/factory"

func (c Container) FixtureFactory() *factory.FixtureFactory {
	return factory.NewFixtureFactory(
		c.RoundRepository(),
		c.TeamRepository(),
		c.VenueRepository(),
	)
}
