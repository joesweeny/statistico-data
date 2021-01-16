package bootstrap

import (
	"database/sql"
	"fmt"
	"github.com/evalphobia/logrus_sentry"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	spClient "github.com/statistico/statistico-sportmonks-go-client"
	understat "github.com/statistico/statistico-understat-parser"
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
	UnderstatParser  *understat.Parser
}

func BuildContainer(config *Config) *Container {
	c := Container{
		Config: config,
	}

	c.Clock = clock()
	c.Database = databaseConnection(config)
	c.Logger = logger(config)
	c.SportMonksClient = sportMonksClient(config)
	c.UnderstatParser = understatParser(config)

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
			Timeout: 10 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: 30 * time.Second,
		TLSHandshakeTimeout: 15 * time.Second,
	}

	client := &http.Client{
		Timeout:   time.Second * 30,
		Transport: trans,
	}

	c.SetHTTPClient(client)

	return c
}

func understatParser(config *Config) *understat.Parser {
	return &understat.Parser{BaseURL: config.Understat.BaseURL}
}

func logger(config *Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	tags := map[string]string{
		"application": "statistico-data",
	}

	levels := []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}

	hook, err := logrus_sentry.NewWithTagsSentryHook(config.Sentry.DSN, tags, levels)

	if err == nil {
		hook.Timeout = 20 * time.Second
		hook.StacktraceConfiguration.Enable = true
		hook.StacktraceConfiguration.IncludeErrorBreadcrumb = true
		logger.AddHook(hook)
	}

	return logger
}

func clock() clockwork.Clock {
	return clockwork.NewRealClock()
}
