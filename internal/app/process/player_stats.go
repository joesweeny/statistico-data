package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
	"time"
)

const playerStatsByDate = "player-stats:by-date"
const playerStatsBySeasonId = "player-stats:by-season-id"
const playerStatsByCompetitionId = "player-stats:by-competition-id"

type PlayerStatsProcessor struct {
	playerStatsRepo app.PlayerStatsRepository
	competitionRepo app.CompetitionRepository
	seasonRepo      app.SeasonRepository
	requester       app.PlayerStatRequester
	clock           clockwork.Clock
	logger          *logrus.Logger
}

func (p PlayerStatsProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case playerStatsByDate:
		go p.processByDate(option, done)
	case playerStatsBySeasonId:
		id, _ := strconv.Atoi(option)
		go p.processSeason(uint64(id), done)
	case playerStatsByCompetitionId:
		id, _ := strconv.Atoi(option)
		go p.processCompetition(uint64(id), done)
	default:
		p.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (p PlayerStatsProcessor) processByDate(date string, done chan bool) {
	d, err := time.Parse("2006-01-02", date)

	if err != nil {
		p.logger.Fatalf("Error parsing date in player stats processor: %s", err.Error())
		return
	}

	ids, err := p.competitionRepo.IDs()

	if err != nil {
		p.logger.Fatalf("Error fetching competition IDs in player stats processor: %s", err.Error())
		return
	}

	ch := p.requester.PlayerStatsByDate(d, ids)

	go p.persistStats(ch, done)
}

func (p PlayerStatsProcessor) processSeason(seasonID uint64, done chan bool) {
	ch := p.requester.PlayerStatsBySeasonIDs([]uint64{seasonID})

	go p.persistStats(ch, done)
}

func (p PlayerStatsProcessor) processCompetition(competitionID uint64, done chan bool) {
	seasons, err := p.seasonRepo.ByCompetitionId(competitionID, "name_asc")

	if err != nil {
		p.logger.Fatalf("Error fetching seasons in player stats processor: %s", err.Error())
		return
	}

	var ids []uint64

	for _, s := range seasons {
		ids = append(ids, s.ID)
	}

	ch := p.requester.PlayerStatsBySeasonIDs(ids)

	go p.persistStats(ch, done)
}

func (p PlayerStatsProcessor) persistStats(ch <-chan *app.PlayerStats, done chan bool) {
	for stats := range ch {
		p.persist(stats)
	}

	done <- true
}

func (p PlayerStatsProcessor) persist(x *app.PlayerStats) {
	_, err := p.playerStatsRepo.ByFixtureAndPlayer(x.FixtureID, x.PlayerID)

	if err != nil {
		if err := p.playerStatsRepo.Insert(x); err != nil {
			p.logger.Warningf("Error '%s' occurred when inserting player stats struct: %+v\n,", err.Error(), *x)
		}

		return
	}

	if err := p.playerStatsRepo.Update(x); err != nil {
		p.logger.Warningf("Error '%s' occurred when updating player stats struct: %+v\n,", err.Error(), *x)
	}

	return
}

func NewPlayerStatsProcessor(
	r app.PlayerStatsRepository,
	c app.CompetitionRepository,
	s app.SeasonRepository,
	q app.PlayerStatRequester,
	cl clockwork.Clock,
	log *logrus.Logger,
) *PlayerStatsProcessor {
	return &PlayerStatsProcessor{
		playerStatsRepo: r,
		competitionRepo: c,
		seasonRepo: s,
		requester: q,
		clock: cl,
		logger: log,
	}
}
