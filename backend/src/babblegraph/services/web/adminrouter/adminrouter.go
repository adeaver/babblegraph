package adminrouter

import (
	"babblegraph/util/env"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterAdminRouter(r *mux.Router) error {
	staticFileDirName := env.MustEnvironmentVariable("STATIC_DIR")
	r.PathPrefix("/ops").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/ops_index.html", staticFileDirName))
	})
	return nil
}
