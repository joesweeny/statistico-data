package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/statistico/statistico-data/internal/app/performance"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildTeamsQuery(t *testing.T) {
	t.Run("builds query for goals scored for home venue greater than total", func(t *testing.T) {
		t.Helper()

		filter := performance.StatFilter{
			Seasons: []uint64{43, 992},
			Stat:    "goals",
			Action:  "for",
			Metric:  "gte",
			Measure: "total",
			Value:   2.0,
			Venue:   "home",
			Games:   5,
		}

		sql := "SELECT team_id, team_name " +
			"FROM (SELECT team_id, team_name, goals, rank() over (partition by team_id order by date desc) " +
			"FROM home_stats_for WHERE season_id IN ($1,$2)) AS ranked " +
			"WHERE rank <= $3 AND goals >= $4 " +
			"GROUP BY team_id, team_name " +
			"HAVING COUNT(*) = $5"

		bindings := []interface{}{uint64(43), uint64(992), uint8(5), float32(2), uint8(5)}

		assertCorrectSql(t, &filter, sql, bindings)
	})

	t.Run("builds query for goals scored for away venue less than total", func(t *testing.T) {
		t.Helper()

		filter := performance.StatFilter{
			Seasons: []uint64{43, 992},
			Stat:    "goals",
			Action:  "for",
			Metric:  "lte",
			Measure: "total",
			Value:   1.0,
			Venue:   "away",
			Games:   3,
		}

		sql := "SELECT team_id, team_name " +
			"FROM (SELECT team_id, team_name, goals, rank() over (partition by team_id order by date desc) " +
			"FROM away_stats_for WHERE season_id IN ($1,$2)) AS ranked " +
			"WHERE rank <= $3 AND goals <= $4 " +
			"GROUP BY team_id, team_name " +
			"HAVING COUNT(*) = $5"

		bindings := []interface{}{uint64(43), uint64(992), uint8(3), float32(1), uint8(3)}

		assertCorrectSql(t, &filter, sql, bindings)
	})

	t.Run("builds query for corners against away venue greater than average", func(t *testing.T) {
		t.Helper()

		filter := performance.StatFilter{
			Seasons: []uint64{},
			Stat:    "corners",
			Action:  "against",
			Metric:  "gte",
			Measure: "average",
			Value:   1.0,
			Venue:   "away",
			Games:   3,
		}

		sql := "SELECT team_id, team_name " +
			"FROM (SELECT team_id, team_name, corners, rank() over (partition by team_id order by date desc) " +
			"FROM away_stats_against) AS ranked " +
			"WHERE rank <= $1 " +
			"GROUP BY team_id, team_name " +
			"HAVING AVG(corners) >= $2"

		bindings := []interface{}{uint8(3), float32(1)}

		assertCorrectSql(t, &filter, sql, bindings)
	})

	t.Run("builds query for possession against home venue less than average", func(t *testing.T) {
		t.Helper()

		filter := performance.StatFilter{
			Seasons: []uint64{},
			Stat:    "possession",
			Action:  "against",
			Metric:  "lte",
			Measure: "average",
			Value:   1.0,
			Venue:   "home",
			Games:   3,
		}

		sql := "SELECT team_id, team_name " +
			"FROM (SELECT team_id, team_name, possession, rank() over (partition by team_id order by date desc) " +
			"FROM home_stats_against) AS ranked " +
			"WHERE rank <= $1 " +
			"GROUP BY team_id, team_name " +
			"HAVING AVG(possession) <= $2"

		bindings := []interface{}{uint8(3), float32(1)}

		assertCorrectSql(t, &filter, sql, bindings)
	})

	t.Run("builds query for shots on goal combined at home venue greater than total", func(t *testing.T) {
		t.Helper()

		filter := performance.StatFilter{
			Seasons: []uint64{45, 222},
			Stat:    "shots_on_goal",
			Action:  "combined",
			Metric:  "gte",
			Measure: "total",
			Value:   5,
			Venue:   "home",
			Games:   3,
		}

		sql := "SELECT team_id, team_name " +
			"FROM (SELECT team_id, team_name, SUM(shots_on_goal) as shots_on_goal, rank() over (partition by team_id order by date desc) " +
			"FROM (SELECT * FROM home_stats_for UNION SELECT * FROM home_stats_against) AS stats " +
			"WHERE season_id IN ($1,$2) GROUP BY team_id, team_name, date) AS ranked " +
			"WHERE rank <= $3 AND shots_on_goal >= $4 " +
			"GROUP BY team_id, team_name " +
			"HAVING COUNT(*) = $5"

		bindings := []interface{}{uint64(45), uint64(222), uint8(3), float32(5), uint8(3)}

		assertCorrectSql(t, &filter, sql, bindings)
	})

	t.Run("builds query for shots total combined at away venue greater than average", func(t *testing.T) {
		t.Helper()

		filter := performance.StatFilter{
			Seasons: []uint64{45, 222},
			Stat:    "shots_total",
			Action:  "combined",
			Metric:  "gte",
			Measure: "average",
			Value:   5,
			Venue:   "away",
			Games:   3,
		}

		sql := "SELECT team_id, team_name " +
			"FROM (SELECT team_id, team_name, SUM(shots_total) as shots_total, rank() over (partition by team_id order by date desc) " +
			"FROM (SELECT * FROM away_stats_for UNION SELECT * FROM away_stats_against) AS stats " +
			"WHERE season_id IN ($1,$2) GROUP BY team_id, team_name, date) AS ranked " +
			"WHERE rank <= $3 " +
			"GROUP BY team_id, team_name " +
			"HAVING AVG(shots_total) >= $4"

		bindings := []interface{}{uint64(45), uint64(222), uint8(3), float32(5)}

		assertCorrectSql(t, &filter, sql, bindings)
	})

	t.Run("builds query for xg combined at home venue greater than total", func(t *testing.T) {
		t.Helper()

		filter := performance.StatFilter{
			Seasons: []uint64{45, 222},
			Stat:    "xg",
			Action:  "combined",
			Metric:  "gte",
			Measure: "total",
			Value:   5,
			Venue:   "home",
			Games:   3,
		}

		sql := "SELECT team_id, team_name " +
			"FROM (SELECT team_id, team_name, SUM(xg) as xg, rank() over (partition by team_id order by date desc) " +
			"FROM (SELECT * FROM home_stats_for UNION SELECT * FROM home_stats_against) AS stats " +
			"WHERE season_id IN ($1,$2) GROUP BY team_id, team_name, date) AS ranked " +
			"WHERE rank <= $3 AND xg >= $4 " +
			"GROUP BY team_id, team_name " +
			"HAVING COUNT(*) = $5"

		bindings := []interface{}{uint64(45), uint64(222), uint8(3), float32(5), uint8(3)}

		assertCorrectSql(t, &filter, sql, bindings)
	})

	t.Run("builds query for goals for at home and away greater than total", func(t *testing.T) {
		t.Helper()

		filter := performance.StatFilter{
			Seasons: []uint64{45, 222},
			Stat:    "goals",
			Action:  "for",
			Metric:  "gte",
			Measure: "total",
			Value:   5,
			Venue:   "home_away",
			Games:   3,
		}

		sql := "SELECT team_id, team_name " +
			"FROM (SELECT team_id, team_name, goals, rank() over (partition by team_id order by date desc) " +
			"FROM (SELECT * FROM home_stats_for UNION SELECT * FROM away_stats_for) AS stats " +
			"WHERE season_id IN ($1,$2)) AS ranked " +
			"WHERE rank <= $3 AND goals >= $4 " +
			"GROUP BY team_id, team_name " +
			"HAVING COUNT(*) = $5"

		bindings := []interface{}{uint64(45), uint64(222), uint8(3), float32(5), uint8(3)}

		assertCorrectSql(t, &filter, sql, bindings)
	})

	t.Run("builds query for goals against at home and away less than average", func(t *testing.T) {
		t.Helper()

		filter := performance.StatFilter{
			Seasons: []uint64{},
			Stat:    "goals",
			Action:  "against",
			Metric:  "lte",
			Measure: "average",
			Value:   5,
			Venue:   "home_away",
			Games:   3,
		}

		sql := "SELECT team_id, team_name " +
			"FROM (SELECT team_id, team_name, goals, rank() over (partition by team_id order by date desc) " +
			"FROM (SELECT * FROM home_stats_against UNION SELECT * FROM away_stats_against) AS stats) AS ranked " +
			"WHERE rank <= $1 " +
			"GROUP BY team_id, team_name " +
			"HAVING AVG(goals) <= $2"

		bindings := []interface{}{uint8(3), float32(5)}

		assertCorrectSql(t, &filter, sql, bindings)
	})

	t.Run("builds query for xg combined at home and away greater than total", func(t *testing.T) {
		t.Helper()

		filter := performance.StatFilter{
			Seasons: []uint64{16036,5},
			Stat:    "xg",
			Action:  "combined",
			Metric:  "gte",
			Measure: "total",
			Value:   5,
			Venue:   "home_away",
			Games:   3,
		}

		sql := "SELECT team_id, team_name " +
			"FROM (SELECT team_id, team_name, xg, rank() over (partition by team_id order by date desc) " +
			"FROM ((SELECT team_id, team_name, SUM(xg) as xg, date " +
			"FROM (SELECT * FROM home_stats_for UNION SELECT * FROM home_stats_against) AS home_stats " +
			"WHERE season_id IN (16036,5) " +
			"GROUP BY team_id, team_name, date " +
			"UNION " +
			"SELECT team_id, team_name, SUM(xg) as xg, date " +
			"FROM (SELECT * FROM away_stats_for UNION SELECT * FROM away_stats_against) AS away_stats " +
			"WHERE season_id IN (16036,5) " +
			"GROUP BY team_id, team_name, date)) combined) AS ranked " +
			"WHERE rank <= $1 AND xg >= $2 " +
			"GROUP BY team_id, team_name " +
			"HAVING COUNT(*) = $3"

		bindings := []interface{}{uint8(3), float32(5), uint8(3)}

		assertCorrectSql(t, &filter, sql, bindings)
	})
}

func assertCorrectSql(t *testing.T, filter *performance.StatFilter, expected string, bindings []interface{}) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := buildTeamsQuery(builder, filter)

	sql, args, err := query.ToSql()

	if err != nil {
		t.Fatalf("Expected nil, got %s", err.Error())
	}

	assert.Equal(t, expected, sql)
	assert.Equal(t, bindings, args)
}

