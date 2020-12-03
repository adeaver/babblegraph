package main

import (
	"babblegraph/util/env"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	staticFileDirName := env.MustEnvironmentVariable("STATIC_DIR")

	// TODO: put API router in
	// apiRouter := r.PathPrefix("/api")
	r.PathPrefix("/dist").Handler(http.FileServer(http.Dir(staticFileDirName)))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", staticFileDirName))
	})

	http.ListenAndServe(":8080", r)
}
