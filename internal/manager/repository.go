package manager

import "github.com/joesweeny/statshub/internal/model"

type Repository interface {
	Insert(m *model.Manager) error
	Update(m *model.Manager) error
	Id(id int) (*model.Manager, error)
}
