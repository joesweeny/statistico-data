package stats

import (
	"github.com/joesweeny/sportmonks-go-client"
	"testing"
	"github.com/stretchr/testify/assert"
)

var factory = TeamFactory{Clock: clock}

func TestCreateTeamStats(t *testing.T) {
	t.Run("a new team stats struct is hydrated", func(t *testing.T) {
		t.Helper()

		m := factory.createTeamStats(newClientTeamStats())

		a := assert.New(t)

		a.Equal(960, m.TeamID)
		a.Equal(34019, m.FixtureID)
		a.Equal(5, *m.TeamShots.Total)
		a.Equal(5, *m.TeamShots.OnGoal)
		a.Equal(5, *m.TeamShots.OffGoal)
		a.Equal(5, *m.TeamPasses.Total)
		a.Equal(5, *m.TeamPasses.Accuracy)
		a.Nil(m.TeamPasses.Percentage)
		a.Nil(m.TeamAttacks.Total)
		a.Nil(m.TeamAttacks.Dangerous)
		a.Equal(5, *m.Fouls)
		a.Equal(0, *m.Corners)
		a.Equal(5, *m.Offsides)
		a.Nil(m.Possession)
		a.Nil(m.YellowCards)
		a.Nil(m.RedCards)
		a.Equal(0, *m.Saves)
		a.Equal(0, *m.Substitutions)
		a.Nil(m.GoalKicks)
		a.Nil(m.GoalAttempts)
		a.Equal(5, *m.FreeKicks)
		a.Nil(m.ThrowIns)
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.CreatedAt.String())
		a.Equal("2019-01-14 11:25:00 +0000 UTC", m.UpdatedAt.String())
	})
}

func newClientTeamStats() *sportmonks.TeamStats {
	total := 5
	zero := 0
	shots := sportmonks.TeamShots{
		Total:   &total,
		Ongoal:  &total,
		Offgoal: &total,
	}

	passes := sportmonks.TeamPasses{
		Total:    &total,
		Accurate: &total,
	}

	return &sportmonks.TeamStats{
		TeamID:    960,
		FixtureID: 34019,
		Shots:     shots,
		Passes:    passes,
		Fouls: &total,
		Corners: &zero,
		Offsides: &total,
		Saves: &zero,
		Substitutions: &zero,
		FreeKick: &total,
	}
}