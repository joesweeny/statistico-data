package rest

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
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

	json, err := os.Open("./opt/api/openapi.json")

	defer json.Close()

	if err != nil {
		_, _ = w.Write([]byte("Internal server error"))
		return
	}

	contents, _ := ioutil.ReadAll(json)
	_, _ = w.Write(contents)
}
