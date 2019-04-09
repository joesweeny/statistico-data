package season

import (
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

		clientSeason.IsCurrentSeason = false

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
	var round = 10
	var stage = 567
	return &sportmonks.Season{
		ID:              100,
		Name:            "2018-2019",
		LeagueID:        231,
		IsCurrentSeason: true,
		CurrentRoundID:  &round,
		CurrentStageID:  &stage,
		Fixtures: struct {
			Data []sportmonks.Fixture `json:"data"`
		}{},
	}
}
