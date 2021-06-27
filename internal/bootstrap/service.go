package bootstrap

import (
	"github.com/statistico/statistico-football-data/internal/app/grpc"
)

func (c Container) CompetitionService() *grpc.CompetitionService {
	return grpc.NewCompetitionService(c.CompetitionRepository(), c.Logger)
}

func (c Container) EventService() *grpc.EventService {
	return grpc.NewEventService(c.EventRepository(), c.Logger)
}

func (c Container) FixtureService() *grpc.FixtureService {
	return grpc.NewFixtureService(c.FixtureRepository(), c.ProtoFixtureFactory(), c.Logger)
}

func (c Container) ResultService() *grpc.ResultService {
	return grpc.NewResultService(c.FixtureRepository(), c.ProtoResultFactory(), c.Logger)
}

func (c Container) PlayerStatsService() *grpc.PlayerStatsService {
	return grpc.NewPlayerStatsService(c.FixtureRepository(), c.ProtoPlayerStatsFactory(), c.Logger)
}

func (c Container) SeasonService() *grpc.SeasonService {
	return grpc.NewSeasonService(c.SeasonRepository(), c.Logger)
}
func (c Container) TeamService() *grpc.TeamService {
	return grpc.NewTeamService(c.TeamRepository(), c.Logger)
}

func (c Container) TeamStatsService() *grpc.TeamStatsService {
	return grpc.NewTeamStatsService(
		c.FixtureRepository(),
		c.TeamStatsRepository(),
		c.FixtureTeamXGRepository(),
		c.Logger,
	)
}
