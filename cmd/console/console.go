package main

import (
	"flag"
	"fmt"
	"github.com/joesweeny/statshub/internal/bootstrap"
	"github.com/joesweeny/statshub/internal/config"
	"os"
	"time"
)

const competition = "competition"
const country = "country"
const fixture = "fixture"
const fixtureCurrentSeason = "fixture:current-season"
const round = "round"
const roundCurrentSeason = "round:current-season"
const season = "season"
const squad = "squad"
const team = "team"
const teamCurrentTeam = "team:current-season"
const venue = "venue"
const venueCurrentSeason = "venue:current-season"

var option = flag.String("option", "", "Provide the model name to process")

func main() {
	app := bootstrap.Bootstrap{Config: config.GetConfig()}

	flag.Parse()

	var service bootstrap.Service

	switch *option {
	case competition:
		service = app.CompetitionService()
		break
	case country:
		service = app.CountryService()
		break
	case fixture:
		service = app.FixtureService()
		break
	case round:
		service = app.RoundService()
		break
	case season:
		service = app.SeasonService()
	case squad:
		service = app.SquadService()
		break
	case team:
		service = app.TeamService()
	case venue:
		service = app.VenueService()
	default:
		fmt.Println("The option provided is not supported")
		os.Exit(1)
	}

	start := time.Now()

	if err := service.Process(); err != nil {
		fail(option, err)
		os.Exit(1)
	}

	fmt.Printf("Processing complete for %s\n", *option)

	elapsed := time.Since(start)

	fmt.Printf("%s command took %s\n", *option, elapsed)

	os.Exit(0)
}

func fail(model *string, err error) {
	if err != nil {
		fmt.Printf("Error when processing %s: %s\n", *model, err.Error())
		os.Exit(1)
	}
}
