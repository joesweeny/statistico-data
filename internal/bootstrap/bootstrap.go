package bootstrap

import (
	"database/sql"
	"fmt"
	"github.com/joesweeny/sportmonks-go-client"
	"github.com/joesweeny/statshub/internal/config"
	"github.com/jonboulle/clockwork"
	"log"
	"os"
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

	conn.SetMaxOpenConns(50)
	conn.SetMaxIdleConns(25)

	return conn
}

func (b Bootstrap) sportmonksClient() (*sportmonks.Client, error) {
	s := b.Config.Services.SportsMonks
	return sportmonks.NewClient(s.BaseUri, s.ApiKey, logger())
}

func logger() *log.Logger {
	return log.New(os.Stdout, "Error: ", 0)
}

func clock() clockwork.Clock {
	return clockwork.NewRealClock()
}
