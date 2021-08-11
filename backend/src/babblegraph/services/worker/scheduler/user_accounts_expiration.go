package scheduler

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func expireUserAccounts() error {
	var userIDs []users.UserID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userIDs, err = useraccounts.GetUserIDsForExpiredSubscriptionQuery(tx)
		return err
	}); err != nil {
		return err
	}
	var numExpiredSubscription int64
	for _, userID := range userIDs {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			if err := useraccounts.ExpireSubscriptionForUser(tx, userID); err != nil {
				return err
			}
			return bgstripe.CancelSubscription(tx, userID)
		}); err != nil {
			sentry.CaptureException(fmt.Errorf("Error expiring subscription for user %s: %s", userID, err.Error()))
			continue
		}
		numExpiredSubscription++
	}
	log.Println(fmt.Sprintf("Expired %d subscriptions", numExpiredSubscription))
	return nil
}
