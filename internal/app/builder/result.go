package builder

import (
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/app/converter"
	"github.com/statistico/statistico-data/internal/app/fetch"
	"github.com/statistico/statistico-data/internal/app/proto"
)

type ResultBuilder struct {
	fetcher fetch.EntityFetcher
}

func (b ResultBuilder) BuildProtoResult(f *app.Fixture) (*proto.Result, error) {
	r, err := b.fetcher.ResultByID(f.ID)

	if err != nil {
		return nil, err
	}

	s, err := b.fetcher.SeasonByID(f.SeasonID)

	if err != nil {
		return nil, err
	}

	c, err := b.fetcher.CompetitionByID(s.CompetitionID)

	if err != nil {
		return nil, err
	}

	home, err := b.fetcher.TeamByID(f.HomeTeamID)

	if err != nil {
		return nil, err
	}

	away, err := b.fetcher.TeamByID(f.AwayTeamID)

	if err != nil {
		return nil, err
	}

	p := proto.Result{
		Id:          int64(f.ID),
		Competition: converter.CompetitionToProto(c),
		Season:      converter.SeasonToProto(s),
		DateTime:    f.Date.Unix(),
		MatchData:   converter.ToMatchData(home, away, r),
	}
	if f.VenueID != nil {
		v, err := b.fetcher.VenueByID(*f.VenueID)

		if err == nil {
			p.Venue = converter.VenueToProto(v)
		}
	}

	if f.RoundID != nil {
		r, err := b.fetcher.RoundByID(*f.RoundID)

		if err == nil {
			p.Round = converter.RoundToProto(r)
		}
	}

	return &p, nil
}
