package bootstrap

import "github.com/statistico/statistico-data/internal/app/rest"

func (c Container) FixtureHandler() *rest.FixtureHandler {
	return rest.NewFixtureHandler(
		c.FixtureRepository(),
		c.FixtureFactory(),
		c.Logger,
	)
}
