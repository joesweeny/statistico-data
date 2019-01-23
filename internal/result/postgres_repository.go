package result

import (
	"database/sql"
	"github.com/joesweeny/statshub/internal/model"
	_ "github.com/lib/pq"
)

type PostgresResultRepository struct {
	Connection *sql.DB
}

func (p *PostgresResultRepository) Insert(r *model.Result) error {
	query := `
	INSERT INTO sportmonks_result (fixture_id, pitch_condition, home_formation, away_formation, home_score, away_score,
	home_pen_score, away_pen_score, half_time_score, full_time_score, extra_time_score, home_league_position,
	away_league_position, minutes, seconds, added_time, extra_time, injury_time, created_at, updated_at) VALUES ($1, $2, 
	$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`

	_, err := p.Connection.Exec(
		query,
		r.FixtureID,
		r.PitchCondition,
		r.HomeFormation,
		r.AwayFormation,
		r.HomeScore,
		r.AwayScore,
		r.HomePenScore,
		r.AwayPenScore,
		r.HalfTimeScore,
		r.FullTimeScore,
		r.ExtraTimeScore,
		r.HomeLeaguePosition,
		r.AwayLeaguePosition,
		r.Minutes,
		r.Seconds,
		r.AddedTime,
		r.ExtraTime,
		r.InjuryTime,
		r.CreatedAt.Unix(),
		r.UpdatedAt.Unix(),
	)

	return err
}
