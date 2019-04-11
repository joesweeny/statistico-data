package proto

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data/internal/model"
	pbCompetition "github.com/statistico/statistico-data/internal/proto/competition"
	pbResult "github.com/statistico/statistico-data/internal/proto/result"
	pbRound "github.com/statistico/statistico-data/internal/proto/round"
	pbSeason "github.com/statistico/statistico-data/internal/proto/season"
	pbTeam "github.com/statistico/statistico-data/internal/proto/team"
	pbVenue "github.com/statistico/statistico-data/internal/proto/venue"
	"time"
)

func TeamToProto(t *model.Team) *pbTeam.Team {
	var x pbTeam.Team
	x.Id = int64(t.ID)
	x.Name = t.Name

	return &x
}

func CompetitionToProto(c *model.Competition) *pbCompetition.Competition {
	var x pbCompetition.Competition
	x.Id = int64(c.ID)
	x.Name = c.Name
	x.IsCup = &wrappers.BoolValue{
		Value: c.IsCup,
	}

	return &x
}

func RoundToProto(r *model.Round) *pbRound.Round {
	return &pbRound.Round{
		Id:        int64(r.ID),
		Name:      r.Name,
		SeasonId:  int64(r.SeasonID),
		StartDate: r.StartDate.Format(time.RFC3339),
		EndDate:   r.EndDate.Format(time.RFC3339),
	}
}

func SeasonToProto(s *model.Season) *pbSeason.Season {
	var x pbSeason.Season
	x.Id = int64(s.ID)
	x.Name = s.Name
	x.IsCurrent = &wrappers.BoolValue{
		Value: s.IsCurrent,
	}

	return &x
}

func VenueToProto(v *model.Venue) *pbVenue.Venue {
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

func ToMatchData(home *model.Team, away *model.Team, res *model.Result) *pbResult.MatchData {
	return &pbResult.MatchData{
		HomeTeam: TeamToProto(home),
		AwayTeam: TeamToProto(away),
		Stats:    ToMatchStats(res),
	}
}

func ToMatchStats(res *model.Result) *pbResult.MatchStats {
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

	if res.Seconds != nil {
		stats.Seconds = &wrappers.Int32Value{
			Value: int32(*res.Seconds),
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
