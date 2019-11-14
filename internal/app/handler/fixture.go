package handler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/converter"
	"github.com/statistico/statistico-data/internal/app/proto"
	"time"
)

type FixtureHandler struct {
	RoundRepo       app.RoundRepository
	TeamRepo        app.TeamRepository
	VenueRepo       app.VenueRepository
	Logger *logrus.Logger
}

func (h FixtureHandler) HandleFixture(f *app.Fixture) (*proto.Fixture, error) {
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
		HomeTeam:    converter.TeamToProto(home),
		AwayTeam:    converter.TeamToProto(away),
		DateTime:    &proto.Date{
			Utc:    f.Date.Unix(),
			Rfc:    f.Date.Format(time.RFC3339),
		},
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

	return &p, nil
}
