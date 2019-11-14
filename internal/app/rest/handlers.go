package rest

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func routePath(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Welcome to the Statistico Data API")
}

func healthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Healthcheck OK")
}

func competitionFixtures(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	comp := ps.ByName("competition")
	sea := ps.ByName("season")

	query := r.URL.Query()

	after := query.Get("date_after")
	before := query.Get("date_before")

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Competition ID %s. Season %s \n", comp, sea)

	if after != "" {
		_, _ = fmt.Fprintf(w, "After date %s", after)
	}

	if before != "" {
		_, _ = fmt.Fprintf(w, "Before date %s", before)
	}
}
