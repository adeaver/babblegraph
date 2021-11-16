package adminrouter

import (
	"babblegraph/services/web/adminrouter/api/user"
	"babblegraph/services/web/router"
	"babblegraph/util/env"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterAdminRouter(r *mux.Router) error {
	staticFileDirName := env.MustEnvironmentVariable("STATIC_DIR")
	s := r.PathPrefix("/ops").Subrouter()
	if err := router.WithAPIRouter(s, "api", []router.RouteGroup{
		user.Routes,
	}); err != nil {
		return err
	}
	s.PathPrefix("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/ops_index.html", staticFileDirName))
	})
	return nil
}
