package clientrouter

import (
	"babblegraph/services/web/clientrouter/api"
	"babblegraph/services/web/clientrouter/api/billing"
	"babblegraph/services/web/clientrouter/api/blog"
	"babblegraph/services/web/clientrouter/api/language"
	"babblegraph/services/web/clientrouter/api/podcasts"
	"babblegraph/services/web/clientrouter/api/ses"
	"babblegraph/services/web/clientrouter/api/token"
	"babblegraph/services/web/clientrouter/api/user"
	"babblegraph/services/web/clientrouter/api/useraccounts"
	"babblegraph/services/web/clientrouter/api/utm"
	"babblegraph/services/web/clientrouter/index"
	"babblegraph/services/web/router"

	"github.com/gorilla/mux"
)

func RegisterClientRouter(r *mux.Router) error {
	if err := registerAPI(r); err != nil {
		return err
	}
	if err := router.WithAPIRouter(r, "api", []router.RouteGroup{
		useraccounts.Routes,
		billing.Routes,
		user.Routes,
		podcasts.Routes,
	}); err != nil {
		return err
	}
	return index.RegisterIndexRoutes(r, []index.IndexPage{})
}

func registerAPI(r *mux.Router) error {
	api.CreateNewAPIRouter(r)
	if err := user.RegisterRouteGroups(); err != nil {
		return err
	}
	if err := ses.RegisterRouteGroups(); err != nil {
		return err
	}
	if err := utm.RegisterRouteGroups(); err != nil {
		return err
	}
	if err := language.RegisterRouteGroups(); err != nil {
		return err
	}
	if err := token.RegisterRouteGroups(); err != nil {
		return err
	}
	if err := blog.RegisterRouteGroups(); err != nil {
		return err
	}
	return nil
}
