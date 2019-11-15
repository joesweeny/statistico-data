package bootstrap

import "github.com/statistico/statistico-data/internal/app/rest"

func (c Container) RestFixtureFactory() *rest.FixtureFactory {
	return rest.NewFixtureFactory(c.RoundRepository(), c.TeamRepository(), c.VenueRepository())
}
