package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

type PlayerStatsRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (p *PlayerStatsRepository) Insert(a *app.PlayerStats) error {
	query := `INSERT INTO sportmonks_player_stats (fixture_id, player_id, team_id, position, formation_position, substitute,
	shots_total, shots_on_goal, goals_scored, goals_conceded, fouls_drawn, fouls_committed, yellow_cards, red_card,
	crosses_total, crosses_accuracy, passes_total, passes_accuracy, assists, offsides, saves, pen_scored, pen_missed, pen_saved,
	pen_committed, pen_won, hit_woodwork, tackles, blocks, interceptions, clearances, minutes_played, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24,
	$25, $26, $27, $28, $29, $30, $31, $32, $33, $34)`

	_, err := p.connection.Exec(
		query,
		a.FixtureID,
		a.PlayerID,
		a.TeamID,
		a.Position,
		a.FormationPosition,
		a.IsSubstitute,
		a.PlayerShots.Total,
		a.PlayerShots.OnGoal,
		a.PlayerGoals.Scored,
		a.PlayerGoals.Conceded,
		a.PlayerFouls.Drawn,
		a.PlayerFouls.Committed,
		a.YellowCards,
		a.RedCard,
		a.PlayerCrosses.Total,
		a.PlayerCrosses.Accuracy,
		a.PlayerPasses.Total,
		a.PlayerPasses.Accuracy,
		a.Assists,
		a.Offsides,
		a.Saves,
		a.PlayerPenalties.Scored,
		a.PlayerPenalties.Missed,
		a.PlayerPenalties.Saved,
		a.PlayerPenalties.Committed,
		a.PlayerPenalties.Won,
		a.HitWoodwork,
		a.Tackles,
		a.Blocks,
		a.Interceptions,
		a.Clearances,
		a.MinutesPlayed,
		p.clock.Now().Unix(),
		p.clock.Now().Unix(),
	)

	return err
}

