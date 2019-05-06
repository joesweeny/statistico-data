package team_stats

import (
	"github.com/statistico/statistico-data/internal/model"
	pbTeamStats "github.com/statistico/statistico-data/internal/proto/stats/team"
	"github.com/statistico/statistico-data/internal/proto"
)

func HandleTeamStats(t []*model.TeamStats) ([]*pbTeamStats.TeamStats) {
	var stats []*pbTeamStats.TeamStats

	for _, team := range t {
		s := proto.TeamStatsToProto(team)
		stats = append(stats, s)
	}

	return stats
}
