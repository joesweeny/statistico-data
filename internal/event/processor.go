package event

import (
	"github.com/joesweeny/sportmonks-go-client"
	"log"
)

type Processor struct {
	Repository
	Factory
	Logger *log.Logger
}

func (p Processor) ProcessGoalEvent(s *sportmonks.GoalEvent) {
	if _, err := p.Repository.GoalEventById(s.ID); err != ErrNotFound {
		return
	}

	created := p.Factory.createGoalEvent(s)

	if err := p.InsertGoalEvent(created); err != nil {
		log.Printf("Error '%s' occurred when inserting Goal Event struct: %+v\n,", err.Error(), created)
	}

	return
}

func (p Processor) ProcessSubstitutionEvent(s *sportmonks.SubstitutionEvent) {
	if _, err := p.Repository.SubstitutionEventById(s.ID); err != ErrNotFound {
		return
	}

	created := p.Factory.createSubstitutionEvent(s)

	if err := p.InsertSubstitutionEvent(created); err != nil {
		log.Printf("Error '%s' occurred when inserting Substitution Event struct: %+v\n,", err.Error(), created)
	}

	return
}
