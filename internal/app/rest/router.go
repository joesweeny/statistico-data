package rest

import (
	"github.com/julienschmidt/httprouter"
)

func Router() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", routePath)
	router.GET("/healthcheck", healthCheck)
	router.GET("/competition/:competition/season/:season/fixture", competitionFixtures)

	return router
}
