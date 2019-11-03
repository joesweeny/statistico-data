package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
)

const team = "team"
const teamCurrentSeason = "team:current-season"

// TeamProcessor fetches data from external data source using the TeamRequester
// before persisting to the storage engine using the TeamRepository.
type TeamProcessor struct {
	teamRepo app.TeamRepository
	seasonRepo app.SeasonRepository
	requester app.TeamRequester
	logger     *logrus.Logger
}

func (t TeamProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case team:
		go t.processAllSeasons(done)
	case teamCurrentSeason:
		go t.processCurrentSeason(done)
	default:
		t.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (t TeamProcessor) processAllSeasons(done chan bool) {
	ids, err := t.seasonRepo.IDs()

	if err != nil {
		t.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := t.requester.TeamsBySeasonIDs(ids)

	go t.persistTeams(ch, done)
}

func (t TeamProcessor) processCurrentSeason(done chan bool) {
	ids, err := t.seasonRepo.CurrentSeasonIDs()

	if err != nil {
		t.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := t.requester.TeamsBySeasonIDs(ids)

	go t.persistTeams(ch, done)
}

func (t TeamProcessor) persistTeams(ch <-chan *app.Team, done chan bool) {
	for team := range ch {
		t.persist(team)
	}

	done <- true
}

func (t TeamProcessor) persist(x *app.Team) {
	_, err := t.teamRepo.ByID(x.ID)

	if err != nil {
		if err := t.teamRepo.Insert(x); err != nil {
			t.logger.Warningf("Error '%s' occurred when inserting team struct: %+v\n,", err.Error(), *x)
		}

		return
	}

	if err := t.teamRepo.Update(x); err != nil {
		t.logger.Warningf("Error '%s' occurred when updating team struct: %+v\n,", err.Error(), *x)
	}

	return
}

func NewTeamProcessor(t app.TeamRepository, s app.SeasonRepository, r app.TeamRequester, log *logrus.Logger) *TeamProcessor {
	return &TeamProcessor{teamRepo: t, seasonRepo: s, requester: r, logger: log}
}
