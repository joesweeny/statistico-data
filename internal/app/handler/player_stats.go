package handler

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/converter"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
)

func HandlePlayerStats(p []*app.PlayerStats) []*proto.PlayerStats {
	var stats []*proto.PlayerStats

	for _, player := range p {
		s := converter.PlayerStatsToProto(player)
		stats = append(stats, s)
	}

	return stats
}

func HandleStartingLineupPlayers(p []*app.PlayerStats) []*proto.LineupPlayer {
	var lineup []*proto.LineupPlayer

	for _, player := range p {
		if !player.IsSubstitute {
			l := converter.PlayerStatsToLineupPlayerProto(player)
			lineup = append(lineup, l)
		}
	}

	return lineup
}

func HandleSubstituteLineupPlayers(p []*app.PlayerStats) []*proto.LineupPlayer {
	var lineup []*proto.LineupPlayer

	for _, player := range p {
		if player.IsSubstitute {
			l := converter.PlayerStatsToLineupPlayerProto(player)
			lineup = append(lineup, l)
		}
	}

	return lineup
}
