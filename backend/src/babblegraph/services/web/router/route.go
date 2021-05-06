package router

import (
	"babblegraph/model/useraccounts"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Route struct {
	Path                        string
	Handler                     RouteHandler
	ShouldLogBody               bool
	TrackEventWithID            *string
	WithValidAuthorizationLevel []useraccounts.SubscriptionLevel
}

type RouteHandler func(reqBody []byte) (_resp interface{}, _err error)

func makeMuxRouter(processRequest RouteHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(fmt.Sprintf("ERROR: %s", err.Error()))
			writeErrorJSONResponse(w, errorResponse{
				Message: "Request is not valid",
			})
			return
		}
		resp, err := processRequest(body)
		if err != nil {
			log.Println(fmt.Sprintf("ERROR: %s", err.Error()))
			writeErrorJSONResponse(w, errorResponse{
				Message: "Error processing request",
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		writeJSONResponse(w, resp)
	}
}
