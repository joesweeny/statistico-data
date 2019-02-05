package squad

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

func TestCreateSquad(t *testing.T) {
	t.Run("a new domain squad struct is hydrated", func(t *testing.T) {
		t.Helper()

		m := f.createSquad(43, 999, newClientSquad(45, 4501, 2))

		a := assert.New(t)
		a.Equal(43, m.SeasonID)
		a.Equal(999, m.TeamID)
		a.Equal([]int{45, 4501, 2}, m.PlayerIDs)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.UpdatedAt.String())
	})
}

func TestUpdateSquad(t *testing.T) {
	t.Run("update an existing squad struct", func(t *testing.T) {
		t.Helper()

		m := f.createSquad(43, 999, newClientSquad(45, 4501, 2))

		clock.Advance(10 * time.Minute)

		updated := f.updateSquad(newClientSquad(480, 650, 2), m)

		a := assert.New(t)
		a.Equal(43, updated.SeasonID)
		a.Equal(999, updated.TeamID)
		a.Equal([]int{480, 650, 2}, updated.PlayerIDs)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", updated.CreatedAt.String())
		a.Equal("2019-01-14 11:35:00 +0000 UTC", updated.UpdatedAt.String())
	})
}

func newClientSquad(id1, id2, id3 int) *[]sportmonks.SquadPlayer {
	player1 := sportmonks.SquadPlayer{
		PlayerID:   id1,
		PositionID: 2,
		Number:     34,
	}

	player2 := sportmonks.SquadPlayer{
		PlayerID:   id2,
		PositionID: 1,
		Number:     13,
	}

	player3 := sportmonks.SquadPlayer{
		PlayerID:   id3,
		PositionID: 4,
		Number:     30,
	}

	players := []sportmonks.SquadPlayer{player1, player2, player3}

	return &players
}
