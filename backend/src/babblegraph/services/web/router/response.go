package router

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func writeErrorJSONResponse(w http.ResponseWriter, body errorResponse) {
	w.WriteHeader(http.StatusBadRequest)
	writeJSONResponse(w, body)
}

func writeJSONResponse(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
