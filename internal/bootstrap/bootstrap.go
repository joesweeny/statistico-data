package bootstrap

import (
	"github.com/joesweeny/statshub/internal/config"
	"database/sql"
	"fmt"
	"github.com/joesweeny/sportmonks-go-client"
	"log"
	"os"
	"github.com/jonboulle/clockwork"
)

type Bootstrap struct {
	Config *config.Config
}

func (b Bootstrap) databaseConnection() *sql.DB {
	db := b.Config.Database

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)

	conn, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	return conn
}

func (b Bootstrap) sportmonksClient() (*sportmonks.Client, error) {
	s := b.Config.Services.SportsMonks
	return sportmonks.NewClient(s.BaseUri, s.ApiKey)
}

func logger() *log.Logger {
	return log.New(os.Stdout, "Error: ", 0)
}

func clock() clockwork.Clock {
	return clockwork.NewRealClock()
}