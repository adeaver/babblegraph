package adminrouter

import (
	"babblegraph/model/admin"
	"babblegraph/services/web/adminrouter/api/auth"
	"babblegraph/services/web/adminrouter/api/billing"
	"babblegraph/services/web/adminrouter/api/blog"
	"babblegraph/services/web/adminrouter/api/content"
	"babblegraph/services/web/adminrouter/api/podcasts"
	"babblegraph/services/web/adminrouter/api/usermetrics"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func RegisterAdminRouter(r *mux.Router) error {
	staticFileDirName := env.MustEnvironmentVariable("STATIC_DIR")
	s := r.PathPrefix("/ops").Subrouter()
	if err := router.WithAPIRouter(s, "api", []router.RouteGroup{
		auth.Routes,
		usermetrics.Routes,
		blog.Routes,
		content.Routes,
		billing.Routes,
		podcasts.Routes,
	}); err != nil {
		return err
	}
	s.PathPrefix("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL == nil || !(r.URL.Path == "/ops" || r.URL.Path == "/ops/" || strings.HasPrefix(r.URL.Path, "/ops/register")) {
			var hasAuth bool
			for _, cookie := range r.Cookies() {
				if cookie.Name == admin.AccessTokenCookieName {
					token := cookie.Value
					if err := database.WithTx(func(tx *sqlx.Tx) error {
						if _, err := admin.ValidateAccessTokenAndGetUserID(tx, token); err != nil {
							return err
						}
						return nil
					}); err == nil {
						hasAuth = true
					}
				}
			}
			if !hasAuth {
				http.Redirect(w, r, "/ops", http.StatusTemporaryRedirect)
				return
			}
		}
		http.ServeFile(w, r, fmt.Sprintf("%s/ops_index.html", staticFileDirName))
		return
	})
	return nil
}
