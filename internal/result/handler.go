package result

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data/internal/competition"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/statistico/statistico-data/internal/round"
	"github.com/statistico/statistico-data/internal/proto"
	"github.com/statistico/statistico-data/internal/season"
	"github.com/statistico/statistico-data/internal/team"
	"github.com/statistico/statistico-data/internal/venue"
	pbResult "github.com/statistico/statistico-data/internal/proto/result"
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

func (h Handler) HandleResult(f *model.Fixture, r *model.Result) (*pbResult.Result, error) {
	s, err := h.SeasonRepo.Id(f.SeasonID)

	if err != nil {
		e := fmt.Sprintf("Error when retrieving Result: FixtureID %d, Season ID %d", r.FixtureID, f.SeasonID)
		err = errors.New(e)
		h.logError(err)
		return nil, err
	}

	c, err := h.CompetitionRepo.GetById(s.LeagueID)

	if err != nil {
		e := fmt.Sprintf("Error when retrieving Result: FixtureID %d, Competition ID %d", r.FixtureID, s.LeagueID)
		err = errors.New(e)
		h.logError(err)
		return nil, err
	}

	home, err := h.TeamRepo.GetById(f.HomeTeamID)

	if err != nil {
		e := fmt.Sprintf("Error when retrieving Result: FixtureID %d, Home Team ID %d", r.FixtureID, f.HomeTeamID)
		err = errors.New(e)
		h.logError(err)
		return nil, err
	}

	away, err := h.TeamRepo.GetById(f.AwayTeamID)

	if err != nil {
		e := fmt.Sprintf("Error when retrieving Result: FixtureID %d, Away Team ID %d", r.FixtureID, f.AwayTeamID)
		err = errors.New(e)
		h.logError(err)
		return nil, err
	}

	p := pbResult.Result{
		Id:          int64(r.FixtureID),
		Competition: proto.CompetitionToProto(c),
		Season:      proto.SeasonToProto(s),
		DateTime:    f.Date.Unix(),
		MatchData:   proto.ToMatchData(home, away, r),
	}

	if f.RoundID != nil {
		rd, err := h.RoundRepo.GetById(*f.RoundID)

		if err != nil {
			e := fmt.Sprintf("Error when retrieving Result: FixtureID %d, Round ID %d", r.FixtureID, f.RoundID)
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
			e := fmt.Sprintf("Error when retrieving Result: FixtureID %d, Venue ID %d", r.FixtureID, f.VenueID)
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
