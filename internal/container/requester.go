package container

import "github.com/statistico/statistico-data/internal/app/sportmonks"

func (c Container) CountryRequester() *sportmonks.CountryRequester {
	return sportmonks.NewCountryRequester(c.SportMonksClient)
}
