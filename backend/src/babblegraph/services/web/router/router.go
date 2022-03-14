package router

import (
	"fmt"

	"github.com/gorilla/mux"
)

func WithIndexRouter(router *mux.Router, routes []IndexRoute) error {
	paths := make(map[string]bool)
	for _, r := range routes {
		r := r
		if _, ok := paths[r.Path.Text]; ok {
			return fmt.Errorf("Duplicate paths found: %s", r.Path.Text)
		}
		if r.Path.UseAsPrefix {
			router.PathPrefix(r.Path.Text).HandlerFunc(r.makeMuxRoute()).Methods("GET")
		} else {
			router.HandleFunc(r.Path.Text, r.makeMuxRoute()).Methods("GET")
		}
	}
	return nil
}

func WithAPIRouter(r *mux.Router, apiPrefix string, routeGroups []RouteGroup) error {
	routeGroupPrefixes := make(map[string]bool)
	apiRouter := r.PathPrefix(fmt.Sprintf("/%s", apiPrefix)).Subrouter()
	for _, rg := range routeGroups {
		rg := rg
		if _, ok := routeGroupPrefixes[rg.Prefix]; ok {
			return fmt.Errorf("Duplicate route group prefix %s", rg.Prefix)
		}
		routeGroupPrefixes[rg.Prefix] = true
		routeNames := make(map[string]bool)
		for _, r := range rg.Routes {
			r := r
			if _, ok := routeNames[r.Path]; ok {
				return fmt.Errorf("Duplicate path %s for route group %s", r.Path, rg.Prefix)
			}
			routeNames[r.Path] = true
			apiRouter.HandleFunc(fmt.Sprintf("/%s/%s", rg.Prefix, r.Path), r.makeMuxRoute()).Methods("POST")
		}
	}
	return nil
}
