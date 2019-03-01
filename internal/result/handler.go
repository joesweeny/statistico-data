package result

import (
	"github.com/joesweeny/statistico-data/internal/competition"
	"github.com/joesweeny/statistico-data/internal/season"
	"github.com/joesweeny/statistico-data/internal/team"
	"github.com/joesweeny/statistico-data/internal/venue"
	"github.com/joesweeny/statistico-data/internal/model"
	pb "github.com/joesweeny/statistico-data/proto/result"
	"github.com/golang/protobuf/ptypes/wrappers"
)

type Handler struct {
	CompetitionRepo competition.Repository
	SeasonRepo season.Repository
	TeamRepo team.Repository
	VenueRepo venue.Repository
}

func (h Handler) HandleResult(f *model.Fixture, r *model.Result) (*pb.Result, error) {
	s, err := h.SeasonRepo.Id(f.SeasonID)

	if err != nil {
		return nil, err
	}

	c, err := h.CompetitionRepo.GetById(s.LeagueID)

	if err != nil {
		return nil, err
	}

	home, err := h.TeamRepo.GetById(f.HomeTeamID)

	if err != nil {
		return nil, err
	}

	away, err := h.TeamRepo.GetById(f.AwayTeamID)

	if err != nil {
		return nil, err
	}

	proto := pb.Result{
		Id: int64(r.FixtureID),
		Competition: competitionToProto(c),
		Season: seasonToProto(s),
		DateTime: f.Date.Unix(),
		MatchData: toMatchData(home, away, r),
	}

	if f.VenueID != nil {
		v, err := h.VenueRepo.GetById(*f.VenueID)

		if err != nil {
			return nil, err
		}

		proto.Venue = venueToProto(v)
	}

	if f.RefereeID != nil {
		ref := wrappers.Int64Value{}
		ref.Value = int64(*f.RefereeID)
		proto.RefereeId = &ref
	}

	return &proto, nil
}

func teamToProto(t *model.Team) *pb.Team {
	var x pb.Team
	x.Id = int64(t.ID)
	x.Name = t.Name

	return &x
}

func competitionToProto(c *model.Competition) *pb.Competition {
	var x pb.Competition
	x.Id = int64(c.ID)
	x.Name = c.Name
	x.IsCup = &wrappers.BoolValue{
		Value: c.IsCup,
	}

	return &x
}

func seasonToProto(s *model.Season) *pb.Season {
	var x pb.Season
	x.Id = int64(s.ID)
	x.Name = s.Name
	x.IsCurrent = &wrappers.BoolValue{
		Value: s.IsCurrent,
	}

	return &x
}

func venueToProto(v *model.Venue) *pb.Venue {
	id := wrappers.Int64Value{
		Value: int64(v.ID),
	}
	name := wrappers.StringValue{
		Value: v.Name,
	}
	ven := pb.Venue{}
	ven.Id = &id
	ven.Name = &name

	return &ven
}

func toMatchData(home *model.Team, away *model.Team, res *model.Result) *pb.MatchData {
	return &pb.MatchData{
		HomeTeam: teamToProto(home),
		AwayTeam: teamToProto(away),
		Stats: toMatchStats(res),
	}
}

func toMatchStats(res *model.Result) *pb.MatchStats {
	stats := pb.MatchStats{
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
