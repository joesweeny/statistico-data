package team

import "github.com/statistico/statistico-data/internal/model"

type Repository interface {
	Insert(t *model.Team) error
	Update(t *model.Team) error
	GetById(id int) (*model.Team, error)
}
