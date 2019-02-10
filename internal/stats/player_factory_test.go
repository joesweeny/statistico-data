package stats

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var t = time.Date(2019, 01, 14, 11, 25, 00, 00, time.UTC)
var clock = clockwork.NewFakeClockAt(t)
var f = PlayerFactory{clock}

func TestCreatePlayerStats(t *testing.T) {
	t.Run("a new player stats struct is hydrated", func(t *testing.T) {
		t.Helper()

		m := f.createPlayerStats(newClientLineupPlayer(), false)

		a := assert.New(t)

		a.Equal(1203, m.FixtureID)
		a.Equal(20918, m.PlayerID)
		a.Equal(55, m.TeamID)
		a.Equal("MF", *m.Position)
		a.Nil(m.FormationPosition)
		a.False(m.IsSubstitute)
		a.Equal(10, *m.PlayerShots.Total)
		a.Equal(10, *m.PlayerShots.OnGoal)
		a.Nil(m.PlayerGoals.Scored)
		a.Nil(m.PlayerGoals.Conceded)
		a.Equal(0, *m.PlayerFouls.Drawn)
		a.Equal(0, *m.PlayerFouls.Committed)
		a.Nil(m.YellowCards)
		a.Nil(m.RedCard)
		a.Equal(0, *m.PlayerPenalties.Committed)
		a.Equal(0, *m.PlayerPenalties.Won)
		a.Equal(0, *m.PlayerPenalties.Scored)
		a.Equal(0, *m.PlayerPenalties.Saved)
		a.Equal(0, *m.PlayerPenalties.Missed)
		a.Equal(10, *m.PlayerCrosses.Total)
		a.Equal(10, *m.PlayerCrosses.Accuracy)
		a.Equal(0, *m.PlayerPasses.Total)
		a.Equal(0, *m.PlayerPasses.Accuracy)
		a.Equal(0, *m.Assists)
		a.Equal(0, *m.Offsides)
		a.Equal(0, *m.Saves)
		a.Equal(0, *m.HitWoodwork)
		a.Nil(m.Tackles)
		a.Nil(m.Blocks)
		a.Nil(m.Interceptions)
		a.Nil(m.Clearances)
		a.Nil(m.MinutesPlayed)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.UpdatedAt.String())
	})
}

func newClientLineupPlayer() *sportmonks.LineupPlayer {
	var i int
	pos := "MF"
	num := 4
	total := 10
	shots := sportmonks.PlayerShots{
		ShotsTotal:  &total,
		ShotsOnGoal: &total,
	}
	fouls := sportmonks.PlayerFouls{
		Drawn:     &i,
		Committed: &i,
	}
	passes := sportmonks.PlayerPasses{
		TotalCrosses:    &total,
		CrossesAccuracy: &total,
		Passes:          &i,
		PassesAccuracy:  &i,
	}
	extra := sportmonks.ExtraPlayerStats{
		Assists:      &i,
		Offsides:     &i,
		Saves:        &i,
		PenScored:    &i,
		PenMissed:    &i,
		PenSaved:     &i,
		PenCommitted: &i,
		PenWon:       &i,
		HitWoodwork:  &i,
	}
	stats := sportmonks.PlayerStats{
		Shots:             shots,
		Fouls:             fouls,
		Passes:            passes,
		ExtraPlayersStats: extra,
	}

	return &sportmonks.LineupPlayer{
		TeamID:     55,
		FixtureID:  1203,
		PlayerID:   20918,
		PlayerName: "Mark Noble",
		Number:     &num,
		Position:   &pos,
		Stats:      stats,
	}
}
