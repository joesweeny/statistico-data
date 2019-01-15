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

		newCountry := f.create(newClientCountry(), uuid.FromStringOrNil("794f464a-a2e3-4d31-ac90-b7e255a28030"))

		a := assert.New(t)

		a.Equal("794f464a-a2e3-4d31-ac90-b7e255a28030", newCountry.ID.String())
		a.Equal(180, newCountry.ExternalID)
		a.Equal("England", newCountry.Name)
		a.Equal("Europe", newCountry.Continent)
		a.Equal("ENG", newCountry.ISO)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", newCountry.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", newCountry.UpdatedAt.String())
	})
}

func TestFactoryUpdate(t *testing.T) {
	t.Run("updates an existing country struct", func (t *testing.T) {
		t.Helper()

		clientCountry := newClientCountry()

		newCountry := f.create(clientCountry, uuid.FromStringOrNil("794f464a-a2e3-4d31-ac90-b7e255a28030"))

		clock.Advance(10 * time.Minute)

		clientCountry.Name = "United Kingdom"

		updated := f.update(clientCountry, newCountry)

		a := assert.New(t)

		a.Equal(180, updated.ExternalID)
		a.Equal("United Kingdom", updated.Name)
		a.Equal("Europe", updated.Continent)
		a.Equal("ENG", updated.ISO)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", updated.CreatedAt.String())
		a.Equal("2019-01-14 11:35:00 +0000 UTC", updated.UpdatedAt.String())
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