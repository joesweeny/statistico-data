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

func HandleStartingLineupPlayers(p []*model.PlayerStats) ([]*pbPlayerStats.LineupPlayer) {
	var lineup []*pbPlayerStats.LineupPlayer

	for _, player := range p {
		if player.IsSubstitute == false {
			l := proto.PlayerStatsToLineupPlayerProto(player)
			lineup = append(lineup, l)
		}
	}

	return lineup
}

func HandleSubstituteLineupPlayers(p []*model.PlayerStats) ([]*pbPlayerStats.LineupPlayer) {
	var lineup []*pbPlayerStats.LineupPlayer

	for _, player := range p {
		if player.IsSubstitute == true {
			l := proto.PlayerStatsToLineupPlayerProto(player)
			lineup = append(lineup, l)
		}
	}

	return lineup
}
