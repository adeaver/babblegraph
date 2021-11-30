package router

import (
	"fmt"
	"log"
	"net/http"
)

type RouteGroup struct {
	Prefix string
	Routes []Route
}

type Route struct {
	Path    string
	Handler RequestHandler
}

type RequestHandler func(r *Request) (interface{}, error)

func (r Route) makeMuxRoute() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		wrappedRequest := &Request{
			r: req,
		}
		status := http.StatusOK
		resp, err := r.Handler(wrappedRequest)
		switch {
		case err != nil:
			log.Println(fmt.Sprintf("ERROR: %s", err.Error()))
			writeErrorJSONResponse(w, errorResponse{
				Message: "Error processing request",
			})
			return
		case wrappedRequest.respStatus != nil:
			status = *wrappedRequest.respStatus
		}
		if len(wrappedRequest.respCookies) != 0 {
			for _, cookie := range wrappedRequest.respCookies {
				http.SetCookie(w, cookie)
			}
		}
		w.WriteHeader(status)
		writeJSONResponse(w, resp)
	}
}
