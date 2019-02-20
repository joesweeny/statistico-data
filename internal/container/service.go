package container

import (
	"github.com/joesweeny/statshub/internal/fixture"
	"github.com/joesweeny/statshub/internal/competition"
	"github.com/joesweeny/statshub/internal/season"
	"github.com/joesweeny/statshub/internal/team"
	"github.com/joesweeny/statshub/internal/venue"
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
