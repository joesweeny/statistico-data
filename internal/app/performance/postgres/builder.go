package postgres

import (
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/statistico/statistico-football-data/internal/app/performance"
	"strings"
)

const (
	actionCombined = "combined"
	actionFor = "for"
	actionAgainst = "against"
	average = "average"
	away = "away"
	greaterThanEqualTo = "gte"
	home = "home"
	homeAway = "home_away"
	lessThanEqualTo = "lte"
	total = "total"
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
	if measure == total {
		if metric == greaterThanEqualTo {
			b = b.Where(sq.GtOrEq{stat: value})
		}

		if metric == lessThanEqualTo {
			b = b.Where(sq.LtOrEq{stat: value})
		}

		b = b.Having(sq.Eq{"COUNT(*)": games})
	}

	if measure == average {
		column := fmt.Sprintf("AVG(%s)", stat)

		if metric == greaterThanEqualTo {
			b = b.Having(sq.GtOrEq{column: value})
		}

		if metric == lessThanEqualTo {
			b = b.Having(sq.LtOrEq{column: value})
		}
	}

	return b
}

func buildSubSelect(stat, venue, action string, seasons []uint64) sq.SelectBuilder {
	if action == actionCombined && venue != homeAway {
		stat = fmt.Sprintf("SUM(%s) as %s", stat, stat)
	}

	b := sq.Select("team_id", "team_name", stat, "rank() over (partition by team_id order by date desc)")

	if action == actionCombined && venue == homeAway {
		return buildHomeAwayCombinedQuery(b, stat, seasons)
	}

	if venue == home {
		if action == actionFor {
			b = b.From("home_stats_for")
		}

		if action == actionAgainst {
			b = b.From("home_stats_against")
		}

		if action == actionCombined {
			b = b.FromSelect(buildUnionSubSelect(venue, action), "stats")
			b = b.GroupBy("team_id, team_name", "date")
		}
	}

	if venue == away {
		if action == actionFor {
			b = b.From("away_stats_for")
		}

		if action == actionAgainst {
			b = b.From("away_stats_against")
		}

		if action == actionCombined {
			b = b.FromSelect(buildUnionSubSelect(venue, action), "stats")
			b = b.GroupBy("team_id, team_name", "date")
		}
	}

	if venue == homeAway {
		if action == actionFor {
			b = b.FromSelect(buildUnionSubSelect(venue, action), "stats")
		}

		if action == actionAgainst {
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

	if venue == home {
		b = b.From("home_stats_for UNION SELECT * FROM home_stats_against")
	}

	if venue == away {
		b = b.From("away_stats_for UNION SELECT * FROM away_stats_against")
	}

	if venue == homeAway {
		if action == actionFor {
			b = b.From("home_stats_for UNION SELECT * FROM away_stats_for")
		}

		if action == actionAgainst {
			b = b.From("home_stats_against UNION SELECT * FROM away_stats_against")
		}
	}

	return b
}

func buildHomeAwayCombinedQuery(b sq.SelectBuilder, stat string, seasons []uint64) sq.SelectBuilder {
	stat = fmt.Sprintf("SUM(%s) as %s", stat, stat)

	homeStats := sq.Select("team_id", "team_name", stat, "date").
		FromSelect(buildUnionSubSelect(home, ""), "home_stats").
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