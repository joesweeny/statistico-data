package stats

import (
	"log"
	"github.com/joesweeny/sportmonks-go-client"
)

type PlayerProcessor struct {
	PlayerRepository
	PlayerFactory
	Logger *log.Logger
}

func (p PlayerProcessor) ProcessPlayerStats(s *sportmonks.LineupPlayer, isSub bool) {
	_, err := p.PlayerRepository.ByFixtureAndPlayer(s.FixtureID, s.PlayerID)

	if err == ErrNotFound {
		created := p.PlayerFactory.createPlayerStats(s, isSub)

		if err := p.PlayerRepository.InsertPlayerStats(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Player Stats struct: %+v\n,", err.Error(), created)
		}

		return
	}
}
