package container

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/sportmonks"
)

func (c Container) CompetitionRequester() app.CompetitionRequester {
	return sportmonks.NewCompetitionRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) CountryRequester() app.CountryRequester {
	return sportmonks.NewCountryRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) EventRequester() app.EventRequester {
	return sportmonks.NewEventRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) FixtureRequester() app.FixtureRequester {
	return sportmonks.NewFixtureRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) RoundRequester() app.RoundRequester {
	return sportmonks.NewRoundRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) ResultRequester() app.ResultRequester {
	return sportmonks.NewResultRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) PlayerRequester() app.PlayerRequester {
	return sportmonks.NewPlayerRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) PlayerStatsRequester() app.PlayerStatRequester {
	return sportmonks.NewPlayerStatsRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) SeasonRequester() app.SeasonRequester {
	return sportmonks.NewSeasonRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) SquadRequester() app.SquadRequester {
	return sportmonks.NewSquadRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) TeamRequester() app.TeamRequester {
	return sportmonks.NewTeamRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) TeamStatsRequester() app.TeamStatsRequester {
	return sportmonks.NewTeamStatsRequester(c.SportMonksClient, c.NewLogger)
}

func (c Container) VenueRequester() app.VenueRequester {
	return sportmonks.NewVenueRequester(c.SportMonksClient, c.NewLogger)
}
