package scheduler

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/async"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

func expireUserAccounts(c async.Context) {
	var userIDs []users.UserID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userIDs, err = useraccounts.GetUserIDsForExpiredSubscriptionQuery(tx)
		return err
	}); err != nil {
		c.Errorf("Error getting IDs to expire: %s", err.Error())
		return
	}
	var numExpiredSubscription int64
	for _, userID := range userIDs {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			if err := useraccounts.ExpireSubscriptionForUser(tx, userID); err != nil {
				return err
			}
			return bgstripe.CancelSubscription(tx, userID)
		}); err != nil {
			c.Errorf("Error expiring subscription for user %s: %s", userID, err.Error())
			continue
		}
		numExpiredSubscription++
	}
	c.Infof("Expired %d subscriptions", numExpiredSubscription)
	return
}
