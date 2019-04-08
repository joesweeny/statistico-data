package event

import "github.com/statistico/statistico-data/internal/model"

type Repository interface {
	InsertGoalEvent(m *model.GoalEvent) error
	InsertSubstitutionEvent(m *model.SubstitutionEvent) error
	GoalEventById(id int) (*model.GoalEvent, error)
	SubstitutionEventById(id int) (*model.SubstitutionEvent, error)
}
