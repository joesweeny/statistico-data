package stats

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/model"
	"github.com/jonboulle/clockwork"
)

type PlayerFactory struct {
	Clock clockwork.Clock
}

func (f PlayerFactory) createPlayerStats(s *sportmonks.LineupPlayer, sub bool) *model.PlayerStats {
	return &model.PlayerStats{
		FixtureID:         s.FixtureID,
		PlayerID:          s.PlayerID,
		TeamID:            s.TeamID,
		Position:          s.Position,
		FormationPosition: s.FormationPosition,
		IsSubstitute:      sub,
		PlayerShots:       *handleShots(&s.Stats.Shots),
		PlayerGoals:       *handleGoals(&s.Stats.Goals),
		PlayerFouls:       *handleFouls(&s.Stats.Fouls),
		YellowCards:       s.Stats.Cards.YellowCards,
		RedCard:           s.Stats.Cards.RedCards,
		PlayerCrosses:     *handleCrosses(&s.Stats.Passes),
		PlayerPasses:      *handlePasses(&s.Stats.Passes),
		Assists:           s.Stats.ExtraPlayersStats.Assists,
		Offsides:          s.Stats.ExtraPlayersStats.Offsides,
		Saves:             s.Stats.ExtraPlayersStats.Saves,
		PlayerPenalties:   *handlePenalties(&s.Stats.ExtraPlayersStats),
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

func (f PlayerFactory) updatePlayerStats(s *sportmonks.LineupPlayer, m *model.PlayerStats) *model.PlayerStats {
	m.Position = s.Position
	m.FormationPosition = s.FormationPosition
	m.PlayerShots = *handleShots(&s.Stats.Shots)
	m.PlayerGoals = *handleGoals(&s.Stats.Goals)
	m.PlayerFouls = *handleFouls(&s.Stats.Fouls)
	m.YellowCards = s.Stats.Cards.YellowCards
	m.RedCard = s.Stats.Cards.RedCards
	m.PlayerCrosses = *handleCrosses(&s.Stats.Passes)
	m.PlayerPasses = *handlePasses(&s.Stats.Passes)
	m.Assists = s.Stats.ExtraPlayersStats.Assists
	m.Offsides = s.Stats.ExtraPlayersStats.Offsides
	m.Saves = s.Stats.ExtraPlayersStats.Saves
	m.PlayerPenalties = *handlePenalties(&s.Stats.ExtraPlayersStats)
	m.HitWoodwork = s.Stats.ExtraPlayersStats.HitWoodwork
	m.Tackles = s.Stats.ExtraPlayersStats.Tackles
	m.Blocks = s.Stats.ExtraPlayersStats.Blocks
	m.Interceptions = s.Stats.ExtraPlayersStats.Interceptions
	m.Clearances = s.Stats.ExtraPlayersStats.Clearances
	m.MinutesPlayed = s.Stats.ExtraPlayersStats.MinutesPlayed
	m.UpdatedAt = f.Clock.Now()

	return m
}

func handleShots(s *sportmonks.PlayerShots) *model.PlayerShots {
	return &model.PlayerShots{
		Total:  s.ShotsTotal,
		OnGoal: s.ShotsOnGoal,
	}
}

func handleGoals(s *sportmonks.PlayerGoals) *model.PlayerGoals {
	return &model.PlayerGoals{
		Scored:   s.Scored,
		Conceded: s.Conceded,
	}
}

func handleFouls(s *sportmonks.PlayerFouls) *model.PlayerFouls {
	return &model.PlayerFouls{
		Drawn:     s.Drawn,
		Committed: s.Committed,
	}
}

func handleCrosses(s *sportmonks.PlayerPasses) *model.PlayerCrosses {
	return &model.PlayerCrosses{
		Total:    s.TotalCrosses,
		Accuracy: s.CrossesAccuracy,
	}
}

func handlePasses(s *sportmonks.PlayerPasses) *model.PlayerPasses {
	return &model.PlayerPasses{
		Total:    s.Passes,
		Accuracy: s.PassesAccuracy,
	}
}

func handlePenalties(s *sportmonks.ExtraPlayerStats) *model.PlayerPenalties {
	return &model.PlayerPenalties{
		Scored:    s.PenScored,
		Missed:    s.PenMissed,
		Saved:     s.PenSaved,
		Committed: s.PenCommitted,
		Won:       s.PenWon,
	}
}
