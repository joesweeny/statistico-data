package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
)

const competition = "competition"

// CompetitionProcessor is used to process data from an external data source to this applications
// chosen data store.
type CompetitionProcessor struct {
	repository app.CompetitionRepository
	requester app.CompetitionRequester
	logger *logrus.Logger
}

// Process fetches data from external an external data source using the CompetitionRequester
// before persisting to the storage engine using the CompetitionRepository.
func (p CompetitionProcessor) Process(command string, option string, done chan bool) {
	if command != competition {
		p.logger.Fatalf("Command %s is not supported", command)
	}

	ch := p.requester.Competitions()

	go p.persistCompetitions(ch, done)
}

func (p CompetitionProcessor) persistCompetitions(ch <-chan *app.Competition, done chan bool) {
	for competition := range ch {
		p.persist(competition)
	}

	done <- true
}

func (p CompetitionProcessor) persist(c *app.Competition) {
	_, err := p.repository.ByID(c.ID)

	if err != nil {
		if err := p.repository.Insert(c); err != nil {
			p.logger.Warningf("Error '%s' occurred when inserting competition struct: %+v\n,", err.Error(), *c)
		}

		return
	}

	if err := p.repository.Update(c); err != nil {
		p.logger.Warningf("Error '%s' occurred when updating competition struct: %+v\n,", err.Error(), *c)
	}

	return
}

func NewCompetitionProcessor(r app.CompetitionRepository, c app.CompetitionRequester, log *logrus.Logger) *CompetitionProcessor {
	return &CompetitionProcessor{repository: r, requester: c, logger: log}
}
