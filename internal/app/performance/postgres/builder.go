package postgres

import (
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/statistico/statistico-data/internal/app/performance"
	"strings"
)

func buildTeamsQuery(s sq.StatementBuilderType, f *performance.StatFilter) sq.SelectBuilder {
	venue := f.Venue
	action := f.Action
	metric := f.Metric
	measure := f.Measure
	games := f.Games
	value := f.Value
	stat := f.Stat
	seasons := f.Seasons

	b := s.Select("team_id, team_name").
		FromSelect(buildSubSelect(stat, venue, action, seasons), "ranked").
		Where(sq.LtOrEq{"rank": games}).
		GroupBy("team_id, team_name")

	return parseWhereHavingClause(b, games, value, measure, stat, metric)
}

func parseWhereHavingClause(b sq.SelectBuilder, games uint8, value float32, measure, stat, metric string) sq.SelectBuilder {
	if measure == "total" {
		if metric == "gte" {
			b = b.Where(sq.GtOrEq{stat: value})
		}

		if metric == "lte" {
			b = b.Where(sq.LtOrEq{stat: value})
		}

		b = b.Having(sq.Eq{"COUNT(*)": games})
	}

	if measure == "average" {
		column := fmt.Sprintf("AVG(%s)", stat)

		if metric == "gte" {
			b = b.Having(sq.GtOrEq{column: value})
		}

		if metric == "lte" {
			b = b.Having(sq.LtOrEq{column: value})
		}
	}

	return b
}

func buildSubSelect(stat, venue, action string, seasons []uint64) sq.SelectBuilder {
	if action == "combined" && venue != "home_away" {
		stat = fmt.Sprintf("SUM(%s) as %s", stat, stat)
	}

	b := sq.Select("team_id", "team_name", stat, "rank() over (partition by team_id order by date desc)")

	if action == "combined" && venue == "home_away" {
		return buildHomeAwayCombinedQuery(b, stat, seasons)
	}

	if venue == "home" {
		if action == "for" {
			b = b.From("home_stats_for")
		}

		if action == "against" {
			b = b.From("home_stats_against")
		}

		if action == "combined" {
			b = b.FromSelect(buildUnionSubSelect(venue, action), "stats")
			b = b.GroupBy("team_id, team_name", "date")
		}
	}

	if venue == "away" {
		if action == "for" {
			b = b.From("away_stats_for")
		}

		if action == "against" {
			b = b.From("away_stats_against")
		}

		if action == "combined" {
			b = b.FromSelect(buildUnionSubSelect(venue, action), "stats")
			b = b.GroupBy("team_id, team_name", "date")
		}
	}

	if venue == "home_away" {
		if action == "for" {
			b = b.FromSelect(buildUnionSubSelect(venue, action), "stats")
		}

		if action == "against" {
			b = b.FromSelect(buildUnionSubSelect(venue, action), "stats")
		}
	}

	if len(seasons) > 0 {
		b = b.Where(sq.Eq{"season_id": seasons})
	}

	return b
}

func buildUnionSubSelect(venue, action string) sq.SelectBuilder {
	b := sq.Select("*")

	if venue == "home" {
		b = b.From("home_stats_for UNION SELECT * FROM home_stats_against")
	}

	if venue == "away" {
		b = b.From("away_stats_for UNION SELECT * FROM away_stats_against")
	}

	if venue == "home_away" {
		if action == "for" {
			b = b.From("home_stats_for UNION SELECT * FROM away_stats_for")
		}

		if action == "against" {
			b = b.From("home_stats_against UNION SELECT * FROM away_stats_against")
		}
	}

	return b
}

func buildHomeAwayCombinedQuery(b sq.SelectBuilder, stat string, seasons []uint64) sq.SelectBuilder {
	stat = fmt.Sprintf("SUM(%s) as %s", stat, stat)

	homeStats := sq.Select("team_id", "team_name", stat, "date").
		FromSelect(buildUnionSubSelect("home", ""), "home_stats").
		GroupBy("team_id, team_name", "date")

	awayStats := sq.Select("team_id", "team_name", stat, "date").
		FromSelect(buildUnionSubSelect("away", ""), "away_stats").
		GroupBy("team_id, team_name", "date")

	if len(seasons) > 0 {
		homeStats = homeStats.Where(fmt.Sprintf("season_id IN (%s)", parseSeasonsSlice(seasons)))
		awayStats = awayStats.Where(fmt.Sprintf("season_id IN (%s)", parseSeasonsSlice(seasons)))
	}

	homeSql, _, _ := homeStats.ToSql()
	awaySql, _, _ := awayStats.ToSql()

	combined := fmt.Sprintf("((%s UNION %s)) combined", homeSql, awaySql)

	return b.From(combined)
}

func parseSeasonsSlice(seasons []uint64) string {
	s, _ := json.Marshal(seasons)
	return strings.Trim(string(s), "[]")
}