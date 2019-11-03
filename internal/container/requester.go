package container

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/sportmonks"
)

func (c Container) CompetitionRequester() app.CompetitionRequester {
	return sportmonks.NewCompetitionRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) CountryRequester() app.CountryRequester {
	return sportmonks.NewCountryRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) EventRequester() app.EventRequester {
	return sportmonks.NewEventRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) FixtureRequester() app.FixtureRequester {
	return sportmonks.NewFixtureRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) RoundRequester() app.RoundRequester {
	return sportmonks.NewRoundRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) PlayerRequester() app.PlayerRequester {
	return sportmonks.NewPlayerRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) SeasonRequester() app.SeasonRequester {
	return sportmonks.NewSeasonRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) SquadRequester() app.SquadRequester {
	return sportmonks.NewSquadRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) TeamRequester() app.TeamRequester {
	return sportmonks.NewTeamRequester(c.NewSportMonksClient, c.NewLogger)
}

func (c Container) VenueRequester() app.VenueRequester {
	return sportmonks.NewVenueRequester(c.NewSportMonksClient, c.NewLogger)
}
