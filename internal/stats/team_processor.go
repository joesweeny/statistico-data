package stats

import (
	"github.com/statistico/sportmonks-go-client"
	"log"
)

type TeamProcessor struct {
	TeamRepository
	TeamFactory
	Logger *log.Logger
}

func (p TeamProcessor) ProcessTeamStats(s *sportmonks.TeamStats) {
	x, err := p.TeamRepository.ByFixtureAndTeam(s.FixtureID, s.TeamID)

	if err == ErrNotFound {
		created := p.TeamFactory.createTeamStats(s)

		if err := p.TeamRepository.InsertTeamStats(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Team Stats struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := p.TeamFactory.updateTeamStats(s, x)

	if err := p.TeamRepository.UpdateTeamStats(updated); err != nil {
		log.Printf("Error '%s' occurred when updating Team Stats struct: %+v\n,", err.Error(), updated)
	}

	return
}
