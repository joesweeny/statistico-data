package country

import (
	"github.com/joesweeny/statshub/internal/model"
	"github.com/satori/go.uuid"
)

type repository interface {
	insert(c model.Country) error
	update(c model.Country) error
	getById(u uuid.UUID) (model.Country, error)
	getByExternalId(id int) (model.Country, error)
}