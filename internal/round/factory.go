package round

import (
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/model"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02"

type Factory struct {
	Clock clockwork.Clock
}

func (f Factory) createRound(s *sportmonks.Round) (*model.Round, error) {
	start, err := time.Parse(dateFormat, s.Start)

	if err != nil {
		return &model.Round{}, err
	}

	end, err := time.Parse(dateFormat, s.End)

	if err != nil {
		return &model.Round{}, err
	}

	round := model.Round{
		ID:        s.ID,
		Name:      strconv.Itoa(s.Name),
		SeasonID:  s.SeasonID,
		StartDate: start,
		EndDate:   end,
		CreatedAt: f.Clock.Now(),
		UpdatedAt: f.Clock.Now(),
	}

	return &round, nil
}

func (f Factory) updateRound(s *sportmonks.Round, m *model.Round) (*model.Round, error) {
	start, err := time.Parse(dateFormat, s.Start)

	if err != nil {
		return &model.Round{}, err
	}

	end, err := time.Parse(dateFormat, s.End)

	if err != nil {
		return &model.Round{}, err
	}

	m.Name = strconv.Itoa(s.Name)
	m.StartDate = start
	m.EndDate = end
	m.UpdatedAt = f.Clock.Now()

	return m, nil
}
