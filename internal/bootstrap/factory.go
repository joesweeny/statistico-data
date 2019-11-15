package bootstrap

import (
	"github.com/statistico/statistico-data/internal/app/grpc"
	"github.com/statistico/statistico-data/internal/app/rest"
)

func (c Container) RestFixtureFactory() *rest.FixtureFactory {
	return rest.NewFixtureFactory(c.RoundRepository(), c.TeamRepository(), c.VenueRepository())
}

func (c Container) ProtoFixtureFactory() *grpc.FixtureFactory {
	return grpc.NewFixtureFactory(c.RoundRepository(), c.TeamRepository(), c.VenueRepository())
}
