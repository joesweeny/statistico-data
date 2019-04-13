package proto

import (
	"github.com/statistico/statistico-data/internal/model"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPlayerStatsToProto(t *testing.T) {
	t.Run("a new PlayerStats proto struct is hydrated", func(t *testing.T) {
		var (
			goals = 2
			assists = 1
			onGoal = 3
		)

		stats := newPlayerStats(&goals, &assists, &onGoal)

		proto := PlayerStatsToProto(stats)

		a := assert.New(t)
		a.Equal(uint64(77), proto.PlayerId)
		a.Equal(int32(5), proto.ShotsTotal.GetValue())
		a.Equal(int32(3), proto.ShotsOnGoal.GetValue())
		a.Equal(int32(2), proto.GoalsScored.GetValue())
		a.Equal(int32(0), proto.GoalsConceded.GetValue())
		a.Equal(int32(1), proto.Assists.GetValue())
	})

	t.Run("nullable fields are handled", func(t *testing.T) {
		stats := &model.PlayerStats{
			PlayerID:        77,
		}

		proto := PlayerStatsToProto(stats)

		a := assert.New(t)
		a.Equal(uint64(77), proto.PlayerId)
		a.Equal(int32(0), proto.ShotsTotal.GetValue())
		a.Equal(int32(0), proto.ShotsOnGoal.GetValue())
		a.Equal(int32(0), proto.GoalsScored.GetValue())
		a.Equal(int32(0), proto.GoalsConceded.GetValue())
		a.Equal(int32(0), proto.Assists.GetValue())
	})
}

func newPlayerStats(goals *int, assists *int, onGoal *int) *model.PlayerStats {
	shots := 5
	conceded := 0
	return &model.PlayerStats{
		PlayerID:        77,
		PlayerShots:     model.PlayerShots{
			Total: 	&shots,
			OnGoal: onGoal,
		},
		PlayerGoals:     model.PlayerGoals{
			Scored: goals,
			Conceded: &conceded,
		},
		Assists: assists,
	}
}