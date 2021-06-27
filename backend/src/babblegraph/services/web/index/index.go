package index

import (
	"babblegraph/model/utm"
	"babblegraph/services/web/middleware"
	"babblegraph/util/database"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func HandleServeIndexPage(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.LogRequestWithoutBody(r)
		utmParameters := utm.GetParametersForRequest(r)
		if utmParameters != nil {
			trackingID, err := utm.GetTrackingIDForRequest(r)
			switch {
			case err != nil:
				log.Println(fmt.Sprintf("Error on getting tracking ID for request: %s", err.Error()))
			case trackingID != nil:
				go func() {
					if err := database.WithTx(func(tx *sqlx.Tx) error {
						return utm.RegisterUTMPageHit(tx, *trackingID, *utmParameters)
					}); err != nil {
						log.Println(fmt.Sprintf("Error registering UTM page hit: %s", err.Error()))
					}
				}()
				http.SetCookie(w, &http.Cookie{
					Name:  utm.UTMTrackingIDCookieName,
					Value: trackingID.Str(),
				})
			default:
				log.Println("Unknown error making tracking ID for request, continuing...")
			}
		}
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", staticFileDirName))
	}
}
