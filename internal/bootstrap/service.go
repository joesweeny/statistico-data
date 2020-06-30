package bootstrap

import (
	"github.com/statistico/statistico-data/internal/app/grpc"
)

func (c Container) FixtureService() *grpc.FixtureService {
	return grpc.NewFixtureService(c.FixtureRepository(), c.ProtoFixtureFactory(), c.Logger)
}

func (c Container) PerformanceService() *grpc.PerformanceService {
	return grpc.NewPerformanceService(c.StatReader())
}

func (c Container) ResultService() *grpc.ResultService {
	return grpc.NewResultService(c.FixtureRepository(), c.ProtoResultFactory(), c.Logger)
}

func (c Container) PlayerStatsService() *grpc.PlayerStatsService {
	return grpc.NewPlayerStatsService(c.FixtureRepository(), c.ProtoPlayerStatsFactory(), c.Logger)
}

func (c Container) TeamService() *grpc.TeamService {
	return grpc.NewTeamService(c.TeamRepository(), c.Logger)
}

func (c Container) TeamStatsService() *grpc.TeamStatsService {
	return grpc.NewTeamStatsService(
		c.FixtureRepository(),
		c.FixtureTeamXGRepository(),
		c.ProtoTeamStatsFactory(),
		c.Logger,
	)
}
