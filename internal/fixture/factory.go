package fixture

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/jonboulle/clockwork"
	"time"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createFixture(s *sportmonks.Fixture) *model.Fixture {
	return &model.Fixture{
		ID:         s.ID,
		SeasonID:   s.SeasonID,
		RoundID:    s.RoundID,
		VenueID:    s.VenueID,
		HomeTeamID: s.LocalteamID,
		AwayTeamID: s.VisitorteamID,
		RefereeID:  s.RefereeID,
		Date:       time.Unix(s.Time.StartingAt.Timestamp, 0),
		CreatedAt:  f.Clock.Now(),
		UpdatedAt:  f.Clock.Now(),
	}
}

func (f Factory) updateFixture(s *sportmonks.Fixture, m *model.Fixture) *model.Fixture {
	m.RoundID = s.RoundID
	m.VenueID = s.VenueID
	m.RefereeID = s.RefereeID
	m.Date = time.Unix(s.Time.StartingAt.Timestamp, 0)
	m.UpdatedAt = f.Clock.Now()

	return m
}
