package proto

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data/internal/app"
	pbCompetition "github.com/statistico/statistico-data/internal/proto/competition"
	pbResult "github.com/statistico/statistico-data/internal/proto/result"
	pbRound "github.com/statistico/statistico-data/internal/proto/round"
	pbSeason "github.com/statistico/statistico-data/internal/proto/season"
	pbPlayerStats "github.com/statistico/statistico-data/internal/proto/stats/player"
	pbTeamStats "github.com/statistico/statistico-data/internal/proto/stats/team"
	pbTeam "github.com/statistico/statistico-data/internal/proto/team"
	pbVenue "github.com/statistico/statistico-data/internal/proto/venue"
	"time"
)

func TeamToProto(t *app.Team) *pbTeam.Team {
	var x pbTeam.Team
	x.Id = int64(t.ID)
	x.Name = t.Name

	return &x
}

func CompetitionToProto(c *app.Competition) *pbCompetition.Competition {
	var x pbCompetition.Competition
	x.Id = int64(c.ID)
	x.Name = c.Name
	x.IsCup = &wrappers.BoolValue{
		Value: c.IsCup,
	}

	return &x
}

func PlayerStatsToLineupPlayerProto(p *app.PlayerStats) *pbPlayerStats.LineupPlayer {
	player := pbPlayerStats.LineupPlayer{
		PlayerId:     uint64(p.PlayerID),
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

func PlayerStatsToProto(p *app.PlayerStats) *pbPlayerStats.PlayerStats {
	stats := pbPlayerStats.PlayerStats{
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

func RoundToProto(r *app.Round) *pbRound.Round {
	return &pbRound.Round{
		Id:        int64(r.ID),
		Name:      r.Name,
		SeasonId:  int64(r.SeasonID),
		StartDate: r.StartDate.Format(time.RFC3339),
		EndDate:   r.EndDate.Format(time.RFC3339),
	}
}

func SeasonToProto(s *app.Season) *pbSeason.Season {
	var x pbSeason.Season
	x.Id = int64(s.ID)
	x.Name = s.Name
	x.IsCurrent = &wrappers.BoolValue{
		Value: s.IsCurrent,
	}

	return &x
}

func TeamStatsToProto(t *app.TeamStats) *pbTeamStats.TeamStats {
	stats := pbTeamStats.TeamStats{
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

func VenueToProto(v *app.Venue) *pbVenue.Venue {
	id := wrappers.Int64Value{
		Value: int64(v.ID),
	}
	name := wrappers.StringValue{
		Value: v.Name,
	}
	ven := pbVenue.Venue{}
	ven.Id = &id
	ven.Name = &name

	return &ven
}

func ToMatchData(home *app.Team, away *app.Team, res *app.Result) *pbResult.MatchData {
	return &pbResult.MatchData{
		HomeTeam: TeamToProto(home),
		AwayTeam: TeamToProto(away),
		Stats:    ToMatchStats(res),
	}
}

func ToMatchStats(res *app.Result) *pbResult.MatchStats {
	stats := pbResult.MatchStats{
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
