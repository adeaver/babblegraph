package router

import (
	"babblegraph/model/users"
	"babblegraph/services/web/middleware"
	"fmt"
	"net/http"

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
	Prefix              string
	Routes              []Route
	AuthenticatedRoutes []AuthenticatedRoute
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
		if err := registerRoute(registerRouteInput{
			shouldLogBody:    r.ShouldLogBody,
			trackEventWithID: r.TrackEventWithID,
			muxRoute:         makeMuxRouter(r.Handler),
			path:             fmt.Sprintf("/%s/%s/%s", apiPrefix, rg.Prefix, r.Path),
		}); err != nil {
			return err
		}
	}
	for _, ar := range rg.AuthenticatedRoutes {
		// Shadow ar so that the closure below
		// returns the right thing
		ar := ar
		if err := registerRoute(registerRouteInput{
			shouldLogBody:    ar.ShouldLogBody,
			trackEventWithID: ar.TrackEventWithID,
			muxRoute: middleware.WithAuthorizationLevelVerification(ar.ValidAuthorizationLevels, func(userID users.UserID) func(http.ResponseWriter, *http.Request) {
				return makeAuthenticatedMuxRouter(userID, ar.Handler)
			}),
			path: fmt.Sprintf("/%s/%s/%s", apiPrefix, rg.Prefix, ar.Path),
		}); err != nil {
			return err
		}
	}
	return nil
}

type registerRouteInput struct {
	shouldLogBody    bool
	trackEventWithID *string
	path             string
	muxRoute         func(http.ResponseWriter, *http.Request)
}

func registerRoute(input registerRouteInput) error {
	if _, ok := a.routeNames[input.path]; ok {
		return fmt.Errorf("Duplicate paths %s", input.path)
	}
	muxRoute := input.muxRoute
	if input.shouldLogBody {
		muxRoute = middleware.WithBodyLogger(muxRoute)
	} else {
		muxRoute = middleware.WithoutBodyLogger(muxRoute)
	}
	if input.trackEventWithID != nil {
		muxRoute = middleware.WithTrackingIDCapture(*input.trackEventWithID, muxRoute)
	}
	a.r.HandleFunc(input.path, a.sentryHandler.HandleFunc(muxRoute)).Methods("POST")
	a.routeNames[input.path] = true
	return nil
}
