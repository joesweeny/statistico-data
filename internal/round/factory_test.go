package round

import (
	"time"
	"github.com/jonboulle/clockwork"
	"testing"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/stretchr/testify/assert"
	"github.com/joesweeny/statshub/internal/model"
)

var t = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)
var clock = clockwork.NewFakeClockAt(t)
var f = Factory{clock}

func TestCreateRound(t *testing.T) {
	t.Run("a new domain round is hydrated", func(t *testing.T) {
		t.Helper()

		m, err := f.createRound(newClientRound("2019-03-12", "2019-03-19"))

		if err != nil {
			t.Fatalf("Test failed: want nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(54, m.ID)
		a.Equal("2", m.Name)
		a.Equal(45, m.SeasonID)
		a.Equal("2019-03-12 00:00:00 +0000 UTC", m.StartDate.String())
		a.Equal("2019-03-19 00:00:00 +0000 UTC", m.EndDate.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.UpdatedAt.String())
	})

	t.Run("returns error if date cannot be parsed", func(t *testing.T) {
		t.Helper()

		m, err := f.createRound(newClientRound("12th March 2019", "2019-03-19"))

		if err == nil {
			t.Fatalf("Test failed: expected %s, got nil", err.Error())
		}

		e := model.Round{}

		if *m != e {
			t.Fatalf("Test failed: expected %+v\n got %+v\n", e, m)
		}
	})
}

func TestUpdateRound(t *testing.T) {
	t.Run("updates an existing round struct", func(t *testing.T) {
		t.Helper()

		clientRound := newClientRound("2019-03-12", "2019-03-19")

		m, err := f.createRound(clientRound)

		if err != nil {
			t.Fatalf("Test failed: want nil, got %s", err.Error())
		}

		clock.Advance(30 * time.Minute)

		clientRound.Start = "2019-03-13"
		clientRound.End = "2019-03-20"

		updated, err := f.updateRound(clientRound, m)

		if err != nil {
			t.Fatalf("Test failed: want nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(54, updated.ID)
		a.Equal("2", updated.Name)
		a.Equal(45, updated.SeasonID)
		a.Equal("2019-03-13 00:00:00 +0000 UTC", updated.StartDate.String())
		a.Equal("2019-03-20 00:00:00 +0000 UTC", updated.EndDate.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", updated.CreatedAt.String())
		a.Equal("2019-01-14 11:55:00 +0000 UTC", updated.UpdatedAt.String())
	})

	t.Run("returns an error if date cannot be parsed", func(t *testing.T) {
		t.Helper()

		clientRound := newClientRound("2019-03-12", "2019-03-19")

		m, err := f.createRound(clientRound)

		if err != nil {
			t.Fatalf("Test failed: want nil, got %s", err.Error())
		}

		clientRound.Start = "13th March 2019"

		updated, err := f.updateRound(clientRound, m)

		if err == nil {
			t.Fatalf("Test failed: expected %s, got nil", err.Error())
		}

		e := model.Round{}

		if *updated != e {
			t.Fatalf("Test failed: expected %+v\n got %+v\n", e, m)
		}
	})
}

func newClientRound(start, end string) *sportmonks.Round {
	return &sportmonks.Round{
		ID:       54,
		Name:     2,
		LeagueID: 9801,
		SeasonID: 45,
		StageID:  3,
		Start:    start,
		End:      end,
	}
}