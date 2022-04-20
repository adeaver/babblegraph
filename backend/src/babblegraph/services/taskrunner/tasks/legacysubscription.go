package tasks

import (
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"time"

	"github.com/jmoiron/sqlx"
)

func MigrateLegacyUsers(c ctx.LogContext) error {
	defaultExpirationTime := time.Date(2022, 12, 1, 6, 0, 0, 0, time.UTC)
	return database.WithTx(func(tx *sqlx.Tx) error {
		u, err := users.GetAllActiveUsers(tx)
		if err != nil {
			return err
		}
		for _, user := range u {
			subscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, user.ID)
			switch {
			case err != nil:
				return err
			case subscriptionLevel != nil:
				// no-op
			default:
				if err := useraccounts.AddSubscriptionLevelForUser(tx, useraccounts.AddSubscriptionLevelForUserInput{
					UserID:            user.ID,
					SubscriptionLevel: useraccounts.SubscriptionLevelLegacy,
					ShouldStartActive: true,
					ExpirationTime:    defaultExpirationTime,
				}); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
