package index

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/services/web/clientrouter/middleware"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func HandleServeLogo(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.LogRequestWithoutBody(r)
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
				return encrypt.WithDecodedToken(token, func(t encrypt.TokenPair) error {
					if t.Key != routes.EmailOpenedKey.Str() {
						return fmt.Errorf("Token has wrong key: %s", t.Key)
					}
					emailRecordID, ok := t.Value.(string)
					if !ok {
						return fmt.Errorf("Token has wrong value type")
					}
					return email.SetEmailFirstOpened(tx, email.ID(emailRecordID))
				})
			}); err != nil {
				log.Println(fmt.Sprintf("Got error handling token %s: %s", token, err.Error()))
			}
		}()
		http.ServeFile(w, r, fmt.Sprintf("%s/logo.png", staticFileDirName))
	}
}
