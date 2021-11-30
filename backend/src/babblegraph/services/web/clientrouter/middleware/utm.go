package middleware

import (
	"babblegraph/model/utm"
	"babblegraph/util/database"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func WithTrackingIDCapture(trackingEventName string, muxRouter func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		go func() {
			for _, cookie := range r.Cookies() {
				if cookie.Name == "uttrid" {
					trackingID := cookie.Value
					if err := database.WithTx(func(tx *sqlx.Tx) error {
						return utm.RegisterEvent(tx, trackingEventName, trackingID)
					}); err != nil {
						log.Println(fmt.Sprintf("Error registering event %s for tracking id %s", trackingEventName, trackingID))
					}
				}
			}
		}()
		muxRouter(w, r)
	}
}
