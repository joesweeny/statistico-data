package postgres_test

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/statistico/statistico-data/internal/performance"
	"github.com/statistico/statistico-data/internal/performance/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildQuery(t *testing.T) {
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

		assertCorrectSql(t, &filter, sql)
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

		assertCorrectSql(t, &filter, sql)
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

		assertCorrectSql(t, &filter, sql)
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

		assertCorrectSql(t, &filter, sql)
	})
}

func assertCorrectSql(t *testing.T, filter *performance.StatFilter, expected string) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	b := builder.Select("team_id", "team_name")

	query := postgres.BuildQuery(b, filter)

	sql, _, err := query.ToSql()

	if err != nil {
		t.Fatalf("Expected nil, got %s", err.Error())
	}

	assert.Equal(t, expected, sql)
}

