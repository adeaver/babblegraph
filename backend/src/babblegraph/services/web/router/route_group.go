package router

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type RouteGroup struct {
	Prefix string
	Routes []Route
}

type Route struct {
	Path    string
	Handler RouteHandler
}

type RouteHandler struct {
	HandleRequestBody RequestBodyHandler
	HandleRawRequest  func(http.ResponseWriter, *http.Request)
}

type RequestBodyHandler func(reqBody []byte) (_resp interface{}, _err error)

func (r Route) makeMuxRoute() func(http.ResponseWriter, *http.Request) {
	switch {
	case r.Handler.HandleRequestBody != nil:
		return makeMuxRouteForRequestBodyHandler(r.Handler.HandleRequestBody)
	case r.Handler.HandleRawRequest != nil:
		return r.Handler.HandleRawRequest
	default:
		panic(fmt.Sprintf("Route %s has no valid handler", r.Path))
	}
}

func makeMuxRouteForRequestBodyHandler(handler RequestBodyHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(fmt.Sprintf("ERROR: %s", err.Error()))
			writeErrorJSONResponse(w, errorResponse{
				Message: "Request is not valid",
			})
			return
		}
		resp, err := handler(body)
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
