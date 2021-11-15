package adminrouter

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterAdminRouter(r *mux.Router) error {
	r.HandleFunc("/ops", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>babblegraph.com/ops is under construction\n</h1>")
	})
	return nil
}
