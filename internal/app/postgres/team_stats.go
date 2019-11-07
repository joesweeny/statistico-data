package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-data/internal/app"
	"time"
)

type TeamStatsRepository struct {
	connection *sql.DB
	clock      clockwork.Clock
}

func (t *TeamStatsRepository) InsertTeamStats(a *app.TeamStats) error {
	query := `
	INSERT INTO sportmonks_team_stats (fixture_id, team_id, shots_total, shots_on_goal, shots_off_goal, shots_blocked, 
	shots_inside_box, shots_outside_box, passes_total, passes_accuracy, passes_percentage, attacks_total, attacks_dangerous,
	fouls, corners, offsides, possession, yellow_cards, red_cards, saves, substitutions, goal_kicks, goal_attempts, 
	free_kicks, throw_ins, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
	$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)`

	_, err := t.connection.Exec(
		query,
		a.FixtureID,
		a.TeamID,
		a.TeamShots.Total,
		a.TeamShots.OnGoal,
		a.TeamShots.OffGoal,
		a.TeamShots.Blocked,
		a.TeamShots.InsideBox,
		a.TeamShots.OutsideBox,
		a.TeamPasses.Total,
		a.TeamPasses.Accuracy,
		a.TeamPasses.Percentage,
		a.TeamAttacks.Total,
		a.TeamAttacks.Dangerous,
		a.Fouls,
		a.Corners,
		a.Offsides,
		a.Possession,
		a.YellowCards,
		a.RedCards,
		a.Saves,
		a.Substitutions,
		a.GoalKicks,
		a.GoalAttempts,
		a.FreeKicks,
		a.ThrowIns,
		t.clock.Now().Unix(),
		t.clock.Now().Unix(),
	)

	return err
}

func (t *TeamStatsRepository) UpdateTeamStats(a *app.TeamStats) error {
	if _, err := t.ByFixtureAndTeam(a.FixtureID, a.TeamID); err != nil {
		return err
	}

	query := `
	UPDATE sportmonks_team_stats SET shots_total = $3, shots_on_goal = $4, shots_off_goal = $5, 
	shots_blocked = $6, shots_inside_box = $7, shots_outside_box = $8, passes_total = $9, passes_accuracy = $10, 
	passes_percentage = $11, attacks_total = $12, attacks_dangerous = $13, fouls = $14, corners = $15, offsides = $16, 
	possession = $17, yellow_cards = $18, red_cards = $19, saves = $20, substitutions = $21, goal_kicks = $22, 
	goal_attempts = $23, free_kicks = $24, throw_ins = $25, updated_at = $26 where fixture_id = $1 AND team_id = $2`

	_, err := t.connection.Exec(
		query,
		a.FixtureID,
		a.TeamID,
		a.TeamShots.Total,
		a.TeamShots.OnGoal,
		a.TeamShots.OffGoal,
		a.TeamShots.Blocked,
		a.TeamShots.InsideBox,
		a.TeamShots.OutsideBox,
		a.TeamPasses.Total,
		a.TeamPasses.Accuracy,
		a.TeamPasses.Percentage,
		a.TeamAttacks.Total,
		a.TeamAttacks.Dangerous,
		a.Fouls,
		a.Corners,
		a.Offsides,
		a.Possession,
		a.YellowCards,
		a.RedCards,
		a.Saves,
		a.Substitutions,
		a.GoalKicks,
		a.GoalAttempts,
		a.FreeKicks,
		a.ThrowIns,
		t.clock.Now().Unix(),
	)

	return err
}

func (t *TeamStatsRepository) ByFixtureAndTeam(fixtureID, teamID uint64) (*app.TeamStats, error) {
	query := `SELECT * FROM sportmonks_team_stats where fixture_id = $1 AND team_id = $2`
	row := t.connection.QueryRow(query, fixtureID, teamID)

	return rowToStats(row, fixtureID, teamID)
}

func rowToStats(r *sql.Row, fixtureID, teamID uint64) (*app.TeamStats, error) {
	var created int64
	var updated int64

	var a = app.TeamStats{}

	err := r.Scan(
		&a.FixtureID,
		&a.TeamID,
		&a.TeamShots.Total,
		&a.TeamShots.OnGoal,
		&a.TeamShots.OffGoal,
		&a.TeamShots.Blocked,
		&a.TeamShots.InsideBox,
		&a.TeamShots.OutsideBox,
		&a.TeamPasses.Total,
		&a.TeamPasses.Accuracy,
		&a.TeamPasses.Percentage,
		&a.TeamAttacks.Total,
		&a.TeamAttacks.Dangerous,
		&a.Fouls,
		&a.Corners,
		&a.Offsides,
		&a.Possession,
		&a.YellowCards,
		&a.RedCards,
		&a.Saves,
		&a.Substitutions,
		&a.GoalKicks,
		&a.GoalAttempts,
		&a.FreeKicks,
		&a.ThrowIns,
		&created,
		&updated,
	)

	if err != nil {
		return &a, fmt.Errorf("stats for Team ID %d and Fixture ID %d does not exist", teamID, fixtureID)
	}

	a.CreatedAt = time.Unix(created, 0)
	a.UpdatedAt = time.Unix(updated, 0)

	return &a, nil
}

func NewTeamStatsRepository(connection *sql.DB, clock clockwork.Clock) *TeamStatsRepository {
	return &TeamStatsRepository{connection: connection, clock: clock}
}