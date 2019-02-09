package stats

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/jonboulle/clockwork"
)

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createPlayerStats(s *sportmonks.LineupPlayer, sub bool) *model.PlayerStats {
	shots := model.PlayerShots{
		Total:  s.Stats.Shots.ShotsTotal,
		OnGoal: s.Stats.Shots.ShotsOnGoal,
	}

	goals := model.PlayerGoals{
		Scored:   s.Stats.Goals.Scored,
		Conceded: s.Stats.Goals.Conceded,
	}

	fouls := model.PlayerFouls{
		Drawn:     s.Stats.Fouls.Drawn,
		Committed: s.Stats.Fouls.Committed,
	}

	crosses := model.PlayerCrosses{
		Total:    s.Stats.Passes.TotalCrosses,
		Accuracy: s.Stats.Passes.CrossesAccuracy,
	}

	passes := model.PlayerPasses{
		Total:    s.Stats.Passes.Passes,
		Accuracy: s.Stats.Passes.PassesAccuracy,
	}

	penalties := model.PlayerPenalties{
		Scored:    s.Stats.ExtraPlayersStats.PenScored,
		Missed:    s.Stats.ExtraPlayersStats.PenMissed,
		Saved:     s.Stats.ExtraPlayersStats.PenSaved,
		Committed: s.Stats.ExtraPlayersStats.PenCommitted,
		Won:       s.Stats.ExtraPlayersStats.PenWon,
	}

	return &model.PlayerStats{
		FixtureID:         s.FixtureID,
		PlayerID:          s.PlayerID,
		TeamID:            s.TeamID,
		Position:          s.Position,
		FormationPosition: s.FormationPosition,
		IsSubstitute:      sub,
		PlayerShots:       shots,
		PlayerGoals:       goals,
		PlayerFouls:       fouls,
		YellowCards:       s.Stats.Cards.YellowCards,
		RedCard:           s.Stats.Cards.RedCards,
		PlayerCrosses:     crosses,
		PlayerPasses:      passes,
		Assists:           s.Stats.ExtraPlayersStats.Assists,
		Offsides:          s.Stats.ExtraPlayersStats.Offsides,
		Saves:             s.Stats.ExtraPlayersStats.Saves,
		PlayerPenalties:   penalties,
		HitWoodwork:       s.Stats.ExtraPlayersStats.HitWoodwork,
		Tackles:           s.Stats.ExtraPlayersStats.Tackles,
		Blocks:            s.Stats.ExtraPlayersStats.Blocks,
		Interceptions:     s.Stats.ExtraPlayersStats.Interceptions,
		Clearances:        s.Stats.ExtraPlayersStats.Clearances,
		MinutesPlayed:     s.Stats.ExtraPlayersStats.MinutesPlayed,
		CreatedAt:         f.Clock.Now(),
		UpdatedAt:         f.Clock.Now(),
	}
}
