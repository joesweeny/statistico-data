package bootstrap

import (
	"github.com/statistico/statistico-football-data/internal/app"
	"github.com/statistico/statistico-football-data/internal/app/sportmonks"
)

func (c Container) CompetitionRequester() app.CompetitionRequester {
	return sportmonks.NewCompetitionRequester(c.SportMonksClient, c.Logger)
}

func (c Container) CountryRequester() app.CountryRequester {
	return sportmonks.NewCountryRequester(c.SportMonksClient, c.Logger)
}

func (c Container) EventRequester() app.EventRequester {
	return sportmonks.NewEventRequester(c.SportMonksClient, c.Logger)
}

func (c Container) FixtureRequester() app.FixtureRequester {
	return sportmonks.NewFixtureRequester(c.SportMonksClient, c.Logger)
}

func (c Container) RoundRequester() app.RoundRequester {
	return sportmonks.NewRoundRequester(c.SportMonksClient, c.Logger)
}

func (c Container) ResultRequester() app.ResultRequester {
	return sportmonks.NewResultRequester(c.SportMonksClient, c.Logger)
}

func (c Container) PlayerRequester() app.PlayerRequester {
	return sportmonks.NewPlayerRequester(c.SportMonksClient, c.Logger)
}

func (c Container) PlayerStatsRequester() app.PlayerStatRequester {
	return sportmonks.NewPlayerStatsRequester(c.SportMonksClient, c.Logger)
}

func (c Container) SeasonRequester() app.SeasonRequester {
	return sportmonks.NewSeasonRequester(c.SportMonksClient, c.Logger)
}

func (c Container) SquadRequester() app.SquadRequester {
	return sportmonks.NewSquadRequester(c.SportMonksClient, c.Logger)
}

func (c Container) TeamRequester() app.TeamRequester {
	return sportmonks.NewTeamRequester(c.SportMonksClient, c.Logger)
}

func (c Container) TeamStatsRequester() app.TeamStatsRequester {
	return sportmonks.NewTeamStatsRequester(c.SportMonksClient, c.Logger)
}

func (c Container) VenueRequester() app.VenueRequester {
	return sportmonks.NewVenueRequester(c.SportMonksClient, c.Logger)
}
