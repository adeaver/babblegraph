package index

import (
	"babblegraph/services/web/clientrouter/middleware"
	"babblegraph/util/env"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterIndexRoutes(r *mux.Router, indexPages []IndexPage) error {
	staticFileDirName := env.MustEnvironmentVariable("STATIC_DIR")
	registeredIndexRoutes := make(map[string]bool)
	for _, page := range indexPages {
		if _, ok := registeredIndexRoutes[page.RouteName]; ok {
			return fmt.Errorf("Route %s is already registered", page.RouteName)
		}
		registeredIndexRoutes[page.RouteName] = true
		r.HandleFunc(page.RouteName, func(w http.ResponseWriter, r *http.Request) {
			if page.HandleAuthorization == nil {
				// TODO: replace this with the title
				HandleServeIndexPage(staticFileDirName)(w, r)
				return
			}
			middleware.WithAuthorizationCheck(w, r, *page.HandleAuthorization)
		})
	}
	r.HandleFunc("/logout", HandleLogout())
	r.HandleFunc("/login", HandleLoginPage(staticFileDirName))
	r.HandleFunc("/signup/{token}", HandleCreateUserPage(staticFileDirName))
	r.HandleFunc("/blog/{blog_path}", HandleServeBlogPost(staticFileDirName))
	r.HandleFunc("/dist/{token}/logo.png", HandleServeLogo(staticFileDirName))
	r.PathPrefix("/dist").Handler(http.StripPrefix("/dist", http.FileServer(http.Dir(staticFileDirName))))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/favicon.ico", staticFileDirName))
	})
	r.HandleFunc("/article/{token}", handleArticleRoute(staticFileDirName))
	r.HandleFunc("/a/{token}", handleArticleHTMLPassthrough())
	r.HandleFunc("/out/{token}", handleArticleOutLink())
	r.PathPrefix("/").HandlerFunc(HandleServeIndexPage(staticFileDirName))
	return nil
}

// An Index Page is any frontend
// page that requires some backend calculation
// This can be something like page title, authentication,
// etc.

type IndexPage struct {
	RouteName           string
	RouteTitle          *string
	HandleAuthorization *middleware.WithAuthorizationCheckInput
}
