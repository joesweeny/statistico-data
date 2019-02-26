package container

import (
	"github.com/joesweeny/statistico-data/internal/fixture"
	"github.com/joesweeny/statistico-data/internal/competition"
	"github.com/joesweeny/statistico-data/internal/season"
	"github.com/joesweeny/statistico-data/internal/team"
	"github.com/joesweeny/statistico-data/internal/venue"
)

func (c Container) FixtureService() *fixture.Service {
	return &fixture.Service{
		Repository: &fixture.PostgresFixtureRepository{Connection: c.Database},
		Handler: fixture.Handler{
			CompetitionRepo: &competition.PostgresCompetitionRepository{Connection: c.Database},
			SeasonRepo: &season.PostgresSeasonRepository{Connection: c.Database},
			TeamRepo: &team.PostgresTeamRepository{Connection: c.Database},
			VenueRepo: &venue.PostgresVenueRepository{Connection: c.Database},
		},
		Logger: c.Logger,
	}
}
