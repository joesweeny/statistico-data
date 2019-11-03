package proto

import (
	"github.com/statistico/statistico-data/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerStatsToProto(t *testing.T) {
	t.Run("a new PlayerStats proto struct is hydrated", func(t *testing.T) {
		var (
			goals   = 2
			assists = 1
			onGoal  = 3
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
			PlayerID: 77,
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

func TestPlayerStatsToLineupPlayerProto(t *testing.T) {
	pos := "M"
	form := 8

	player := model.PlayerStats{
		PlayerID:          105,
		Position:          &pos,
		IsSubstitute:      false,
		FormationPosition: &form,
	}

	t.Run("a new LineupPlayer proto struct is hydrated", func(t *testing.T) {
		pl := PlayerStatsToLineupPlayerProto(&player)

		a := assert.New(t)
		a.Equal(uint64(105), pl.PlayerId)
		a.Equal("M", pl.Position)
		a.False(pl.IsSubstitute)
		a.Equal(uint32(8), pl.FormationPosition.GetValue())
	})

	t.Run("nullable fields are handled", func(t *testing.T) {
		player.FormationPosition = nil

		pl := PlayerStatsToLineupPlayerProto(&player)

		a := assert.New(t)
		a.Equal(uint64(105), pl.PlayerId)
		a.Equal("M", pl.Position)
		a.False(pl.IsSubstitute)
		a.Equal(uint32(0), pl.FormationPosition.GetValue())
	})
}

func TestTeamStatsToProto(t *testing.T) {
	t.Run("returns a proto team stats struct", func(t *testing.T) {
		m := newTeamStats()

		proto := TeamStatsToProto(m)

		a := assert.New(t)
		a.Equal(uint64(850), proto.TeamId)
		a.Equal(uint32(5), proto.ShotsTotal.GetValue())
		a.Nil(proto.ShotsOnGoal)
		a.Nil(proto.ShotsOffGoal)
		a.Equal(uint32(2), proto.ShotsBlocked.GetValue())
		a.Nil(proto.ShotsInsideBox)
		a.Nil(proto.ShotsOutsideBox)
		a.Nil(proto.PassesTotal)
		a.Nil(proto.PassesAccuracy)
		a.Nil(proto.PassesPercentage)
		a.Nil(proto.AttacksTotal)
		a.Equal(uint32(30), proto.AttacksDangerous.GetValue())
		a.Nil(proto.Fouls)
		a.Equal(uint32(15), proto.Corners.GetValue())
		a.Nil(proto.Offsides)
		a.Nil(proto.Possession)
		a.Nil(proto.YellowCards)
		a.Equal(uint32(1), proto.RedCards.GetValue())
		a.Nil(proto.Saves)
		a.Nil(proto.Substitutions)
		a.Equal(uint32(12), proto.GoalKicks.GetValue())
		a.Nil(proto.GoalAttempts)
		a.Nil(proto.FreeKicks)
		a.Nil(proto.ThrowIns)
	})
}

func newPlayerStats(goals *int, assists *int, onGoal *int) *model.PlayerStats {
	shots := 5
	conceded := 0
	return &model.PlayerStats{
		PlayerID: 77,
		PlayerShots: model.PlayerShots{
			Total:  &shots,
			OnGoal: onGoal,
		},
		PlayerGoals: model.PlayerGoals{
			Scored:   goals,
			Conceded: &conceded,
		},
		Assists: assists,
	}
}

func newTeamStats() *model.TeamStats {
	total := 5
	blocked := 2
	corners := 15
	redCards := 1
	dangerous := 30
	goalKicks := 12

	return &model.TeamStats{
		TeamID: 850,
		TeamShots: model.TeamShots{
			Total:   &total,
			Blocked: &blocked,
		},
		Corners:  &corners,
		RedCards: &redCards,
		TeamAttacks: model.TeamAttacks{
			Dangerous: &dangerous,
		},
		GoalKicks: &goalKicks,
	}
}
