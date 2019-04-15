package player_stats

import (
	"github.com/statistico/statistico-data/internal/model"
	"github.com/statistico/statistico-data/internal/proto"
	pbPlayerStats "github.com/statistico/statistico-data/internal/proto/stats"
)

func HandlePlayerStats(p []*model.PlayerStats) ([]*pbPlayerStats.PlayerStats) {
	var stats []*pbPlayerStats.PlayerStats

	for _, player := range p {
		s := proto.PlayerStatsToProto(player)
		stats = append(stats, s)
	}

	return stats
}

func HandleLineupPlayers(p []*model.PlayerStats) ([]*pbPlayerStats.LineupPlayer) {
	var lineup []*pbPlayerStats.LineupPlayer

	for _, player := range p {
		l := proto.PlayerStatsToLineupPlayerProto(player)
		lineup = append(lineup, l)
	}

	return lineup
}
