package player_stats

import (
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/fixture"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const playerStatsByResultId = "player-stats:by-result-id"
const playerStatsBySeasonId = "player-stats:by-season-id"
const playerStatsToday = "player-stats:today"
const callLimit = 1500

var counter int
var waitGroup sync.WaitGroup

type Processor struct {
	PlayerRepository
	PlayerFactory
	Logger *log.Logger
	FixtureRepo fixture.Repository
	Client          *sportmonks.Client
}

func (p Processor) Process(command string, option string, done chan bool) {
	switch command {
	case playerStatsByResultId:
		for _, id := range strings.Split(option, ",") {
			id, _ := strconv.Atoi(id)
			go p.byId(done, id)
		}
	case playerStatsBySeasonId:
		id, _ := strconv.Atoi(option)
		go p.bySeasonId(done, id)
	case playerStatsToday:
		go p.statsToday(done)
	default:
		p.Logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (p Processor) byId(done chan bool, id int) {
	fix, err := p.FixtureRepo.ById(uint64(id))

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Fixture ID: %d, %s", id, err.Error())
		return
	}

	ids := []int{fix.ID}

	go p.processStats(ids, done)
}

func (p Processor) bySeasonId(done chan bool, id int) {
	// Adding a Clock.Now() here is a bit hacky. Redo by dynamically handling this
	fix, err := p.FixtureRepo.BySeasonId(int64(id), p.Clock.Now())

	if err != nil {
		p.Logger.Fatalf("Error when retrieving fixtures for Season ID: %d, %s", id, err.Error())
		return
	}

	var ids []int

	for _, f := range fix {
		ids = append(ids, f.ID)
	}

	go p.processStats(ids, done)
}

func (p Processor) statsToday(done chan bool) {
	now := p.Clock.Now()
	y, m, d := now.Date()

	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

	ids, err := p.FixtureRepo.IdsBetween(from, to)

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	go p.processStats(ids, done)
}

func (p Processor) processStats(ids []int, done chan bool) {
	results := make(chan sportmonks.Fixture, len(ids))

	go p.callClient(ids, results, done, &counter)
	go p.parseStats(results, done)
}

func (p Processor) callClient(ids []int, ch chan<- sportmonks.Fixture, done chan bool, c *int) {
	q := []string{"lineup,bench"}

	for _, id := range ids {
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

func (p Processor) parseStats(ch <-chan sportmonks.Fixture, done chan bool) {
	for x := range ch {
		p.handleStats(x)
	}

	waitGroup.Wait()

	done <- true
}

func (p Processor) handleStats(fix sportmonks.Fixture) {
	for _, player := range fix.Lineup.Data {
		waitGroup.Add(1)

		go func(stats sportmonks.LineupPlayer) {
			p.processPlayerStats(&stats, false)
			defer waitGroup.Done()
		}(player)
	}

	for _, player := range fix.Bench.Data {
		waitGroup.Add(1)

		go func(stats sportmonks.LineupPlayer) {
			p.processPlayerStats(&stats, true)
			defer waitGroup.Done()
		}(player)
	}
}

func (p Processor) processPlayerStats(s *sportmonks.LineupPlayer, isSub bool) {
	x, err := p.PlayerRepository.ByFixtureAndPlayer(uint64(s.FixtureID), uint64(s.PlayerID))

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
