package index

import (
	"babblegraph/actions/email"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func HandleServeLogo(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		router.LogRequestWithoutBody(r)
		// In order to collect information about whether an email was opened, we pass
		// the logo hero image with a URL of the above format. This is done because
		// 1) Some clients ban zero width images
		// 2) I imagine some clients ban url parameters
		routeVars := mux.Vars(r)
		go func() {
			// This is done in a go routine because we do not want it to
			// affect the speed with which requests are handled.
			token, ok := routeVars["token"]
			if !ok {
				return
			}
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				return email.HandleDailyEmailOpenToken(tx, token)
			}); err != nil {
				log.Println(fmt.Sprintf("Got error handling token %s: %s", token, err.Error()))
			}
		}()
		http.ServeFile(w, r, fmt.Sprintf("%s/logo.png", staticFileDirName))
	}
}