func (p *PlayerStatsRepository) Update(a *app.PlayerStats) error {
	if _, err := p.ByFixtureAndPlayer(a.FixtureID, a.PlayerID); err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_player_stats SET position = $3, formation_position = $4, substitute = $5, shots_total = $6, 
	shots_on_goal = $7, goals_scored = $8, goals_conceded = $9, fouls_drawn = $10, fouls_committed = $11, yellow_cards = $12, 
	red_card = $13, crosses_total = $14, crosses_accuracy = $15, passes_total = $16, passes_accuracy = $17, assists = $18, 
	offsides = $19, saves = $20, pen_scored = $21, pen_missed = $22, pen_saved = $23, pen_committed = $24, pen_won = $25, 
	hit_woodwork = $26, tackles = $27, blocks = $28, interceptions = $29, clearances = $30, minutes_played = $31, 
	updated_at = $32 WHERE fixture_id = $1 AND player_id = $2`

	_, err := p.connection.Exec(
		query,
		a.FixtureID,
		a.PlayerID,
		a.Position,
		a.FormationPosition,
		a.IsSubstitute,
		a.PlayerShots.Total,
		a.PlayerShots.OnGoal,
		a.PlayerGoals.Scored,
		a.PlayerGoals.Conceded,
		a.PlayerFouls.Drawn,
		a.PlayerFouls.Committed,
		a.YellowCards,
		a.RedCard,
		a.PlayerCrosses.Total,
		a.PlayerCrosses.Accuracy,
		a.PlayerPasses.Total,
		a.PlayerPasses.Accuracy,
		a.Assists,
		a.Offsides,
		a.Saves,
		a.PlayerPenalties.Scored,
		a.PlayerPenalties.Missed,
		a.PlayerPenalties.Saved,
		a.PlayerPenalties.Committed,
		a.PlayerPenalties.Won,
		a.HitWoodwork,
		a.Tackles,
		a.Blocks,
		a.Interceptions,
		a.Clearances,
		a.MinutesPlayed,
		p.clock.Now().Unix(),
	)

	return err
}

func (p *PlayerStatsRepository) ByFixtureAndPlayer(fixtureID, playerID uint64) (*app.PlayerStats, error) {
	query := `SELECT * FROM sportmonks_player_stats WHERE fixture_id = $1 AND player_id = $2`
	row := p.connection.QueryRow(query, fixtureID, playerID)

	return rowToPlayerStats(row, fixtureID, playerID)
}

func (p *PlayerStatsRepository) ByFixtureAndTeam(fixtureID, teamID uint64) ([]*app.PlayerStats, error) {
	query := `SELECT * FROM sportmonks_player_stats WHERE fixture_id = $1 AND team_id = $2 order by formation_position ASC`
	rows, err := p.connection.Query(query, fixtureID, teamID)

	if err != nil {
		return []*app.PlayerStats{}, fmt.Errorf("player stats for Team ID %d and Fixture ID %d does not exist", teamID, fixtureID)
	}

	var (
		created int64
		updated int64
		stats   []*app.PlayerStats
	)

	defer rows.Close()

	for rows.Next() {
		var m app.PlayerStats

		err := rows.Scan(
			&m.FixtureID,
			&m.PlayerID,
			&m.TeamID,
			&m.Position,
			&m.FormationPosition,
			&m.IsSubstitute,
			&m.PlayerShots.Total,
			&m.PlayerShots.OnGoal,
			&m.PlayerGoals.Scored,
			&m.PlayerGoals.Conceded,
			&m.PlayerFouls.Drawn,
			&m.PlayerFouls.Committed,
			&m.YellowCards,
			&m.RedCard,
			&m.PlayerCrosses.Total,
			&m.PlayerCrosses.Accuracy,
			&m.PlayerPasses.Total,
			&m.PlayerPasses.Accuracy,
			&m.Assists,
			&m.Offsides,
			&m.Saves,
			&m.PlayerPenalties.Scored,
			&m.PlayerPenalties.Missed,
			&m.PlayerPenalties.Saved,
			&m.PlayerPenalties.Committed,
			&m.PlayerPenalties.Won,
			&m.HitWoodwork,
			&m.Tackles,
			&m.Blocks,
			&m.Interceptions,
			&m.Clearances,
			&m.MinutesPlayed,
			&created,
			&updated,
		)

		if err != nil {
			return stats, err
		}

		m.CreatedAt = time.Unix(created, 0)
		m.UpdatedAt = time.Unix(updated, 0)

		stats = append(stats, &m)
	}

	return stats, nil
}

func rowToPlayerStats(r *sql.Row, fixtureID, playerID uint64) (*app.PlayerStats, error) {
	var created int64
	var updated int64

	m := app.PlayerStats{}

	err := r.Scan(
		&m.FixtureID,
		&m.PlayerID,
		&m.TeamID,
		&m.Position,
		&m.FormationPosition,
		&m.IsSubstitute,
		&m.PlayerShots.Total,
		&m.PlayerShots.OnGoal,
		&m.PlayerGoals.Scored,
		&m.PlayerGoals.Conceded,
		&m.PlayerFouls.Drawn,
		&m.PlayerFouls.Committed,
		&m.YellowCards,
		&m.RedCard,
		&m.PlayerCrosses.Total,
		&m.PlayerCrosses.Accuracy,
		&m.PlayerPasses.Total,
		&m.PlayerPasses.Accuracy,
		&m.Assists,
		&m.Offsides,
		&m.Saves,
		&m.PlayerPenalties.Scored,
		&m.PlayerPenalties.Missed,
		&m.PlayerPenalties.Saved,
		&m.PlayerPenalties.Committed,
		&m.PlayerPenalties.Won,
		&m.HitWoodwork,
		&m.Tackles,
		&m.Blocks,
		&m.Interceptions,
		&m.Clearances,
		&m.MinutesPlayed,
		&created,
		&updated,
	)

	if err != nil {
		return &m, fmt.Errorf("player stats for Player ID %d and Fixture ID %d does not exist", playerID, fixtureID)
	}

	m.CreatedAt = time.Unix(created, 0)
	m.UpdatedAt = time.Unix(updated, 0)

	return &m, nil
}

func NewPlayerStatsRepository(connection *sql.DB, clock clockwork.Clock) *PlayerStatsRepository {
	return &PlayerStatsRepository{connection: connection, clock: clock}
}
