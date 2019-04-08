package stats

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/statistico/statistico-data/internal/model"
	"time"
)

var ErrNotFound = errors.New("not found")

type PostgresTeamStatsRepository struct {
	Connection *sql.DB
}

func (p *PostgresTeamStatsRepository) InsertTeamStats(m *model.TeamStats) error {
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
		m.TeamShots.Total,
		m.TeamShots.OnGoal,
		m.TeamShots.OffGoal,
		m.TeamShots.Blocked,
		m.TeamShots.InsideBox,
		m.TeamShots.OutsideBox,
		m.TeamPasses.Total,
		m.TeamPasses.Accuracy,
		m.TeamPasses.Percentage,
		m.TeamAttacks.Total,
		m.TeamAttacks.Dangerous,
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

func (p *PostgresTeamStatsRepository) UpdateTeamStats(m *model.TeamStats) error {
	if _, err := p.ByFixtureAndTeam(m.FixtureID, m.TeamID); err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_team_stats SET shots_total = $3, shots_on_goal = $4, shots_off_goal = $5, 
	shots_blocked = $6, shots_inside_box = $7, shots_outside_box = $8, passes_total = $9, passes_accuracy = $10, 
	passes_percentage = $11, attacks_total = $12, attacks_dangerous = $13, fouls = $14, corners = $15, offsides = $16, 
	possession = $17, yellow_cards = $18, red_cards = $19, saves = $20, substitutions = $21, goal_kicks = $22, 
	goal_attempts = $23, free_kicks = $24, throw_ins = $25, updated_at = $26 where fixture_id = $1 AND team_id = $2`

	_, err := p.Connection.Exec(
		query,
		m.FixtureID,
		m.TeamID,
		m.TeamShots.Total,
		m.TeamShots.OnGoal,
		m.TeamShots.OffGoal,
		m.TeamShots.Blocked,
		m.TeamShots.InsideBox,
		m.TeamShots.OutsideBox,
		m.TeamPasses.Total,
		m.TeamPasses.Accuracy,
		m.TeamPasses.Percentage,
		m.TeamAttacks.Total,
		m.TeamAttacks.Dangerous,
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
		&m.TeamShots.Total,
		&m.TeamShots.OnGoal,
		&m.TeamShots.OffGoal,
		&m.TeamShots.Blocked,
		&m.TeamShots.InsideBox,
		&m.TeamShots.OutsideBox,
		&m.TeamPasses.Total,
		&m.TeamPasses.Accuracy,
		&m.TeamPasses.Percentage,
		&m.TeamAttacks.Total,
		&m.TeamAttacks.Dangerous,
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
