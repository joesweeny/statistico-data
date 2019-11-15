package factory

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"time"
)

// Convert a domain Team struct into a proto Team struct
func TeamToProto(t *app.Team) *proto.Team {
	var x proto.Team
	x.Id = int64(t.ID)
	x.Name = t.Name

	return &x
}

// Convert a domain Competition struct into a proto Competition struct
func CompetitionToProto(c *app.Competition) *proto.Competition {
	var x proto.Competition
	x.Id = int64(c.ID)
	x.Name = c.Name
	x.IsCup = &wrappers.BoolValue{
		Value: c.IsCup,
	}

	return &x
}

// Convert a domain PlayerStats struct into a proto LineupPlayer struct
func PlayerStatsToLineupPlayerProto(p *app.PlayerStats) *proto.LineupPlayer {
	player := proto.LineupPlayer{
		PlayerId:     p.PlayerID,
		Position:     *p.Position,
		IsSubstitute: p.IsSubstitute,
	}

	if p.FormationPosition != nil {
		player.FormationPosition = &wrappers.UInt32Value{
			Value: uint32(*p.FormationPosition),
		}
	}

	return &player
}

// Convert a domain PlayerStats struct into a proto PlayerStats struct
func PlayerStatsToProto(p *app.PlayerStats) *proto.PlayerStats {
	stats := proto.PlayerStats{
		PlayerId: p.PlayerID,
	}

	shots := p.PlayerShots

	if shots.Total != nil {
		stats.ShotsTotal = &wrappers.Int32Value{
			Value: int32(*shots.Total),
		}
	}

	if shots.OnGoal != nil {
		stats.ShotsOnGoal = &wrappers.Int32Value{
			Value: int32(*shots.OnGoal),
		}
	}

	goals := p.PlayerGoals

	if goals.Scored != nil {
		stats.GoalsScored = &wrappers.Int32Value{
			Value: int32(*goals.Scored),
		}
	}

	if goals.Conceded != nil {
		stats.GoalsConceded = &wrappers.Int32Value{
			Value: int32(*goals.Conceded),
		}
	}

	if p.Assists != nil {
		stats.Assists = &wrappers.Int32Value{
			Value: int32(*p.Assists),
		}
	}

	return &stats
}

// Convert a domain Round struct into a proto Round struct
func RoundToProto(r *app.Round) *proto.Round {
	return &proto.Round{
		Id:        int64(r.ID),
		Name:      r.Name,
		SeasonId:  int64(r.SeasonID),
		StartDate: r.StartDate.Format(time.RFC3339),
		EndDate:   r.EndDate.Format(time.RFC3339),
	}
}

// Convert a domain Season struct into a proto Season struct
func SeasonToProto(s *app.Season) *proto.Season {
	var x proto.Season
	x.Id = int64(s.ID)
	x.Name = s.Name
	x.IsCurrent = &wrappers.BoolValue{
		Value: s.IsCurrent,
	}

	return &x
}

// Convert a domain TeamStats struct into a proto TeamStats struct
func TeamStatsToProto(t *app.TeamStats) *proto.TeamStats {
	stats := proto.TeamStats{
		TeamId: t.TeamID,
	}

	shots := t.TeamShots

	if shots.Total != nil {
		stats.ShotsTotal = &wrappers.UInt32Value{
			Value: uint32(*shots.Total),
		}
	}

	if shots.OnGoal != nil {
		stats.ShotsOnGoal = &wrappers.UInt32Value{
			Value: uint32(*shots.OnGoal),
		}
	}

	if shots.OffGoal != nil {
		stats.ShotsOffGoal = &wrappers.UInt32Value{
			Value: uint32(*shots.OffGoal),
		}
	}

	if shots.Blocked != nil {
		stats.ShotsBlocked = &wrappers.UInt32Value{
			Value: uint32(*shots.Blocked),
		}
	}

	if shots.InsideBox != nil {
		stats.ShotsInsideBox = &wrappers.UInt32Value{
			Value: uint32(*shots.InsideBox),
		}
	}

	if shots.OutsideBox != nil {
		stats.ShotsOutsideBox = &wrappers.UInt32Value{
			Value: uint32(*shots.OutsideBox),
		}
	}

	passes := t.TeamPasses

	if passes.Total != nil {
		stats.PassesTotal = &wrappers.UInt32Value{
			Value: uint32(*passes.Total),
		}
	}

	if passes.Accuracy != nil {
		stats.PassesAccuracy = &wrappers.UInt32Value{
			Value: uint32(*passes.Accuracy),
		}
	}

	if passes.Percentage != nil {
		stats.PassesPercentage = &wrappers.UInt32Value{
			Value: uint32(*passes.Percentage),
		}
	}

	attacks := t.TeamAttacks

	if attacks.Total != nil {
		stats.AttacksTotal = &wrappers.UInt32Value{
			Value: uint32(*attacks.Total),
		}
	}

	if attacks.Dangerous != nil {
		stats.AttacksDangerous = &wrappers.UInt32Value{
			Value: uint32(*attacks.Dangerous),
		}
	}

	if t.Fouls != nil {
		stats.Fouls = &wrappers.UInt32Value{
			Value: uint32(*t.Fouls),
		}
	}

	if t.Corners != nil {
		stats.Corners = &wrappers.UInt32Value{
			Value: uint32(*t.Corners),
		}
	}

	if t.Offsides != nil {
		stats.Offsides = &wrappers.UInt32Value{
			Value: uint32(*t.Offsides),
		}
	}

	if t.Possession != nil {
		stats.Possession = &wrappers.UInt32Value{
			Value: uint32(*t.Possession),
		}
	}

	if t.YellowCards != nil {
		stats.YellowCards = &wrappers.UInt32Value{
			Value: uint32(*t.YellowCards),
		}
	}

	if t.RedCards != nil {
		stats.RedCards = &wrappers.UInt32Value{
			Value: uint32(*t.RedCards),
		}
	}

	if t.Saves != nil {
		stats.Saves = &wrappers.UInt32Value{
			Value: uint32(*t.Saves),
		}
	}

	if t.Substitutions != nil {
		stats.Substitutions = &wrappers.UInt32Value{
			Value: uint32(*t.Substitutions),
		}
	}

	if t.GoalKicks != nil {
		stats.GoalKicks = &wrappers.UInt32Value{
			Value: uint32(*t.GoalKicks),
		}
	}

	if t.GoalAttempts != nil {
		stats.GoalAttempts = &wrappers.UInt32Value{
			Value: uint32(*t.GoalAttempts),
		}
	}

	if t.FreeKicks != nil {
		stats.FreeKicks = &wrappers.UInt32Value{
			Value: uint32(*t.FreeKicks),
		}
	}

	if t.ThrowIns != nil {
		stats.ThrowIns = &wrappers.UInt32Value{
			Value: uint32(*t.ThrowIns),
		}
	}

	return &stats
}

