package player

import "github.com/joesweeny/statshub/internal/model"

type Repository interface {
	Insert(m *model.Player) error
	Update(m *model.Player) error
	Id(id int) (*model.Player, error)
}
