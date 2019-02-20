package stats

import (
	"github.com/joesweeny/sportmonks-go-client"
	"log"
)

type PlayerProcessor struct {
	PlayerRepository
	PlayerFactory
	Logger *log.Logger
}

func (p PlayerProcessor) ProcessPlayerStats(s *sportmonks.LineupPlayer, isSub bool) {
	x, err := p.PlayerRepository.ByFixtureAndPlayer(s.FixtureID, s.PlayerID)

	if err == ErrNotFound {
		created := p.PlayerFactory.createPlayerStats(s, isSub)

		if err := p.PlayerRepository.InsertPlayerStats(created); err != nil {
			log.Printf("Error '%s' occurred when inserting Player Stats struct: %+v\n,", err.Error(), created)
		}

		return
	}

	updated := p.PlayerFactory.updatePlayerStats(s, x)

	if err := p.PlayerRepository.UpdatePlayerStats(updated); err != nil {
		log.Printf("Error '%s' occurred when Updating Player Stats struct: %+v\n,", err.Error(), updated)
	}

	return
}
