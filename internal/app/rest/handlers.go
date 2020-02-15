package rest

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RoutePath(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, "./opt/api/index.html")
}

func HealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Healthcheck OK")
}

func RenderApiDocs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, "./opt/api/openapi.json")
}
