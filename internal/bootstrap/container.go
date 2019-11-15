package bootstrap

import (
	"database/sql"
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	"net"
	"net/http"
	"os"
	"time"
)

type Container struct {
	Clock            clockwork.Clock
	Config           *Config
	Database         *sql.DB
	Logger           *logrus.Logger
	SportMonksClient *spClient.HTTPClient
}

func BuildContainer(config *Config) *Container {
	c := Container{
		Config: config,
	}

	c.Clock = clock()
	c.Database = databaseConnection(config)
	c.Logger = logger()
	c.SportMonksClient = sportMonksClient(config)

	return &c
}

func databaseConnection(config *Config) *sql.DB {
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

func sportMonksClient(config *Config) *spClient.HTTPClient {
	s := config.Services.SportsMonks

	c := spClient.NewDefaultHTTPClient(s.ApiKey)

	trans := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Timeout:   time.Second * 30,
		Transport: trans,
	}

	c.SetHTTPClient(client)

	return c
}

func logger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	return logger
}

func clock() clockwork.Clock {
	return clockwork.NewRealClock()
}
