package routermiddleware

import (
	"babblegraph/model/utm"
	"babblegraph/services/web/router"
	"babblegraph/util/async"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

const (
	utmTrackingIDCookieName = "uttrid"
)

func WithUTMEventTracking(trackingEvent string, handler router.RequestHandler) router.RequestHandler {
	return func(r *router.Request) (interface{}, error) {
		async.WithContext(make(chan error, 1), "utm-tracking", func(c async.Context) {
			for _, cookie := range r.GetCookies() {
				if cookie.Name == "uttrid" {
					trackingID := cookie.Value
					if err := database.WithTx(func(tx *sqlx.Tx) error {
						return utm.RegisterEvent(tx, trackingEvent, trackingID)
					}); err != nil {
						c.Infof("Error registering event %s for tracking id %s", trackingEvent, trackingID)
					}
				}
			}
		}).Start()
		return handler(r)
	}
}
