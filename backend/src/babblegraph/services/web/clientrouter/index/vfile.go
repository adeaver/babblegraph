package index

import (
	"net/http"

	"github.com/gorilla/mux"
)

func HandleVirtualFile(w http.ResponseWriter, r *http.Request) {
	routeVars := mux.Vars(r)
	_, ok := routeVars["fileName"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "https://traffic.megaphone.fm/GLT9301967030.mp3?updated=1642621471", http.StatusFound)
}
