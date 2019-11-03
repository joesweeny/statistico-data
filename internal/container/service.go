package container

import (
	"github.com/statistico/statistico-data/internal/fixture"
	"github.com/statistico/statistico-data/internal/result"
	"github.com/statistico/statistico-data/internal/stats/player"
	"github.com/statistico/statistico-data/internal/stats/team"
)

func (c Container) FixtureService() *fixture.Service {
	return &fixture.Service{
		FixtureRepo: c.FixtureRepository(),
		Handler: fixture.Handler{
			CompetitionRepo: c.CompetitionRepository(),
			RoundRepo:       c.RoundRepository(),
			SeasonRepo:      c.SeasonRepository(),
			TeamRepo:        c.TeamRepository(),
			VenueRepo:       c.VenueRepository(),
			Logger:          c.Logger,
		},
		Logger: c.Logger,
	}
}

func (c Container) ResultService() *result.Service {
	return &result.Service{
		FixtureRepo: c.FixtureRepository(),
		ResultRepo:  &result.PostgresResultRepository{Connection: c.Database},
		Handler: result.Handler{
			CompetitionRepo: c.CompetitionRepository(),
			RoundRepo:       c.RoundRepository(),
			SeasonRepo:      c.SeasonRepository(),
			TeamRepo:        c.TeamRepository(),
			VenueRepo:       c.VenueRepository(),
			Logger:          c.Logger,
		},
		Logger: c.Logger,
	}
}

func (c Container) PlayerStatsService() *player_stats.Service {
	return &player_stats.Service{
		PlayerRepository: &player_stats.PostgresPlayerStatsRepository{Connection: c.Database},
		FixtureRepo:      c.FixtureRepository(),
		Logger:           c.Logger,
	}
}

func (c Container) TeamStatsService() *team_stats.Service {
	return &team_stats.Service{
		TeamRepository:    &team_stats.PostgresTeamStatsRepository{Connection: c.Database},
		FixtureRepository: c.FixtureRepository(),
		Logger:            c.Logger,
	}
}
