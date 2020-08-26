package factory

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/grpc/proto"
	"time"
)

// Convert a domain Team struct into a proto Team struct
func TeamToProto(t *app.Team) *proto.Team {
	team := proto.Team{
		Id:             t.ID,
		Name:           t.Name,
		CountryId:      t.CountryID,
		VenueId:        t.VenueID,
		IsNationalTeam: &wrappers.BoolValue{Value: t.NationalTeam},
	}

	if t.ShortCode != nil {
		team.ShortCode = &wrappers.StringValue{Value: *t.ShortCode}
	}

	if t.Founded != nil {
		team.Founded = &wrappers.UInt64Value{Value: uint64(*t.Founded)}
	}

	if t.Logo != nil {
		team.Logo = &wrappers.StringValue{Value: *t.Logo}
	}

	return &team
}

// Convert a domain Competition struct into a proto Competition struct
func CompetitionToProto(c *app.Competition) *proto.Competition {
	var x proto.Competition
	x.Id = c.ID
	x.Name = c.Name
	x.IsCup = c.IsCup
	x.CountryId = c.CountryID

	return &x
}

// Convert a domain PlayerStats struct into a proto LineupPlayer struct
func playerStatsToLineupPlayerProto(p *app.PlayerStats) *proto.LineupPlayer {
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
func playerStatsToProto(p *app.PlayerStats) *proto.PlayerStats {
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
func roundToProto(r *app.Round) *proto.Round {
	return &proto.Round{
		Id:        r.ID,
		Name:      r.Name,
		SeasonId:  r.SeasonID,
		StartDate: r.StartDate.Format(time.RFC3339),
		EndDate:   r.EndDate.Format(time.RFC3339),
	}
}

// Convert a domain Season struct into a proto Season struct
func SeasonToProto(s *app.Season) *proto.Season {
	var x proto.Season
	x.Id = s.ID
	x.Name = s.Name
	x.IsCurrent = &wrappers.BoolValue{
		Value: s.IsCurrent,
	}

	return &x
}

// Convert a domain TeamStats struct into a proto TeamStats struct
func teamStatsToProto(t *app.TeamStats) *proto.TeamStats {
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

	if t.Goals != nil {
		stats.Goals = &wrappers.UInt32Value{
			Value: uint32(*t.Goals),
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

func TeamStatToProto(s *app.TeamStat) *proto.TeamStat {
	return &proto.TeamStat{
		FixtureId: s.FixtureID,
		Stat:      s.Stat,
		Value:     &wrappers.UInt32Value{
			Value: *s.Value,
		},
	}
}

// Convert a domain Venue struct into a proto Venue struct
func venueToProto(v *app.Venue) *proto.Venue {
	return &proto.Venue{
		Id: v.ID,
		Name: v.Name,
	}
}

// Convert a domain Result struct into a proto MatchStats struct
func toMatchStats(res *app.Result) *proto.MatchStats {
	stats := proto.MatchStats{
		HomeScore: &wrappers.UInt32Value{
			Value: uint32(*res.HomeScore),
		},
		AwayScore: &wrappers.UInt32Value{
			Value: uint32(*res.AwayScore),
		},
	}

	if res.HomePenScore != nil {
		stats.HomePenScore = &wrappers.UInt32Value{
			Value: uint32(*res.HomePenScore),
		}
	}

	if res.AwayPenScore != nil {
		stats.AwayPenScore = &wrappers.UInt32Value{
			Value: uint32(*res.AwayPenScore),
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
		stats.HomeLeaguePosition = &wrappers.UInt32Value{
			Value: uint32(*res.HomeLeaguePosition),
		}
	}

	if res.AwayLeaguePosition != nil {
		stats.AwayLeaguePosition = &wrappers.UInt32Value{
			Value: uint32(*res.AwayLeaguePosition),
		}
	}

	if res.Minutes != nil {
		stats.Minutes = &wrappers.UInt32Value{
			Value: uint32(*res.Minutes),
		}
	}

	if res.AddedTime != nil {
		stats.AddedTime = &wrappers.UInt32Value{
			Value: uint32(*res.AddedTime),
		}
	}

	if res.ExtraTime != nil {
		stats.ExtraTime = &wrappers.UInt32Value{
			Value: uint32(*res.ExtraTime),
		}
	}

	if res.InjuryTime != nil {
		stats.InjuryTime = &wrappers.UInt32Value{
			Value: uint32(*res.InjuryTime),
		}
	}

	return &stats
}

func handlePlayerStats(p []*app.PlayerStats) []*proto.PlayerStats {
	var stats []*proto.PlayerStats

	for _, player := range p {
		s := playerStatsToProto(player)
		stats = append(stats, s)
	}

	return stats
}

func handleStartingLineupPlayers(p []*app.PlayerStats) []*proto.LineupPlayer {
	var lineup []*proto.LineupPlayer

	for _, player := range p {
		if !player.IsSubstitute {
			l := playerStatsToLineupPlayerProto(player)
			lineup = append(lineup, l)
		}
	}

	return lineup
}

func handleSubstituteLineupPlayers(p []*app.PlayerStats) []*proto.LineupPlayer {
	var lineup []*proto.LineupPlayer

	for _, player := range p {
		if player.IsSubstitute {
			l := playerStatsToLineupPlayerProto(player)
			lineup = append(lineup, l)
		}
	}

	return lineup
}
