package main

import (
	"babblegraph/services/web/api/user"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting babblegraph web server")
	r := mux.NewRouter()
	staticFileDirName := env.MustEnvironmentVariable("STATIC_DIR")

	if err := setupDatabases(); err != nil {
		log.Fatal(err.Error())
	}

	if err := registerAPI(r); err != nil {
		log.Fatal(err.Error())
	}

	r.PathPrefix("/dist").Handler(http.StripPrefix("/dist", http.FileServer(http.Dir(staticFileDirName))))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", staticFileDirName))
	})

	http.ListenAndServe(":8080", r)
}

func setupDatabases() error {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		return fmt.Errorf("Error setting up main-db: %s", err.Error())
	}
	return nil
}

func registerAPI(r *mux.Router) error {
	router.CreateNewAPIRouter(r)
	if err := user.RegisterRouteGroups(); err != nil {
		return err
	}
	return nil
}