// Convert a domain Venue struct into a proto Venue struct
func VenueToProto(v *app.Venue) *proto.Venue {
	id := wrappers.Int64Value{
		Value: int64(v.ID),
	}
	name := wrappers.StringValue{
		Value: v.Name,
	}
	ven := proto.Venue{}
	ven.Id = &id
	ven.Name = &name

	return &ven
}

// Convert a domain Team and Result struct data into a proto MatchData struct
func ToMatchData(home *app.Team, away *app.Team, res *app.Result) *proto.MatchData {
	return &proto.MatchData{
		HomeTeam: TeamToProto(home),
		AwayTeam: TeamToProto(away),
		Stats:    ToMatchStats(res),
	}
}

// Convert a domain Result struct into a proto MatchStats struct
func ToMatchStats(res *app.Result) *proto.MatchStats {
	stats := proto.MatchStats{
		HomeScore: &wrappers.Int32Value{
			Value: int32(*res.HomeScore),
		},
		AwayScore: &wrappers.Int32Value{
			Value: int32(*res.AwayScore),
		},
	}

	if res.HomePenScore != nil {
		stats.HomePenScore = &wrappers.Int32Value{
			Value: int32(*res.HomePenScore),
		}
	}

	if res.AwayPenScore != nil {
		stats.AwayPenScore = &wrappers.Int32Value{
			Value: int32(*res.AwayPenScore),
		}
	}

	if res.PitchCondition != nil {
		stats.Pitch = &wrappers.StringValue{
			Value: *res.PitchCondition,
		}
	}

	if res.HomeFormation != nil {
		stats.HomeFormation = &wrappers.StringValue{
			Value: *res.HomeFormation,
		}
	}

	if res.AwayFormation != nil {
		stats.AwayFormation = &wrappers.StringValue{
			Value: *res.AwayFormation,
		}
	}

	if res.HalfTimeScore != nil {
		stats.HalfTimeScore = &wrappers.StringValue{
			Value: *res.HalfTimeScore,
		}
	}

	if res.FullTimeScore != nil {
		stats.FullTimeScore = &wrappers.StringValue{
			Value: *res.FullTimeScore,
		}
	}

	if res.ExtraTimeScore != nil {
		stats.ExtraTimeScore = &wrappers.StringValue{
			Value: *res.ExtraTimeScore,
		}
	}

	if res.HomeLeaguePosition != nil {
		stats.HomeLeaguePosition = &wrappers.Int32Value{
			Value: int32(*res.HomeLeaguePosition),
		}
	}

	if res.AwayLeaguePosition != nil {
		stats.AwayLeaguePosition = &wrappers.Int32Value{
			Value: int32(*res.AwayLeaguePosition),
		}
	}

	if res.Minutes != nil {
		stats.Minutes = &wrappers.Int32Value{
			Value: int32(*res.Minutes),
		}
	}

	if res.AddedTime != nil {
		stats.AddedTime = &wrappers.Int32Value{
			Value: int32(*res.AddedTime),
		}
	}

	if res.ExtraTime != nil {
		stats.ExtraTime = &wrappers.Int32Value{
			Value: int32(*res.ExtraTime),
		}
	}

	if res.InjuryTime != nil {
		stats.InjuryTime = &wrappers.Int32Value{
			Value: int32(*res.InjuryTime),
		}
	}

	return &stats
}

func HandlePlayerStats(p []*app.PlayerStats) []*proto.PlayerStats {
	var stats []*proto.PlayerStats

	for _, player := range p {
		s := PlayerStatsToProto(player)
		stats = append(stats, s)
	}

	return stats
}

func HandleStartingLineupPlayers(p []*app.PlayerStats) []*proto.LineupPlayer {
	var lineup []*proto.LineupPlayer

	for _, player := range p {
		if !player.IsSubstitute {
			l := PlayerStatsToLineupPlayerProto(player)
			lineup = append(lineup, l)
		}
	}

	return lineup
}

func HandleSubstituteLineupPlayers(p []*app.PlayerStats) []*proto.LineupPlayer {
	var lineup []*proto.LineupPlayer

	for _, player := range p {
		if player.IsSubstitute {
			l := PlayerStatsToLineupPlayerProto(player)
			lineup = append(lineup, l)
		}
	}

	return lineup
}