package result

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/statistico/statistico-data/internal/proto"
	pbResult "github.com/statistico/statistico-data/internal/proto/result"
	"log"
)

type Handler struct {
	CompetitionRepo app.CompetitionRepository
	RoundRepo       app.RoundRepository
	SeasonRepo      app.SeasonRepository
	TeamRepo        app.TeamRepository
	VenueRepo       app.VenueRepository
	Logger          *log.Logger
}

func (h Handler) HandleResult(f *model.Fixture, r *model.Result) (*pbResult.Result, error) {
	s, err := h.SeasonRepo.ByID(int64(f.SeasonID))

	if err != nil {
		e := fmt.Errorf("error when retrieving Result: FixtureID %d, Season ID %d", r.FixtureID, f.SeasonID)
		h.Logger.Println(e)
		return nil, e
	}

	c, err := h.CompetitionRepo.ByID(s.CompetitionID)

	if err != nil {
		e := fmt.Errorf("error when retrieving Result: FixtureID %d, Competition ID %d", r.FixtureID, s.CompetitionID)
		h.Logger.Println(e)
		return nil, e
	}

	home, err := h.TeamRepo.ByID(int64(f.HomeTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving Result: FixtureID %d, Home Team ID %d", r.FixtureID, f.HomeTeamID)
		h.Logger.Println(e)
		return nil, e
	}

	away, err := h.TeamRepo.ByID(int64(f.AwayTeamID))

	if err != nil {
		e := fmt.Errorf("error when retrieving Result: FixtureID %d, Away Team ID %d", r.FixtureID, f.AwayTeamID)
		h.Logger.Println(e)
		return nil, e
	}

	p := pbResult.Result{
		Id:          int64(r.FixtureID),
		Competition: proto.CompetitionToProto(c),
		Season:      proto.SeasonToProto(s),
		DateTime:    f.Date.Unix(),
		MatchData:   proto.ToMatchData(home, away, r),
	}

	if f.RoundID != nil {
		rd, err := h.RoundRepo.ByID(int64(*f.RoundID))

		if err != nil {
			e := fmt.Errorf("error when retrieving Result: FixtureID %d, Round ID %d", r.FixtureID, f.RoundID)
			h.Logger.Println(e)
			p.Round = nil
		} else {
			p.Round = proto.RoundToProto(rd)
		}
	}

	if f.VenueID != nil {
		v, err := h.VenueRepo.GetById(int64(*f.VenueID))

		if err != nil {
			e := fmt.Errorf("error when retrieving Result: FixtureID %d, Venue ID %d", r.FixtureID, f.VenueID)
			h.Logger.Println(e)
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
