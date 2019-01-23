package season

import (
	"time"
	"github.com/jonboulle/clockwork"
	"testing"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/stretchr/testify/assert"
)

var t = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)
var clock = clockwork.NewFakeClockAt(t)
var f = Factory{clock}

func TestFactoryCreateSeason(t *testing.T) {
	t.Run("a new domain season struct is hydrated", func(t *testing.T) {
		t.Helper()

		s := f.createSeason(newClientSeason())

		a := assert.New(t)
		a.Equal(100, s.ID)
		a.Equal("2018-2019", s.Name)
		a.Equal(231, s.LeagueID)
		a.True(s.IsCurrent)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", s.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", s.UpdatedAt.String())
	})

	t.Run("updates an existing season struct", func(t *testing.T) {
		t.Helper()

		clientSeason := newClientSeason()

		s := f.createSeason(clientSeason)

		clock.Advance(10 * time.Minute)

		clientSeason.CurrentSeason = false

		updated := f.updateSeason(clientSeason, s)

		a := assert.New(t)
		a.Equal(100, updated.ID)
		a.Equal("2018-2019", updated.Name)
		a.Equal(231, updated.LeagueID)
		a.False(s.IsCurrent)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", s.CreatedAt.String())
		a.Equal("2019-01-14 11:35:00 +0000 UTC", s.UpdatedAt.String())

	})
}

func newClientSeason() *sportmonks.Season {
	return &sportmonks.Season{
		ID:             100,
		Name:           "2018-2019",
		LeagueID:       231,
		CurrentSeason:  true,
		CurrentRoundID: 10,
		CurrentStageID: 567,
		Fixtures:       struct {
			Data []sportmonks.Fixture `json:"data"`
		} {},
		Results:       struct {
			Data []sportmonks.Fixture `json:"data"`
		} {},
	}
}