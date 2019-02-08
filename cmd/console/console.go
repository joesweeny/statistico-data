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

	var processor container.Processor

	switch *command {
	case Competition:
		processor = app.CompetitionProcessor()
		break
	case Country:
		processor = app.CountryProcessor()
		break
	case Fixture, FixtureCurrentSeason:
		processor = app.FixtureProcessor()
		break
	case Player:
		processor = app.PlayerProcessor()
	case Round, RoundCurrentSeason:
		processor = app.RoundProcessor()
		break
	case Season:
		processor = app.SeasonProcessor()
		break
	case Squad, SquadCurrentSeason:
		processor = app.SquadProcessor()
		break
	case Team, TeamCurrentSeason:
		processor = app.TeamProcessor()
		break
	case Venue, VenueCurrentSeason:
		processor = app.VenueProcessor()
		break
	default:
		fmt.Println("The command provided is not supported")
		os.Exit(1)
	}

	done := make(chan bool)

	start := time.Now()

	fmt.Printf("%s: Processing started for %s\n", start.String(), *command)

	go processor.Process(*command, done)

	<-done

	elapsed := time.Since(start)

	fmt.Printf("Processing complete for %s: Duration %s\n", *command, elapsed)

	os.Exit(0)
}
