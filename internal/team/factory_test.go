package team

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

func TestCreateTeam(t *testing.T) {
	t.Run("a new domain team struct is hydrated", func(t *testing.T) {
		t.Helper()

		m := f.createTeam(newClientTeam())

		a := assert.New(t)
		a.Equal(56, m.ID)
		a.Equal("West Ham United", m.Name)
		a.Equal("WHU", *m.ShortCode)
		a.Equal(8, *m.CountryID)
		a.False(m.NationalTeam)
		a.Equal(1898, *m.Founded)
		a.Equal("/path/to/logo", *m.Logo)
		a.Equal(99, m.VenueID)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.UpdatedAt.String())
	})
}

func TestUpdateTeam(t *testing.T) {
	t.Run("update an existing team struct", func(t *testing.T) {
		t.Helper()

		clientTeam := newClientTeam()

		m := f.createTeam(clientTeam)

		clock.Advance(10 * time.Minute)

		var logo *string
		clientTeam.Name = "West Ham London"
		clientTeam.VenueID = 78908
		clientTeam.LogoPath = logo

		updated := f.updateTeam(clientTeam, m)

		a := assert.New(t)
		a.Equal(56, updated.ID)
		a.Equal("West Ham London", updated.Name)
		a.Equal("WHU", *updated.ShortCode)
		a.Equal(8, *updated.CountryID)
		a.False(updated.NationalTeam)
		a.Equal(1898, *updated.Founded)
		a.Nil(m.Logo)
		a.Equal(78908, m.VenueID)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", updated.CreatedAt.String())
		a.Equal("2019-01-14 11:35:00 +0000 UTC", updated.UpdatedAt.String())
	})
}

func newClientTeam() *sportmonks.Team {
	logo := "/path/to/logo"
	twitter := "westhamunited"

	return &sportmonks.Team{
		ID:           56,
		LegacyID:     34,
		Name:         "West Ham United",
		ShortCode:    "WHU",
		Twitter:      &twitter,
		CountryID:    8,
		NationalTeam: false,
		Founded:      1898,
		LogoPath:     &logo,
		VenueID:      99,
	}
}
