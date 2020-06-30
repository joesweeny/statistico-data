package factory

import (
	"github.com/statistico/statistico-data/internal/app"
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

		proto := playerStatsToProto(stats)

		a := assert.New(t)
		a.Equal(uint64(77), proto.PlayerId)
		a.Equal(int32(5), proto.ShotsTotal.GetValue())
		a.Equal(int32(3), proto.ShotsOnGoal.GetValue())
		a.Equal(int32(2), proto.GoalsScored.GetValue())
		a.Equal(int32(0), proto.GoalsConceded.GetValue())
		a.Equal(int32(1), proto.Assists.GetValue())
	})

	t.Run("nullable fields are handled", func(t *testing.T) {
		stats := &app.PlayerStats{
			PlayerID: 77,
		}

		proto := playerStatsToProto(stats)

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

	player := app.PlayerStats{
		PlayerID:          105,
		Position:          &pos,
		IsSubstitute:      false,
		FormationPosition: &form,
	}

	t.Run("a new LineupPlayer proto struct is hydrated", func(t *testing.T) {
		pl := playerStatsToLineupPlayerProto(&player)

		a := assert.New(t)
		a.Equal(uint64(105), pl.PlayerId)
		a.Equal("M", pl.Position)
		a.False(pl.IsSubstitute)
		a.Equal(uint32(8), pl.FormationPosition.GetValue())
	})

	t.Run("nullable fields are handled", func(t *testing.T) {
		player.FormationPosition = nil

		pl := playerStatsToLineupPlayerProto(&player)

		a := assert.New(t)
		a.Equal(uint64(105), pl.PlayerId)
		a.Equal("M", pl.Position)
		a.False(pl.IsSubstitute)
		a.Equal(uint32(0), pl.FormationPosition.GetValue())
	})
}

func TestTeamToProtoTeam(t *testing.T) {
	t.Run("a new proto Team struct is hydrated", func(t *testing.T) {
		code := "WHU"
		founded := 1895
		logo := "https://logo.com"

		team := app.Team{
			ID:           1,
			Name:         "West Ham United",
			ShortCode:    &code,
			CountryID:    8,
			VenueID:      214,
			NationalTeam: false,
			Founded:      &founded,
			Logo:         &logo,
		}

		proto := TeamToProto(&team)

		a := assert.New(t)
		a.Equal(uint64(1), proto.Id)
		a.Equal("West Ham United", proto.Name)
		a.Equal("WHU", proto.ShortCode.Value)
		a.Equal(uint64(8), proto.CountryId)
		a.Equal(uint64(214), proto.VenueId)
		a.Equal(false, proto.IsNationalTeam)
		a.Equal(uint64(1895), proto.Founded.Value)
		a.Equal("https://logo.com", proto.Logo.Value)
	})

	t.Run("handle null values when converting team struct", func(t *testing.T) {
		team := app.Team{
			ID:           1,
			Name:         "West Ham United",
			CountryID:    8,
			VenueID:      214,
			NationalTeam: false,
		}

		proto := TeamToProto(&team)

		a := assert.New(t)
		a.Equal(uint64(1), proto.Id)
		a.Equal("West Ham United", proto.Name)
		a.Nil(proto.ShortCode)
		a.Equal(uint64(8), proto.CountryId)
		a.Equal(uint64(214), proto.VenueId)
		a.Equal(false, proto.IsNationalTeam)
		a.Nil(proto.Founded)
		a.Nil(proto.Logo)
	})
}

func TestTeamStatsToProto(t *testing.T) {
	t.Run("returns a proto team stats struct", func(t *testing.T) {
		m := newTeamStats()

		proto := teamStatsToProto(m)

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

func newPlayerStats(goals *int, assists *int, onGoal *int) *app.PlayerStats {
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

func newTeamStats() *app.TeamStats {
	total := 5
	blocked := 2
	corners := 15
	redCards := 1
	dangerous := 30
	goalKicks := 12

	return &app.TeamStats{
		TeamID: 850,
		TeamShots: app.TeamShots{
			Total:   &total,
			Blocked: &blocked,
		},
		Corners:  &corners,
		RedCards: &redCards,
		TeamAttacks: app.TeamAttacks{
			Dangerous: &dangerous,
		},
		GoalKicks: &goalKicks,
	}
}
