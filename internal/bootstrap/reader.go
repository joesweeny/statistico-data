package bootstrap

import "github.com/statistico/statistico-football-data/internal/app/performance/postgres"

func (c Container) StatReader() *postgres.StatReader {
	return postgres.NewStatReader(c.Database)
}
