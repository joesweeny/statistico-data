package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
)

const season = "season"

// SeasonProcessor is used to process data from an external data source to this applications
// chosen data store.
type SeasonProcessor struct {
	repository app.SeasonRepository
	requester app.SeasonRequester
	logger *logrus.Logger
}

// Process fetches data from external an external data source using the SeasonRequester
// before persisting to the storage engine using the SeasonRepository.
func (s SeasonProcessor) Process(command string, option string, done chan bool) {
	if command != season {
		s.logger.Fatalf("Command %s is not supported", command)
	}

	ch := s.requester.Seasons()

	go s.persistSeasons(ch, done)
}

func (s SeasonProcessor) persistSeasons(ch <-chan *app.Season, done chan bool) {
	for season := range ch {
		s.persist(season)
	}

	done <- true
}

func (s SeasonProcessor) persist(a *app.Season) {
	_, err := s.repository.ByID(a.ID)

	if err != nil {
		if err := s.repository.Insert(a); err != nil {
			s.logger.Warningf("Error '%s' occurred when inserting season struct: %+v\n,", err.Error(), *c)
		}

		return
	}

	if err := s.repository.Update(a); err != nil {
		s.logger.Warningf("Error '%s' occurred when updating season struct: %+v\n,", err.Error(), *c)
	}

	return
}

func NewSeasonProcessor(r app.SeasonRepository, a app.SeasonRequester, log *logrus.Logger) *SeasonProcessor {
	return &SeasonProcessor{repository: r, requester: a, logger: log}
}