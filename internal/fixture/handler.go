package fixture

import(
	"github.com/joesweeny/statshub/internal/model"
	pb "github.com/joesweeny/statshub/proto/fixture"
	"github.com/joesweeny/statshub/internal/team"
	"github.com/joesweeny/statshub/internal/competition"
	"github.com/joesweeny/statshub/internal/season"
	"github.com/golang/protobuf/ptypes/wrappers"
)

type Handler struct {
	TeamRepo team.Repository
	CompetitionRepo competition.Repository
	SeasonRepo season.Repository
}

func (h Handler) handleFixture(f *model.Fixture) (*pb.Fixture, error) {
	s, err := h.SeasonRepo.Id(f.SeasonID)

	if err != nil {
		return nil, nil
	}

	c, err := h.CompetitionRepo.GetById(s.LeagueID)

	if err != nil {
		return nil, nil
	}

	home, err := h.TeamRepo.GetById(f.HomeTeamID)

	if err != nil {
		return nil, nil
	}

	away, err := h.TeamRepo.GetById(f.AwayTeamID)

	if err != nil {
		return nil, nil
	}

	return &pb.Fixture{
		Id: int64(f.ID),
		Competition: competitionToProto(c),
		Season: seasonToProto(s),
		HomeTeam: teamToProto(home),
		AwayTeam: teamToProto(away),
		VenueId: idToWrapperValue(f.VenueID),
		RefereeId: idToWrapperValue(f.RefereeID),
		DateTime: f.Date.Unix(),
	}, nil
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

func idToWrapperValue(id *int) *wrappers.Int32Value {
	var v wrappers.Int32Value
	v.Value = int32(*id)

	return &v
}
