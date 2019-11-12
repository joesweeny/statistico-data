package main

import (
	"flag"
	"fmt"
	"github.com/statistico/statistico-data/internal/bootstrap"
	"os"
	"time"
)

var command = flag.String("command", "", "Provide the command name to process")
var option = flag.String("option", "", "Optional parameter to pass to command")

func main() {
	app := bootstrap.BuildContainer(bootstrap.BuildConfig())

	flag.Parse()

	var processor bootstrap.Processor

	switch *command {
	case Competition:
		processor = app.CompetitionProcessor()
		break
	case Country:
		processor = app.CountryProcessor()
		break
	case EventsByResultId, EventsBySeasonId, EventsToday:
		processor = app.EventProcessor()
		break
	case Fixture, FixtureCurrentSeason:
		processor = app.FixtureProcessor()
		break
	case Player:
		processor = app.PlayerProcessor()
		break
	case PlayerStatsByResultId, PlayerStatsBySeasonId, PlayerStatsToday:
		processor = app.PlayerStatsProcessor()
		break
	case Result, ResultById, ResultBySeasonId, ResultToday:
		processor = app.ResultProcessor()
		break
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
	case TeamStatsByResultId, TeamStatsBySeasonId, TeamStatsToday:
		processor = app.TeamStatsProcessor()
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

	go processor.Process(*command, *option, done)

	<-done

	elapsed := time.Since(start)

	fmt.Printf("Processing complete for %s: Duration %s\n", *command, elapsed)

	os.Exit(0)
}
