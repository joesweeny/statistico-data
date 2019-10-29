package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/model"
	"github.com/statistico/statistico-data/internal/squad"
	"log"
)

const callLimit = 1800
const player = "player"

var counter int

type PlayerProcessor struct {
	playerRepo app.PlayerRepository
	squadRepo squad.Repository
	requester app.PlayerRequester
	logger     *logrus.Logger
}

func (p PlayerProcessor) Process(command string, option string, done chan bool) {
	if command != player {
		p.logger.Fatalf("Command %s is not supported", command)
	}

	squads, err := p.squadRepo.All()

	if err != nil {
		p.logger.Fatalf("Error retrieving squad data. Message: %s\n", err.Error())
		return
	}

	ch := make(chan *app.Player, callLimit)

	go p.parseSquads(squads, ch, done, &counter)
	go p.parsePlayers(ch, done)
}

func (p PlayerProcessor) parseSquads(s []model.Squad, ch chan<- *app.Player, done chan bool, c *int) {
	for _, sq := range s {
		for _, id := range sq.PlayerIDs {
			if _, err := p.playerRepo.ByID(int64(id)); err != nil {
				continue
			}

			if *c >= callLimit {
				p.logger.Printf("Api call limited reached %d calls\n", *c)
				done <- true
			}

			res := p.requester.PlayerByID(int64(id))

			*c++

			ch <- res
		}
	}

	close(ch)
}

func (p PlayerProcessor) parsePlayers(ch <-chan *app.Player, done chan bool) {
	for x := range ch {
		p.persist(x)
	}

	done <- true
}

func (p PlayerProcessor) persist(x *app.Player) {
	if err := p.playerRepo.Insert(x); err != nil {
		log.Printf("Error '%s' occurred when inserting Player struct: %+v\n,", err.Error(), &x)
	}
}
