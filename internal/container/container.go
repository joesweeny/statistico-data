package container

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-data/internal/config"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"log"
	"os"
	"time"
)

type Container struct {
	Clock               clockwork.Clock
	Config              *config.Config
	Database            *sql.DB
	Logger              *log.Logger
	NewLogger           *logrus.Logger
	SportMonksClient    *spClient.HTTPClient
}

func Bootstrap(config *config.Config) *Container {
	c := Container{
		Config: config,
	}

	c.Clock = clock()
	c.Database = databaseConnection(config)
	c.Logger = logger()
	c.NewLogger = newLogger()
	c.SportMonksClient = sportMonksClient(config)

	return &c
}

func databaseConnection(config *config.Config) *sql.DB {
	db := config.Database

	dsn := "host=%s port=%s user=%s " +
		"password=%s dbname=%s sslmode=disable"

	psqlInfo := fmt.Sprintf(dsn, db.Host, db.Port, db.User, db.Password, db.Name)

	conn, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	conn.SetMaxOpenConns(50)
	conn.SetMaxIdleConns(25)

	return conn
}

func sportMonksClient(config *config.Config) *spClient.HTTPClient {
	s := config.Services.SportsMonks

	return spClient.NewHTTPClient(s.ApiKey)
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	return logger
}

func clock() clockwork.Clock {
	return clockwork.NewRealClock()
}

func logger() *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("%s : Error: ", time.Now().Format(time.RFC3339)), 0)
}
