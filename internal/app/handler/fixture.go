package handler

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/converter"
	"github.com/statistico/statistico-data/internal/app/proto"
)

type FixtureHandler struct {
	CompetitionRepo app.CompetitionRepository
	RoundRepo       app.RoundRepository
	SeasonRepo      app.SeasonRepository
	TeamRepo        app.TeamRepository
	VenueRepo       app.VenueRepository
	Logger *logrus.Logger
}

func (h FixtureHandler) HandleFixture(f *app.Fixture) (*proto.Fixture, error) {
	s, err := h.SeasonRepo.ByID(uint64(f.SeasonID))

	if err != nil {
		e := fmt.Errorf("error when retrieving Fixture: ID %d, Season ID %d", f.ID, f.SeasonID)
		h.Logger.Println(e)
		return nil, e
	}

	c, err := h.CompetitionRepo.ByID(s.CompetitionID)

	if err != nil {
		e := fmt.Errorf("error when retrieving Fixture: ID %d, Competition ID %d", f.ID, s.CompetitionID)
		h.Logger.Println(e.Error())
		return nil, e
	}

	home, err := h.TeamRepo.ByID(uint64(f.HomeTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving Fixture: ID %d, Home Team ID %d", f.ID, f.HomeTeamID)
		h.Logger.Println(e.Error())
		return nil, e
	}

	away, err := h.TeamRepo.ByID(uint64(f.AwayTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving Fixture: ID %d, Away Team ID %d", f.ID, f.AwayTeamID)
		h.Logger.Println(e.Error())
		return nil, e
	}

	p := proto.Fixture{
		Id:          int64(f.ID),
		Competition: converter.CompetitionToProto(c),
		Season:      converter.SeasonToProto(s),
		HomeTeam:    converter.TeamToProto(home),
		AwayTeam:    converter.TeamToProto(away),
		DateTime:    f.Date.Unix(),
	}

	if f.RoundID != nil {
		rd, err := h.RoundRepo.ByID(uint64(*f.RoundID))

		if err != nil {
			e := fmt.Errorf("error when retrieving Fixture: ID %d, Round ID %d", f.ID, f.RoundID)
			h.Logger.Println(e.Error())
			p.Round = nil
		} else {
			p.Round = converter.RoundToProto(rd)
		}
	}

	if f.VenueID != nil {
		v, err := h.VenueRepo.GetById(uint64(*f.VenueID))

		if err != nil {
			e := fmt.Errorf("error when retrieving Fixture: ID %d, Venue ID %d", f.ID, f.VenueID)
			h.Logger.Println(e.Error())
			p.Venue = nil
		} else {
			p.Venue = converter.VenueToProto(v)
		}
	}

	if f.RefereeID != nil {
		ref := wrappers.Int64Value{}
		ref.Value = int64(*f.RefereeID)
		p.RefereeId = &ref
	}

	return &p, nil
}
