package container

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/sportmonks-go-client"
	"github.com/statistico/statistico-data/internal/config"
	"log"
	"os"
	"time"
)

type Container struct {
	Clock            clockwork.Clock
	Config           *config.Config
	Database         *sql.DB
	Logger           *log.Logger
	SportMonksClient *sportmonks.Client
}

func Bootstrap(config *config.Config) *Container {
	c := Container{
		Config: config,
	}

	c.Clock = clock()
	c.Database = databaseConnection(config)
	c.Logger = logger()
	c.SportMonksClient = sportmonksClient(config)

	return &c
}

func databaseConnection(config *config.Config) *sql.DB {
	db := config.Database

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

func sportmonksClient(config *config.Config) *sportmonks.Client {
	s := config.Services.SportsMonks

	client, err := sportmonks.NewClient(s.BaseUri, s.ApiKey, logger())

	if err != nil {
		panic(err)
	}

	return client
}

func logger() *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("%s : Error: ", time.Now().Format(time.RFC3339)), 0)
}

func clock() clockwork.Clock {
	return clockwork.NewRealClock()
}
