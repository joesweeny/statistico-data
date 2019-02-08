package player

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var t = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)
var clock = clockwork.NewFakeClockAt(t)
var f = Factory{clock}

func TestCreatePlayer(t *testing.T) {
	m := f.createPlayer(newClientPlayer())

	a := assert.New(t)

	a.Equal(58, m.ID)
	a.Equal(462, m.CountryId)
	a.Equal("Felipe", m.FirstName)
	a.Equal("Anderson", m.LastName)
	a.Equal("Sao Paulo", *m.BirthPlace)
	a.Equal("1984-03-12", *m.DateOfBirth)
	a.Equal(3, m.PositionID)
	a.Equal("/path/to/image", *m.Image)
	a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
	a.Equal("2019-01-14 11:25:00 +0000 UTC", m.UpdatedAt.String())
}

func newClientPlayer() *sportmonks.Player {
	return &sportmonks.Player{
		PlayerID:     58,
		TeamID:       999,
		CountryID:    462,
		PositionID:   3,
		CommonName:   "Felipe Anderson",
		FullName:     "Felipe Anderson",
		FirstName:    "Felipe",
		LastName:     "Anderson",
		Nationality:  "Brazilian",
		BirthDate:    "1984-03-12",
		BirthCountry: "Brazil",
		Birthplace:   "Sao Paulo",
		Height:       "A Giant",
		Weight:       "Skinny",
		ImagePath:    "/path/to/image",
	}
}
