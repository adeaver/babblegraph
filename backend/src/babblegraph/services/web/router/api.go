package router

import (
	"fmt"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
)

const apiPrefix string = "api"

type apiRouter struct {
	r             *mux.Router
	prefixes      map[string]bool
	routeNames    map[string]bool
	sentryHandler *sentryhttp.Handler
}

var a *apiRouter = nil

func CreateNewAPIRouter(mainRouter *mux.Router) {
	a = &apiRouter{
		r:          mainRouter,
		prefixes:   make(map[string]bool),
		routeNames: make(map[string]bool),
		sentryHandler: sentryhttp.New(sentryhttp.Options{
			Repanic: true,
		}),
	}
	registerUserAccountsRoutes()
}

type RouteGroup struct {
	Prefix string
	Routes []Route
}

func RegisterRouteGroup(rg RouteGroup) error {
	if a == nil {
		panic("API Router not setup")
	}
	_, ok := a.prefixes[rg.Prefix]
	if ok {
		return fmt.Errorf("Duplicate route groups with prefix %s", rg.Prefix)
	}
	a.prefixes[rg.Prefix] = true
	for _, r := range rg.Routes {
		path := fmt.Sprintf("/%s/%s/%s", apiPrefix, rg.Prefix, r.Path)
		if _, ok := a.routeNames[path]; ok {
			return fmt.Errorf("Duplicate paths %s", path)
		}
		muxRoute := makeMuxRouter(r.Handler)
		if r.ShouldLogBody {
			muxRoute = withBodyLogger(muxRoute)
		} else {
			muxRoute = withoutBodyLogger(muxRoute)
		}
		if r.TrackEventWithID != nil {
			muxRoute = withTrackingIDCapture(*r.TrackEventWithID, muxRoute)
		}
		a.r.HandleFunc(path, a.sentryHandler.HandleFunc(muxRoute)).Methods("POST")
		a.routeNames[path] = true
	}
	return nil
}
