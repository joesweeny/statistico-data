package country

import (
	"testing"
	"github.com/jonboulle/clockwork"
	"time"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/stretchr/testify/assert"
	"github.com/satori/go.uuid"
)

var t = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)
var clock = clockwork.NewFakeClockAt(t)
var f = Factory{clock}

func TestFactoryCreate(t *testing.T) {
	t.Run("a new domain country struct is hydrated", func (t *testing.T) {
		t.Helper()

		newCountry := f.create(newClientCountry())

		assert.IsType(t, uuid.UUID{}, newCountry.ID)
		assert.Equal(t, 180, newCountry.ExternalID)
		assert.Equal(t, "England", newCountry.Name)
		assert.Equal(t, "Europe", newCountry.Continent)
		assert.Equal(t, "ENG", newCountry.ISO)
		assert.Equal(t, "2019-01-14 11:25:00 +0000 UTC", newCountry.CreatedAt.String())
		assert.Equal(t, "2019-01-14 11:25:00 +0000 UTC", newCountry.UpdatedAt.String())
	})
}

func TestFactoryUpdate(t *testing.T) {
	t.Run("updates an existing country struct", func (t *testing.T) {
		t.Helper()

		clientCountry := newClientCountry()

		newCountry := f.create(clientCountry)

		clock.Advance(10 * time.Minute)

		clientCountry.Name = "United Kingdom"

		updated := f.update(clientCountry, newCountry)

		assert.Equal(t, 180, updated.ExternalID)
		assert.Equal(t, "United Kingdom", updated.Name)
		assert.Equal(t, "Europe", updated.Continent)
		assert.Equal(t, "ENG", updated.ISO)
		assert.Equal(t, "2019-01-14 11:25:00 +0000 UTC", updated.CreatedAt.String())
		assert.Equal(t, "2019-01-14 11:35:00 +0000 UTC", updated.UpdatedAt.String())
	})
}

func newClientCountry() sportmonks.Country {
	country := sportmonks.Country{
		ID: 180,
		Name: "England",
		Extra: struct {
			Continent string `json:"continent"`
			SubRegion string `json:"sub_region"`
			WorldRegion string `json:"world_region"`
			Fifa interface{} `json:"fifa,string"`
			ISO string `json:"iso"`
			Longitude string `json:"longitude"`
			Latitude string `json:"latitude"`
		} {
			Continent:   "Europe",
			SubRegion:   "Western Europe",
			WorldRegion: "Europe",
			Fifa:        "ENG",
			ISO:         "ENG",
		},
	}

	return country
}