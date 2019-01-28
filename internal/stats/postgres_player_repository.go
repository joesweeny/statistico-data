package stats

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
)

type PostgresPlayerStatsRepository struct {
	Connection *sql.DB
}

func (p *PostgresPlayerStatsRepository) Insert(m *model.PlayerStats) error {
	query := `INSERT INTO sportmonks_player_stats (fixture_id, player_id, team_id, position, formation_position, substitute,
	shots_total, shots_on_goal, goals_scored, goals_conceded, fouls_drawn, fouls_committed, yellow_cards, red_cards,
	crosses_total, crosses_accuracy, passes_total, passes_accuracy, assists, offsides, saves, pen_scored, pen_missed,
	pen_committed, pen_won, hit_woodwork, tackles, blocks, interceptions, clearances, minutes_played, created_at, updated_at)
	VALUES $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24,
	$25, $26, $27, $28, $29, $30, $31, $32, $33, $34`

	_, err := p.Connection.Exec(
		query,
		m.FixtureID,
		m.PLayerID,
		m.TeamID,
		m.Position,
		m.FormationPosition,
		m.IsSubstitute,
		m.PlayerShots.Total,
		m.PlayerShots.OnGoal,
		m.PlayerGoals.Scored,
		m.PlayerGoals.Conceded,
		m.PlayerFouls.Drawn,
		m.PlayerFouls.Committed,
		m.YellowCards,
		m.RedCard,
		m.PlayerPenalties.Scored,
		m.PlayerPenalties.Missed,
		m.PlayerPenalties.Committed,
		m.PlayerPenalties.Won,
		m.HitWoodwork,
		m.Tackles,
		m.Blocks,
		m.Interceptions,
		m.Clearances,
		m.MinutesPlayed,
		m.CreatedAt.Unix(),
		m.UpdatedAt.Unix(),
	)

	return err
}