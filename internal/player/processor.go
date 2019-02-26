package player

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statistico-data/internal/model"
	"github.com/joesweeny/statistico-data/internal/squad"
	"log"
)

const callLimit = 1800
const player = "player"

var counter int

type Processor struct {
	Repository
	SquadRepo squad.Repository
	Factory
	Client *sportmonks.Client
	Logger *log.Logger
}

func (p Processor) Process(command string, done chan bool) {
	switch command {
	case player:
		go p.allPlayers(done)
	default:
		p.Logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (p Processor) allPlayers(done chan bool) {
	squads, err := p.SquadRepo.All()

	if err != nil {
		p.Logger.Fatalf("Error when retrieving squad data. Message: %s\n", err.Error())
		return
	}

	players := make(chan sportmonks.Player, callLimit)

	go p.parseSquads(squads, players, done, &counter)
	go p.parsePlayers(players, done)
}

func (p Processor) parseSquads(s []model.Squad, ch chan<- sportmonks.Player, done chan bool, c *int) {
	for _, sq := range s {
		for _, id := range sq.PlayerIDs {
			if _, err := p.Repository.Id(id); err != ErrNotFound {
				continue
			}

			if *c >= callLimit {
				p.Logger.Printf("Api call limited reached %d calls\n", *c)
				done <- true
			}

			res, err := p.Client.PlayerById(id, []string{}, 5)

			*c++

			if err != nil {
				p.Logger.Fatalf("Error when calling client '%s", err.Error())
				return
			}

			ch <- res.Data
		}
	}

	close(ch)
}

func (p Processor) parsePlayers(ch <-chan sportmonks.Player, done chan bool) {
	for x := range ch {
		p.persist(&x)
	}

	done <- true
}

func (p Processor) persist(pl *sportmonks.Player) {
	created := p.Factory.createPlayer(pl)

	if err := p.Insert(created); err != nil {
		log.Printf("Error '%s' occurred when inserting Player struct: %+v\n,", err.Error(), created)
	}

	return
}
