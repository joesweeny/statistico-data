package factory

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlePlayerStats(t *testing.T) {
	t.Run("returns a slice of proto PlayerStats structs", func(t *testing.T) {
		var (
			goals   = 1
			assists = 1
			onGoal  = 2
		)

		x := []*app.PlayerStats{
			modelPlayerStats(&goals, &assists, &onGoal),
			modelPlayerStats(&goals, &assists, &onGoal),
			modelPlayerStats(&goals, &assists, &onGoal),
		}

		stats := handlePlayerStats(x)

		assert.Equal(t, 3, len(stats))
	})
}

func TestHandleStartingLineupPlayers(t *testing.T) {
	t.Run("returns a slice of proto LineupPlayer structs who are not substitutes", func(t *testing.T) {
		var (
			playerId1  = 1
			formation1 = 1
			pos1       = "M"
		)

		var (
			playerId2  = 2
			formation2 = 2
			pos2       = "M"
		)

		var (
			playerId3  = 3
			formation3 = 3
			pos3       = "M"
		)

		x := []*app.PlayerStats{
			modelPlayerLineup(playerId1, &formation1, &pos1, false),
			modelPlayerLineup(playerId2, &formation2, &pos2, false),
			modelPlayerLineup(playerId3, &formation3, &pos3, true),
		}

		lineup := handleStartingLineupPlayers(x)

		a := assert.New(t)

		a.Equal(2, len(lineup))

		for i, l := range lineup {
			a.Equal(uint64(i+1), l.PlayerId)
			a.Equal("M", l.Position)
			a.Equal(uint32(i+1), l.FormationPosition.GetValue())
			a.False(l.IsSubstitute)
		}
	})
}

func TestHandleSubstituteLineupPlayers(t *testing.T) {
	t.Run("returns a slice of proto LineupPlayer structs who are substitutes", func(t *testing.T) {
		var (
			playerId1  = 1
			formation1 = 1
			pos1       = "M"
		)

		var (
			playerId2  = 2
			formation2 = 2
			pos2       = "M"
		)

		var (
			playerId3  = 3
			formation3 = 3
			pos3       = "M"
		)

		x := []*app.PlayerStats{
			modelPlayerLineup(playerId1, &formation1, &pos1, true),
			modelPlayerLineup(playerId2, &formation2, &pos2, false),
			modelPlayerLineup(playerId3, &formation3, &pos3, false),
		}

		lineup := handleSubstituteLineupPlayers(x)

		a := assert.New(t)

		a.Equal(1, len(lineup))

		for i, l := range lineup {
			a.Equal(uint64(i+1), l.PlayerId)
			a.Equal("M", l.Position)
			a.Equal(uint32(i+1), l.FormationPosition.GetValue())
			a.True(l.IsSubstitute)
		}
	})
}

func modelPlayerStats(goals *int, assists *int, onGoal *int) *app.PlayerStats {
	shots := 5
	conceded := 0
	return &app.PlayerStats{
		PlayerID: 77,
		PlayerShots: app.PlayerShots{
			Total:  &shots,
			OnGoal: onGoal,
		},
		PlayerGoals: app.PlayerGoals{
			Scored:   goals,
			Conceded: &conceded,
		},
		Assists: assists,
	}
}

func modelPlayerLineup(playerId int, formation *int, position *string, sub bool) *app.PlayerStats {
	return &app.PlayerStats{
		PlayerID:          uint64(playerId),
		Position:          position,
		IsSubstitute:      sub,
		FormationPosition: formation,
	}
}
