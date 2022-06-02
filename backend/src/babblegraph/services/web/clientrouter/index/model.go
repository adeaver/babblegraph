package index

import (
	"babblegraph/services/web/clientrouter/middleware"
	"babblegraph/util/env"
	"encoding/base64"
	"fmt"
	"io/ioutil"
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
	r.HandleFunc("/a/{token}", func(w http.ResponseWriter, r *http.Request) {
		routeVars := mux.Vars(r)
		token, ok := routeVars["token"]
		if !ok {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		url, err := base64.URLEncoding.DecodeString(token)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		resp, err := http.Get(string(url))
		if err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		w.Write(data)
	})
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
