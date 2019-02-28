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
	x.IsCup = c.IsCup

	return &x
}

func seasonToProto(s *model.Season) *pb.Season {
	var x pb.Season
	x.Id = int64(s.ID)
	x.Name = s.Name
	x.IsCurrent = s.IsCurrent

	return &x
}

func venueToProto(v *model.Venue) *pb.Venue {
	id := wrappers.Int64Value{}
	id.Value = int64(v.ID)
	name := wrappers.StringValue{}
	name.Value = v.Name

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
		HomeScore: int32(*res.HomeScore),
		AwayScore: int32(*res.AwayScore),
	}

	if res.HomePenScore != nil {
		a := wrappers.Int32Value{}
		a.Value = int32(*res.HomePenScore)
		stats.HomePenScore = &a
	}

	if res.AwayPenScore != nil {
		a := wrappers.Int32Value{}
		a.Value = int32(*res.AwayPenScore)
		stats.AwayPenScore = &a
	}

	if res.PitchCondition != nil {
		pitch := wrappers.StringValue{}
		pitch.Value = *res.PitchCondition
		stats.Pitch = &pitch
	}

	if res.HomeFormation != nil {
		h := wrappers.StringValue{}
		h.Value = *res.HomeFormation
		stats.HomeFormation = &h
	}

	if res.AwayFormation != nil {
		a := wrappers.StringValue{}
		a.Value = *res.AwayFormation
		stats.AwayFormation = &a
	}

	if res.HalfTimeScore != nil {
		a := wrappers.StringValue{}
		a.Value = *res.HalfTimeScore
		stats.HalfTimeScore = &a
	}

	if res.FullTimeScore != nil {
		a := wrappers.StringValue{}
		a.Value = *res.FullTimeScore
		stats.FullTimeScore = &a
	}

	if res.ExtraTimeScore != nil {
		a := wrappers.StringValue{}
		a.Value = *res.ExtraTimeScore
		stats.ExtraTimeScore = &a
	}

	if res.HomeLeaguePosition != nil {
		a := wrappers.Int32Value{}
		a.Value = int32(*res.HomeLeaguePosition)
		stats.HomeLeaguePosition = &a
	}

	if res.AwayLeaguePosition != nil {
		a := wrappers.Int32Value{}
		a.Value = int32(*res.AwayLeaguePosition)
		stats.AwayLeaguePosition = &a
	}

	if res.Minutes != nil {
		a := wrappers.Int32Value{}
		a.Value = int32(*res.Minutes)
		stats.Minutes = &a
	}

	if res.Seconds != nil {
		a := wrappers.Int32Value{}
		a.Value = int32(*res.Seconds)
		stats.Seconds = &a
	}

	if res.AddedTime != nil {
		a := wrappers.Int32Value{}
		a.Value = int32(*res.AddedTime)
		stats.AddedTime = &a
	}

	if res.ExtraTime != nil {
		a := wrappers.Int32Value{}
		a.Value = int32(*res.ExtraTime)
		stats.ExtraTime = &a
	}

	if res.InjuryTime != nil {
		a := wrappers.Int32Value{}
		a.Value = int32(*res.InjuryTime)
		stats.InjuryTime = &a
	}

	return &stats
}
