package router

import (
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Route struct {
	Path             string
	Handler          RouteHandler
	ShouldLogBody    bool
	TrackEventWithID *string
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

type AuthenticatedRoute struct {
	Path                     string
	Handler                  AuthRouteHandler
	ShouldLogBody            bool
	TrackEventWithID         *string
	ValidAuthorizationLevels []useraccounts.SubscriptionLevel
}

type AuthRouteHandler func(userID users.UserID, reqBody []byte) (_resp interface{}, _err error)

func makeAuthenticatedMuxRouter(userID users.UserID, processRequest AuthRouteHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(fmt.Sprintf("ERROR: %s", err.Error()))
			writeErrorJSONResponse(w, errorResponse{
				Message: "Request is not valid",
			})
			return
		}
		resp, err := processRequest(userID, body)
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
