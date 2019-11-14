package rest

import (
	"encoding/json"
	"net/http"
)

type jsend struct {
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type errorMessage struct {
	Message string `json:"message"`
	Code    int `json:"code"`
}
func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(response)
}

func jsendSuccessResponse(w http.ResponseWriter, payload interface{}) {
	response := jsend{
		Message: "success",
		Data:    payload,
	}

	jsonResponse(w, 200, response)
}

func jsendFailResponse(w http.ResponseWriter, status int, error error) {
	response := jsend{
		Message: "fail",
		Data:    []errorMessage{
			{
				Message: error.Error(),
				Code:    1,
			},
		},
	}

	jsonResponse(w, status, response)
}

func jsendErrorResponse(w http.ResponseWriter, status int, error error) {
	response := jsend{
		Message: "error",
		Data:    []errorMessage{
			{
				Message: error.Error(),
				Code:    1,
			},
		},
	}

	jsonResponse(w, status, response)
}