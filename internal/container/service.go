package container

import (
	"github.com/statistico/statistico-data/internal/competition"
	"github.com/statistico/statistico-data/internal/fixture"
	"github.com/statistico/statistico-data/internal/result"
	"github.com/statistico/statistico-data/internal/round"
	"github.com/statistico/statistico-data/internal/season"
	"github.com/statistico/statistico-data/internal/team"
	"github.com/statistico/statistico-data/internal/venue"
	"github.com/statistico/statistico-data/internal/stats/player"
)

func (c Container) FixtureService() *fixture.Service {
	return &fixture.Service{
		Repository: &fixture.PostgresFixtureRepository{Connection: c.Database},
		Handler: fixture.Handler{
			CompetitionRepo: &competition.PostgresCompetitionRepository{Connection: c.Database},
			RoundRepo:       &round.PostgresRoundRepository{Connection: c.Database},
			SeasonRepo:      &season.PostgresSeasonRepository{Connection: c.Database},
			TeamRepo:        &team.PostgresTeamRepository{Connection: c.Database},
			VenueRepo:       &venue.PostgresVenueRepository{Connection: c.Database},
			Logger:          c.Logger,
		},
		Logger: c.Logger,
	}
}

func (c Container) ResultService() *result.Service {
	return &result.Service{
		FixtureRepo: &fixture.PostgresFixtureRepository{Connection: c.Database},
		ResultRepo:  &result.PostgresResultRepository{Connection: c.Database},
		Handler: result.Handler{
			CompetitionRepo: &competition.PostgresCompetitionRepository{Connection: c.Database},
			RoundRepo:       &round.PostgresRoundRepository{Connection: c.Database},
			SeasonRepo:      &season.PostgresSeasonRepository{Connection: c.Database},
			TeamRepo:        &team.PostgresTeamRepository{Connection: c.Database},
			VenueRepo:       &venue.PostgresVenueRepository{Connection: c.Database},
			Logger:          c.Logger,
		},
		Logger: c.Logger,
	}
}

func (c Container) PlayerStatsService() *player_stats.Service {
	return &player_stats.Service{
		PlayerRepository: &player_stats.PostgresPlayerStatsRepository{Connection: c.Database},
		FixtureRepo: &fixture.PostgresFixtureRepository{Connection: c.Database},
		Logger: c.Logger,
	}
}
