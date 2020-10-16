package process

import (
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"strconv"
	"strings"
	"time"
)

const teamStatsByResultId = "team-stats:by-result-id"
const teamStatsBySeasonId = "team-stats:by-season-id"
const teamStatsToday = "team-stats:today"

type TeamStatsProcessor struct {
	teamStatsRepo app.TeamStatsRepository
	fixtureRepo   app.FixtureRepository
	requester     app.TeamStatsRequester
	clock         clockwork.Clock
	logger        *logrus.Logger
}

func (t TeamStatsProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case teamStatsByResultId:
		for _, id := range strings.Split(option, ",") {
			id, _ := strconv.Atoi(id)
			go t.processByID(done, uint64(id))
		}
	case teamStatsBySeasonId:
		id, _ := strconv.Atoi(option)
		go t.processSeason(done, uint64(id))
	case teamStatsToday:
		go t.processToday(done)
	default:
		t.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (t TeamStatsProcessor) processByID(done chan bool, id uint64) {
	fix, err := t.fixtureRepo.ByID(id)

	if err != nil {
		t.logger.Fatalf("Error when retrieving fixtures for ID: %d, %s", id, err.Error())
		return
	}

	ch := t.requester.TeamStatsByFixtureIDs([]uint64{fix.ID})

	go t.persistStats(ch, done)
}

func (t TeamStatsProcessor) processSeason(done chan bool, seasonID uint64) {
	query := app.FixtureRepositoryQuery{
		SeasonIDs: []uint64{seasonID},
	}

	fix, err := t.fixtureRepo.Get(query)

	if err != nil {
		t.logger.Fatalf("Error when retrieving fixtures for Season ID: %d, %s", seasonID, err.Error())
		return
	}

	var ids []uint64

	for _, f := range fix {
		ids = append(ids, f.ID)
	}

	ch := t.requester.TeamStatsByFixtureIDs(ids)

	go t.persistStats(ch, done)
}

func (t TeamStatsProcessor) processToday(done chan bool) {
	now := t.clock.Now()
	y, m, d := now.Date()

	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Date(y, m, d, 23, 59, 59, 59, now.Location())

	query := app.FixtureRepositoryQuery{
		DateTo:   &to,
		DateFrom: &from,
	}

	ids, err := t.fixtureRepo.GetIDs(query)

	if err != nil {
		t.logger.Fatalf("Error when retrieving fixture ids in event processor: %s", err.Error())
		return
	}

	ch := t.requester.TeamStatsByFixtureIDs(ids)

	go t.persistStats(ch, done)
}

func (t TeamStatsProcessor) persistStats(ch <-chan *app.TeamStats, done chan bool) {
	for result := range ch {
		t.persist(result)
	}

	done <- true
}

func (t TeamStatsProcessor) persist(x *app.TeamStats) {
	_, err := t.teamStatsRepo.ByFixtureAndTeam(x.FixtureID, x.TeamID)

	if err != nil {
		if err := t.teamStatsRepo.InsertTeamStats(x); err != nil {
			t.logger.Warningf("Error '%s' occurred when inserting team stats struct: %+v\n,", err.Error(), *x)
		}

		return
	}

	if err := t.teamStatsRepo.UpdateTeamStats(x); err != nil {
		t.logger.Warningf("Error '%s' occurred when updating team stats struct: %+v\n,", err.Error(), *x)
	}

	return
}

func NewTeamStatsProcessor(
	r app.TeamStatsRepository,
	f app.FixtureRepository,
	q app.TeamStatsRequester,
	c clockwork.Clock,
	log *logrus.Logger,
) *TeamStatsProcessor {
	return &TeamStatsProcessor{teamStatsRepo: r, fixtureRepo: f, requester: q, clock: c, logger: log}
}
