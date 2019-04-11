package fixture

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data/internal/competition"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/statistico/statistico-data/internal/round"
	"github.com/statistico/statistico-data/internal/proto"
	"github.com/statistico/statistico-data/internal/season"
	"github.com/statistico/statistico-data/internal/team"
	"github.com/statistico/statistico-data/internal/venue"
	pbFixture "github.com/statistico/statistico-data/internal/proto/fixture"
	"log"
	"fmt"
)

type Handler struct {
	CompetitionRepo competition.Repository
	RoundRepo       round.Repository
	SeasonRepo      season.Repository
	TeamRepo        team.Repository
	VenueRepo       venue.Repository
	Logger 			*log.Logger
}

func (h Handler) HandleFixture(f *model.Fixture) (*pbFixture.Fixture, error) {
	s, err := h.SeasonRepo.Id(f.SeasonID)

	if err != nil {
		e := fmt.Errorf("error when retrieving Fixture: ID %d, Season ID %d", f.ID, f.SeasonID)
		h.Logger.Println(e)
		return nil, e
	}

	c, err := h.CompetitionRepo.GetById(s.LeagueID)

	if err != nil {
		e := fmt.Errorf("error when retrieving Fixture: ID %d, Competition ID %d", f.ID, s.LeagueID)
		h.Logger.Println(e.Error())
		return nil, e
	}

	home, err := h.TeamRepo.GetById(f.HomeTeamID)

	if err != nil {
		e := fmt.Errorf("error when retrieving Fixture: ID %d, Home Team ID %d", f.ID, f.HomeTeamID)
		h.Logger.Println(e.Error())
		return nil, e
	}

	away, err := h.TeamRepo.GetById(f.AwayTeamID)

	if err != nil {
		e := fmt.Errorf("error when retrieving Fixture: ID %d, Away Team ID %d", f.ID, f.AwayTeamID)
		h.Logger.Println(e.Error())
		return nil, e
	}

	p := pbFixture.Fixture{
		Id:          int64(f.ID),
		Competition: proto.CompetitionToProto(c),
		Season:      proto.SeasonToProto(s),
		HomeTeam:    proto.TeamToProto(home),
		AwayTeam:    proto.TeamToProto(away),
		DateTime:    f.Date.Unix(),
	}

	if f.RoundID != nil {
		rd, err := h.RoundRepo.GetById(*f.RoundID)

		if err != nil {
			e := fmt.Errorf("error when retrieving Fixture: ID %d, Round ID %d", f.ID, f.RoundID)
			h.Logger.Println(e.Error())
			p.Round = nil
		} else {
			p.Round = proto.RoundToProto(rd)
		}
	}

	if f.VenueID != nil {
		v, err := h.VenueRepo.GetById(*f.VenueID)

		if err != nil {
			e := fmt.Errorf("error when retrieving Fixture: ID %d, Venue ID %d", f.ID, f.VenueID)
			h.Logger.Println(e.Error())
			p.Venue = nil
		} else {
			p.Venue = proto.VenueToProto(v)
		}
	}

	if f.RefereeID != nil {
		ref := wrappers.Int64Value{}
		ref.Value = int64(*f.RefereeID)
		p.RefereeId = &ref
	}

	return &p, nil
}
