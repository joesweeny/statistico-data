package main

import (
	"flag"
	"fmt"
	"github.com/joesweeny/statshub/internal/config"
	"github.com/joesweeny/statshub/internal/container"
	"os"
	"time"
)

var command = flag.String("command", "", "Provide the model name to process")

func main() {
	app := container.Bootstrap(config.GetConfig())

	flag.Parse()

	var service container.Service

	switch *command {
	case Competition:
		service = app.CompetitionService()
		break
	case Country:
		service = app.CountryService()
		break
	case Fixture, FixtureCurrentSeason:
		service = app.FixtureService()
		break
	case Round, RoundCurrentSeason:
		service = app.RoundService()
		break
	case Season:
		service = app.SeasonService()
		break
	case Squad, SquadCurrentSeason:
		service = app.SquadService()
		break
	case Team, TeamCurrentSeason:
		service = app.TeamService()
		break
	case Venue, VenueCurrentSeason:
		service = app.VenueService()
		break
	default:
		fmt.Println("The command provided is not supported")
		os.Exit(1)
	}

	done := make(chan bool)

	start := time.Now()

	fmt.Printf("%s: Processing started for %s\n", start.String(), *command)

	service.Process(*command, done)

	<-done

	elapsed := time.Since(start)

	fmt.Printf("Processing complete for %s: Duration %s\n", *command, elapsed)

	os.Exit(0)
}

func fail(model *string, err error) {
	if err != nil {
		fmt.Printf("Error when processing %s: %s\n", *model, err.Error())
		os.Exit(1)
	}
}
