package player_stats

import (
	"github.com/statistico/statistico-data/internal/app"
	proto2 "github.com/statistico/statistico-data/internal/app/proto"
	"github.com/statistico/statistico-data/internal/proto"
)

func HandlePlayerStats(p []*app.PlayerStats) []*proto2.PlayerStats {
	var stats []*proto2.PlayerStats

	for _, player := range p {
		s := proto.PlayerStatsToProto(player)
		stats = append(stats, s)
	}

	return stats
}

func HandleStartingLineupPlayers(p []*app.PlayerStats) []*proto2.LineupPlayer {
	var lineup []*proto2.LineupPlayer

	for _, player := range p {
		if !player.IsSubstitute {
			l := proto.PlayerStatsToLineupPlayerProto(player)
			lineup = append(lineup, l)
		}
	}

	return lineup
}

func HandleSubstituteLineupPlayers(p []*app.PlayerStats) []*proto2.LineupPlayer {
	var lineup []*proto2.LineupPlayer

	for _, player := range p {
		if player.IsSubstitute {
			l := proto.PlayerStatsToLineupPlayerProto(player)
			lineup = append(lineup, l)
		}
	}

	return lineup
}
