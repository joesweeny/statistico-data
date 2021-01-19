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
	case competition:
		processor = app.CompetitionProcessor()
		break
	case country:
		processor = app.CountryProcessor()
		break
	case events, eventsCurrentSeason, eventsBySeasonId:
		processor = app.EventProcessor()
		break
	case fixturesCurrentSeason, fixturesBySeasonId, fixturesByCompetitionId:
		processor = app.FixtureProcessor()
		break
	case fixtureXG, fixtureXGCurrentSeason:
		processor = app.FixtureTeamXGProcessor()
	case player:
		processor = app.PlayerProcessor()
		break
	case playerStatsByDate, playerStatsBySeasonId, playerStatsByCompetitionId:
		processor = app.PlayerStatsProcessor()
		break
	case resultsCurrentSeason, resultsBySeasonId, resultsByCompetitionId:
		processor = app.ResultProcessor()
		break
	case round, roundCurrentSeason:
		processor = app.RoundProcessor()
		break
	case season:
		processor = app.SeasonProcessor()
		break
	case squad, squadCurrentSeason:
		processor = app.SquadProcessor()
		break
	case team, teamCurrentSeason:
		processor = app.TeamProcessor()
		break
	case teamStatsByDate, teamStatsBySeasonId, teamStatsByCompetitionId:
		processor = app.TeamStatsProcessor()
		break
	case venue, venueCurrentSeason:
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
