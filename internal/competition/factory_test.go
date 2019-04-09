package competition

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

func TestFactoryCreateCompetiton(t *testing.T) {
	t.Run("a new domain competition struct is hydrated", func(t *testing.T) {
		t.Helper()

		c := f.createCompetition(newClientLeague())

		a := assert.New(t)

		a.Equal(564, c.ID)
		a.Equal("Serie A", c.Name)
		a.Equal(32, c.CountryID)
		a.Equal(false, c.IsCup)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", c.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", c.UpdatedAt.String())
	})
}

func TestFactoryUpdateCompetition(t *testing.T) {
	t.Run("updates an existing competition struct", func(t *testing.T) {
		t.Helper()

		clientLeague := newClientLeague()

		c := f.createCompetition(clientLeague)

		clock.Advance(10 * time.Minute)

		clientLeague.IsCup = true
		clientLeague.Name = "FA Cup"

		updated := f.updateCompetition(clientLeague, c)

		a := assert.New(t)

		a.Equal(564, updated.ID)
		a.Equal("FA Cup", updated.Name)
		a.Equal(32, updated.CountryID)
		a.Equal(true, updated.IsCup)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", updated.CreatedAt.String())
		a.Equal("2019-01-14 11:35:00 +0000 UTC", updated.UpdatedAt.String())
	})
}

func newClientLeague() *sportmonks.League {
	return &sportmonks.League{
		ID:              564,
		LegacyID:        3491,
		CountryID:       32,
		Name:            "Serie A",
		IsCup:           false,
		CurrentSeasonID: 23,
		CurrentRoundID:  98,
		CurrentStageID:  87,
		LiveStandings:   true,
		Coverage: struct {
			TopscorerGoals   bool `json:"topscorer_goals"`
			TopscorerAssists bool `json:"topscorer_assists"`
			TopscorerCards   bool `json:"topscorer_cards"`
		}{
			TopscorerGoals:   true,
			TopscorerAssists: false,
			TopscorerCards:   true,
		},
		Seasons: struct {
			Data []sportmonks.Season `json:"data"`
		}{
			Data: []sportmonks.Season{},
		},
	}
}
