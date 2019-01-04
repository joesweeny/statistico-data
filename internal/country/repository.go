package country

import (
	"github.com/joesweeny/statshub/internal/model"
	"github.com/satori/go.uuid"
)

type Repository interface {
	Insert(c model.Country) error
	Update(c model.Country) error
	GetById(u uuid.UUID) (model.Country, error)
	GetByExternalId(id int) (model.Country, error)
}