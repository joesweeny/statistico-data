package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"sync"
)

const callLimit = 1800
const player = "player"

var counter int

type PlayerProcessor struct {
	playerRepo app.PlayerRepository
	squadRepo  app.SquadRepository
	requester  app.PlayerRequester
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

func (p PlayerProcessor) parseSquads(s []app.Squad, ch chan<- *app.Player, done chan bool, c *int) {
	var wg sync.WaitGroup

	for _, sq := range s {
		wg.Add(1)

		if *c >= callLimit {
			p.logger.Printf("Api call limited reached %d calls\n", *c)

			close(ch)
		}

		for _, id := range sq.PlayerIDs {
			if _, err := p.playerRepo.ByID(id); err != nil {
				p.logger.Warnf("Failure when fetching squad player data: %s", err.Error())

				pl, err := p.requester.PlayerByID(id)

				if err != nil {
					p.logger.Warnf("Failure when fetching sportmonks player data: %s", err.Error())
				}

				if err == nil {
					ch <- pl
				}

				*c++
			}
		}

		wg.Done()
	}

	wg.Wait()

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
		p.logger.Warnf("Error '%s' occurred inserting player struct when processing: %+v\n,", err.Error(), *x)
	}
}

func NewPlayerProcessor(r app.PlayerRepository, s app.SquadRepository, q app.PlayerRequester, log *logrus.Logger) *PlayerProcessor {
	return &PlayerProcessor{playerRepo: r, squadRepo: s, requester: q, logger: log}
}
