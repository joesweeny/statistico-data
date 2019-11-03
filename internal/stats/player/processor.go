package player_stats

import (
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/model"
	"log"
	"strconv"
	"strings"
	"time"
)

const playerStatsByResultId = "player-stats:by-result-id"
const playerStatsBySeasonId = "player-stats:by-season-id"
const playerStatsToday = "player-stats:today"
const callLimit = 1500

var counter int

type Processor struct {
	PlayerRepository
	PlayerFactory
	Logger      *log.Logger
	FixtureRepo app.FixtureRepository
	Client      *sportmonks.Client
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
	fix, err := p.FixtureRepo.ByID(uint64(id))

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Fixture ID: %d, %s", id, err.Error())
		return
	}

	ids := []uint64{fix.ID}

	go p.processStats(ids, done)
}

func (p Processor) bySeasonId(done chan bool, id int) {
	// Adding a Clock.Now() here is a bit hacky. Redo by dynamically handling this
	fix, err := p.FixtureRepo.BySeasonID(uint64(id), p.Clock.Now())

	if err != nil {
		p.Logger.Fatalf("Error when retrieving fixtures for Season ID: %d, %s", id, err.Error())
		return
	}

	var ids []uint64

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

	ids, err := p.FixtureRepo.IDsBetween(from, to)

	if err != nil {
		p.Logger.Fatalf("Error when retrieving Season IDs: %s", err.Error())
		return
	}

	go p.processStats(ids, done)
}

func (p Processor) processStats(ids []uint64, done chan bool) {
	results := make(chan sportmonks.Fixture, len(ids))
	stats := make(chan *model.PlayerStats, 1000)

	go p.callClient(ids, results, done, &counter)
	go p.parseStats(results, stats)
	go p.processPlayerStats(stats, done)
}

func (p Processor) callClient(ids []uint64, ch chan<- sportmonks.Fixture, done chan bool, c *int) {
	q := []string{"lineup,bench"}

	for _, id := range ids {
		if *c >= callLimit {
			p.Logger.Printf("Api call limited reached %d calls\n", *c)
			break
		}

		res, err := p.Client.FixtureById(int(id), q, 5)

		*c++

		if err != nil {
			p.Logger.Fatalf("Error when calling client '%s", err.Error())
			return
		}

		ch <- res.Data
	}

	close(ch)
}

func (p Processor) parseStats(ch <-chan sportmonks.Fixture, stats chan<- *model.PlayerStats) {
	for x := range ch {
		p.handleStats(x, stats)
	}

	close(stats)
}

func (p Processor) handleStats(fix sportmonks.Fixture, ch chan<- *model.PlayerStats) {
	for _, player := range fix.Lineup.Data {
		ch <- p.PlayerFactory.createPlayerStats(&player, false)
	}

	for _, player := range fix.Bench.Data {
		ch <- p.PlayerFactory.createPlayerStats(&player, true)
	}
}

func (p Processor) processPlayerStats(ch <-chan *model.PlayerStats, done chan bool) {
	for s := range ch {
		_, err := p.PlayerRepository.ByFixtureAndPlayer(uint64(s.FixtureID), uint64(s.PlayerID))

		if err == ErrNotFound {
			if err := p.PlayerRepository.InsertPlayerStats(s); err != nil {
				log.Printf("Error '%s' occurred when inserting Player Stats struct: %+v\n,", err.Error(), s)
			}

			continue
		}

		if err := p.PlayerRepository.UpdatePlayerStats(s); err != nil {
			log.Printf("Error '%s' occurred when Updating Player Stats struct: %+v\n,", err.Error(), s)
		}
	}

	done <- true
}
