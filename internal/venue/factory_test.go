package venue

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

func TestCreateVenue(t *testing.T) {
	t.Run("a new domain venue struct is hydrated", func(t *testing.T) {
		t.Helper()

		m := f.createVenue(newClientVenue())

		a := assert.New(t)
		a.Equal(23, m.ID)
		a.Equal("London Stadium", m.Name)
		a.Equal("Grass", *m.Surface)
		a.Nil(m.Address)
		a.Equal("London", *m.City)
		a.Equal(60000, *m.Capacity)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.UpdatedAt.String())
	})
}

func TestUpdateVenue(t *testing.T) {
	t.Run("updates an existing venue struct", func(t *testing.T) {
		t.Helper()

		clientVenue := newClientVenue()

		m := f.createVenue(clientVenue)

		clock.Advance(10 * time.Minute)

		address := "Stratford"

		clientVenue.Address = &address
		clientVenue.Surface = "Astroturf"

		updated := f.updateVenue(clientVenue, m)

		a := assert.New(t)
		a.Equal(23, updated.ID)
		a.Equal("London Stadium", updated.Name)
		a.Equal("Astroturf", *updated.Surface)
		a.Equal("Stratford", *m.Address)
		a.Equal("London", *m.City)
		a.Equal(60000, *m.Capacity)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
		a.Equal("2019-01-14 11:35:00 +0000 UTC", m.UpdatedAt.String())
	})
}

func newClientVenue() *sportmonks.Venue {
	return &sportmonks.Venue{
		ID:       23,
		Name:     "London Stadium",
		Surface:  "Grass",
		City:     "London",
		Capacity: 60000,
		Image:    "/path/to/image",
	}
}
