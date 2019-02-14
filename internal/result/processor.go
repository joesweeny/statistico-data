package result

import (
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/event"
	"github.com/joesweeny/statshub/internal/fixture"
	"github.com/joesweeny/statshub/internal/stats"
	"log"
	"sync"
)

const result = "result"
const callLimit = 1500

var counter int
var waitGroup sync.WaitGroup

type Processor struct {
	Repository
	FixtureRepo fixture.Repository
	Factory
	Client          *sportmonks.Client
	Logger          *log.Logger
	PlayerProcessor stats.PlayerProcessor
	TeamProcessor   stats.TeamProcessor
	EventProcessor  event.Processor
}

func (p Processor) Process(command string, done chan bool) {
	switch command {
	case result:
		go p.allResults(done)
	default:
		p.Logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (p Processor) allResults(done chan bool) {
	ids, err := p.FixtureRepo.Ids()

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	results := make(chan sportmonks.Fixture, len(ids))

	go p.callClient(ids, results, done, &counter)
	go p.parseResults(results, done)
}

func (p Processor) callClient(ids []int, ch chan<- sportmonks.Fixture, done chan bool, c *int) {
	q := []string{"lineup,bench,stats,goals,substitutions"}

	for _, id := range ids {
		if _, err := p.Repository.GetByFixtureId(id); err != ErrNotFound {
			continue
		}

		if *c >= callLimit {
			p.Logger.Printf("Api call limited reached %d calls\n", *c)
			break
		}

		res, err := p.Client.FixtureById(id, q, 5)

		*c++

		if err != nil {
			p.Logger.Fatalf("Error when calling client '%s", err.Error())
			return
		}

		ch <- res.Data
	}

	close(ch)
}

func (p Processor) parseResults(ch <-chan sportmonks.Fixture, done chan bool) {
	for x := range ch {
		p.handleResult(x)
	}

	waitGroup.Wait()

	done <- true
}

func (p Processor) handleResult(fix sportmonks.Fixture) {
	waitGroup.Add(5)

	go func(stats []sportmonks.TeamStats) {
		p.handleTeams(stats)
		defer waitGroup.Done()
	}(fix.TeamStats.Data)

	go func(players []sportmonks.LineupPlayer) {
		p.handlePlayers(players, false)
		defer waitGroup.Done()
	}(fix.Lineup.Data)

	go func(players []sportmonks.LineupPlayer) {
		p.handlePlayers(players, true)
		defer waitGroup.Done()
	}(fix.Bench.Data)

	go func(goals []sportmonks.GoalEvent) {
		p.handleGoalEvents(goals)
		defer waitGroup.Done()
	}(fix.Goals.Data)

	go func(subs []sportmonks.SubstitutionEvent) {
		p.handleSubstitutionEvents(subs)
		defer waitGroup.Done()
	}(fix.Subs.Data)

	created := p.Factory.createResult(&fix)

	if err := p.Repository.Insert(created); err != nil {
		log.Printf("Error '%s' occurred when inserting Result struct: %+v\n,", err.Error(), created)
	}
}

func (p Processor) handleTeams(t []sportmonks.TeamStats) {
	for _, team := range t {
		p.TeamProcessor.ProcessTeamStats(&team)
	}
}

func (p Processor) handlePlayers(lineups []sportmonks.LineupPlayer, bench bool) {
	for _, player := range lineups {
		p.PlayerProcessor.ProcessPlayerStats(&player, bench)
	}
}

func (p Processor) handleGoalEvents(g []sportmonks.GoalEvent) {
	for _, goal := range g {
		waitGroup.Add(1)

		go func(e sportmonks.GoalEvent) {
			p.EventProcessor.ProcessGoalEvent(&e)
			defer waitGroup.Done()
		}(goal)
	}
}

func (p Processor) handleSubstitutionEvents(s []sportmonks.SubstitutionEvent) {
	for _, sub := range s {
		waitGroup.Add(1)
		
		go func(e sportmonks.SubstitutionEvent) {
			p.EventProcessor.ProcessSubstitutionEvent(&e)
			defer waitGroup.Done()
		}(sub)
	}
}
