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
	"errors"
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
		e := fmt.Sprintf("Error when retrieving Fixture: ID %d, Season ID %d", f.ID, f.SeasonID)
		err = errors.New(e)
		h.logError(err)
		return nil, err
	}

	c, err := h.CompetitionRepo.GetById(s.LeagueID)

	if err != nil {
		e := fmt.Sprintf("Error when retrieving Fixture: ID %d, Competition ID %d", f.ID, s.LeagueID)
		err = errors.New(e)
		h.logError(err)
		return nil, err
	}

	home, err := h.TeamRepo.GetById(f.HomeTeamID)

	if err != nil {
		e := fmt.Sprintf("Error when retrieving Fixture: ID %d, Home Team ID %d", f.ID, f.HomeTeamID)
		err = errors.New(e)
		h.logError(err)
		return nil, err
	}

	away, err := h.TeamRepo.GetById(f.AwayTeamID)

	if err != nil {
		e := fmt.Sprintf("Error when retrieving Fixture: ID %d, Away Team ID %d", f.ID, f.AwayTeamID)
		err = errors.New(e)
		h.logError(err)
		return nil, err
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
			e := fmt.Sprintf("Error when retrieving Fixture: ID %d, Round ID %d", f.ID, f.RoundID)
			err = errors.New(e)
			h.logError(err)
			p.Round = nil
		} else {
			p.Round = proto.RoundToProto(rd)
		}
	}

	if f.VenueID != nil {
		v, err := h.VenueRepo.GetById(*f.VenueID)

		if err != nil {
			e := fmt.Sprintf("Error when retrieving Fixture: ID %d, Venue ID %d", f.ID, f.VenueID)
			err = errors.New(e)
			h.logError(err)
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

func (h Handler) logError(e error) {
	h.Logger.Print(e.Error())
}