package bootstrap

import "github.com/statistico/statistico-data/internal/app/rest"

func (c Container) RestFixtureHandler() *rest.FixtureHandler {
	return rest.NewFixtureHandler(c.FixtureRepository(), c.RestFixtureFactory())
}
