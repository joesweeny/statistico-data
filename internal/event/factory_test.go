package event

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

func TestCreateGoalEvent(t *testing.T) {
	t.Run("a new goal event struct is hydrated", func(t *testing.T) {
		t.Helper()

		m := f.createGoalEvent(newClientGoalEvent())

		a := assert.New(t)

		a.Equal(55, m.ID)
		a.Equal(1, m.TeamID)
		a.Equal(871, m.FixtureID)
		a.Equal(3029, m.PlayerID)
		a.Equal(22, *m.PlayerAssistID)
		a.Equal("1-0", m.Score)
		a.Equal(89, m.Minute)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
	})
}

func TestCreateSubstitutionEvent(t *testing.T) {
	t.Helper()

	m := f.createSubstitutionEvent(newClientSubstitutionEvent())

	a := assert.New(t)

	a.Equal(57, m.ID)
	a.Equal(29, m.TeamID)
	a.Equal(98102, m.FixtureID)
	a.Equal(23, m.PlayerInID)
	a.Equal(30, m.PlayerOutID)
	a.Equal(55, m.Minute)
	a.Nil(m.Injured)
	a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
}

func newClientGoalEvent() *sportmonks.GoalEvent {
	assistId := 22
	assistName := "Kyle Walker"
	return &sportmonks.GoalEvent{
		ID:               55,
		TeamID:           "1",
		Type:             "Goal",
		FixtureID:        871,
		PlayerID:         3029,
		PlayerName:       "Danilo",
		PlayerAssistID:   &assistId,
		PlayerAssistName: &assistName,
		Minute:           89,
		Result:           "1-0",
	}
}

func newClientSubstitutionEvent() *sportmonks.SubstitutionEvent {
	return &sportmonks.SubstitutionEvent{
		ID:            57,
		TeamID:        "29",
		Type:          "Sub",
		FixtureID:     98102,
		PlayerInID:    23,
		PlayerInName:  "Mark Noble",
		PlayerOutID:   30,
		PlayerOutName: "Michael Antonio",
		Minute:        55,
	}
}
