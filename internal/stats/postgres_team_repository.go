package stats

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type PostgresTeamStatsRepository struct {
	Connection *sql.DB
}

func (p *PostgresTeamStatsRepository) Insert(m *model.TeamStats) error {
	query := `
	INSERT INTO sportmonks_team_stats (fixture_id, team_id, shots_total, shots_on_goal, shots_off_goal, shots_blocked, 
	shots_inside_box, shots_outside_box, passes_total, passes_accuracy, passes_percentage, attacks_total, attacks_dangerous,
	fouls, corners, offsides, possession, yellow_cards, red_cards, saves, substitutions, goal_kicks, goal_attempts, 
	free_kicks, throw_ins, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
	$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)`

	_, err := p.Connection.Exec(
		query,
		m.FixtureID,
		m.TeamID,
		m.Shots.Total,
		m.Shots.OnGoal,
		m.Shots.OffGoal,
		m.Shots.Blocked,
		m.Shots.InsideBox,
		m.Shots.OutsideBox,
		m.Passes.Total,
		m.Passes.Accuracy,
		m.Passes.Percentage,
		m.Attacks.Total,
		m.Attacks.Dangerous,
		m.Fouls,
		m.Corners,
		m.Offsides,
		m.Possession,
		m.YellowCards,
		m.RedCards,
		m.Saves,
		m.Substitutions,
		m.GoalKicks,
		m.GoalAttempts,
		m.FreeKicks,
		m.ThrowIns,
		m.CreatedAt.Unix(),
		m.UpdatedAt.Unix(),
	)

	return err
}

func (p *PostgresTeamStatsRepository) ByFixtureAndTeam(fixtureId, teamId int) (*model.TeamStats, error) {
	query := `SELECT * FROM sportmonks_team_stats where fixture_id = $1 AND team_id = $2`
	row := p.Connection.QueryRow(query, fixtureId, teamId)

	return rowToStats(row)
}

func rowToStats(r *sql.Row) (*model.TeamStats, error) {
	var created int64
	var updated int64

	m := model.TeamStats{}

	err := r.Scan(
		&m.FixtureID,
		&m.TeamID,
		&m.Shots.Total,
		&m.Shots.OnGoal,
		&m.Shots.OffGoal,
		&m.Shots.Blocked,
		&m.Shots.InsideBox,
		&m.Shots.OutsideBox,
		&m.Passes.Total,
		&m.Passes.Accuracy,
		&m.Passes.Percentage,
		&m.Attacks.Total,
		&m.Attacks.Dangerous,
		&m.Fouls,
		&m.Corners,
		&m.Offsides,
		&m.Possession,
		&m.YellowCards,
		&m.RedCards,
		&m.Saves,
		&m.Substitutions,
		&m.GoalKicks,
		&m.GoalAttempts,
		&m.FreeKicks,
		&m.ThrowIns,
		&created,
		&updated,
	)

	if err != nil {
		return &m, ErrNotFound
	}

	m.CreatedAt = time.Unix(created, 0)
	m.UpdatedAt = time.Unix(updated, 0)

	return &m, nil
}