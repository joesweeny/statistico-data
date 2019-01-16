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

type Bootstrap struct {}

func getConfig() *config.Config {
	return config.GetConfig()
}

func databaseConnection() *sql.DB {
	c := getConfig()
	db := c.Database

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)

	conn, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	return conn
}

func sportmonksClient() (*sportmonks.Client, error) {
	c := getConfig()
	s := c.Services.SportsMonks
	return sportmonks.NewClient(s.BaseUri, s.ApiKey)
}

func logger() *log.Logger {
	return log.New(os.Stdout, "Error: ", 0)
}

func clock() clockwork.Clock {
	return clockwork.NewRealClock()
}