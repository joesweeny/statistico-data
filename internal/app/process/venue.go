package process

import (
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/app"
	"github.com/statistico/statistico-data/internal/season"
)

const venue = "venue"
const venueCurrentSeason = "venue:current-season"

type VenueProcessor struct {
	venueRepo  app.VenueRepository
	seasonRepo season.Repository
	requester  app.VenueRequester
	logger     *logrus.Logger
}

func (p VenueProcessor) Process(command string, option string, done chan bool) {
	switch command {
	case venue:
		go p.processAllSeasons(done)
	case venueCurrentSeason:
		go p.processCurrentSeason(done)
	default:
		p.logger.Fatalf("Command %s is not supported", command)
		return
	}
}

func (p VenueProcessor) processAllSeasons(done chan bool) {
	ids, err := p.seasonRepo.Ids()

	if err != nil {
		p.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := p.requester.VenuesBySeasonIDs(ids)

	go p.persistVenues(ch, done)
}

func (p VenueProcessor) processCurrentSeason(done chan bool) {
	ids, err := p.seasonRepo.CurrentSeasonIds()

	if err != nil {
		p.logger.Fatalf("Error when retrieving season ids: %s", err.Error())
		return
	}

	ch := p.requester.VenuesBySeasonIDs(ids)

	go p.persistVenues(ch, done)
}

func (p VenueProcessor) persistVenues(ch <-chan *app.Venue, done chan bool) {
	for venue := range ch {
		p.persist(venue)
	}

	done <- true
}

func (p VenueProcessor) persist(v *app.Venue) {
	_, err := p.venueRepo.GetById(v.ID)

	if err != nil {
		if err := p.venueRepo.Insert(v); err != nil {
			p.logger.Warningf("Error '%s' occurred when inserting venue struct: %+v\n,", err.Error(), *v)
		}

		return
	}

	if err := p.venueRepo.Update(v); err != nil {
		p.logger.Warningf("Error '%s' occurred when updating venue struct: %+v\n,", err.Error(), *v)
	}

	return
}

func NewVenueProcessor(r app.VenueRepository, s season.Repository, v app.VenueRequester, log *logrus.Logger) *VenueProcessor {
	return &VenueProcessor{venueRepo: r, seasonRepo: s, requester: v, logger: log}
}
