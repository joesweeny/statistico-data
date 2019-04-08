package stats

import (
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
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
		PlayerShots:       *handlePlayerShots(&s.Stats.Shots),
		PlayerGoals:       *handlePlayerGoals(&s.Stats.Goals),
		PlayerFouls:       *handlePlayerFouls(&s.Stats.Fouls),
		YellowCards:       s.Stats.Cards.YellowCards,
		RedCard:           s.Stats.Cards.RedCards,
		PlayerCrosses:     *handlePlayerCrosses(&s.Stats.Passes),
		PlayerPasses:      *handlePlayerPasses(&s.Stats.Passes),
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
	m.PlayerShots = *handlePlayerShots(&s.Stats.Shots)
	m.PlayerGoals = *handlePlayerGoals(&s.Stats.Goals)
	m.PlayerFouls = *handlePlayerFouls(&s.Stats.Fouls)
	m.YellowCards = s.Stats.Cards.YellowCards
	m.RedCard = s.Stats.Cards.RedCards
	m.PlayerCrosses = *handlePlayerCrosses(&s.Stats.Passes)
	m.PlayerPasses = *handlePlayerPasses(&s.Stats.Passes)
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

func handlePlayerShots(s *sportmonks.PlayerShots) *model.PlayerShots {
	return &model.PlayerShots{
		Total:  s.ShotsTotal,
		OnGoal: s.ShotsOnGoal,
	}
}

func handlePlayerGoals(s *sportmonks.PlayerGoals) *model.PlayerGoals {
	return &model.PlayerGoals{
		Scored:   s.Scored,
		Conceded: s.Conceded,
	}
}

func handlePlayerFouls(s *sportmonks.PlayerFouls) *model.PlayerFouls {
	return &model.PlayerFouls{
		Drawn:     s.Drawn,
		Committed: s.Committed,
	}
}

func handlePlayerCrosses(s *sportmonks.PlayerPasses) *model.PlayerCrosses {
	return &model.PlayerCrosses{
		Total:    s.TotalCrosses,
		Accuracy: s.CrossesAccuracy,
	}
}

func handlePlayerPasses(s *sportmonks.PlayerPasses) *model.PlayerPasses {
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
