package fixture

import(
	"github.com/joesweeny/statshub/internal/model"
	pb "github.com/joesweeny/statshub/proto/fixture"
	"github.com/joesweeny/statshub/internal/team"
	"github.com/joesweeny/statshub/internal/competition"
	"github.com/joesweeny/statshub/internal/season"
	"github.com/joesweeny/statshub/internal/venue"
	"github.com/golang/protobuf/ptypes/wrappers"
)

type Handler struct {
	CompetitionRepo competition.Repository
	SeasonRepo season.Repository
	TeamRepo team.Repository
	VenueRepo venue.Repository
}

func (h Handler) HandleFixture(f *model.Fixture) (*pb.Fixture, error) {
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

	proto := pb.Fixture{
		Id: int64(f.ID),
		Competition: competitionToProto(c),
		Season: seasonToProto(s),
		HomeTeam: teamToProto(home),
		AwayTeam: teamToProto(away),
		DateTime: f.Date.Unix(),
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

	return &x
}

func seasonToProto(s *model.Season) *pb.Season {
	var x pb.Season
	x.Id = int64(s.ID)
	x.Name = s.Name

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
