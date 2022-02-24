package scheduler

import (
	"babblegraph/model/billing"
	"babblegraph/model/useraccounts"
	"babblegraph/model/useraccountsnotifications"
	"babblegraph/model/users"
	"babblegraph/util/async"
	"babblegraph/util/database"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func handleSyncBilling(c async.Context) {
	var userIDsWithExpiredSubscription []users.UserID
	var premiumNewsletterSyncRequests map[billing.PremiumNewsletterSubscriptionID]billing.PremiumNewsletterSubscriptionUpdateType
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userIDsWithExpiredSubscription, err = useraccounts.GetUserIDsForExpiredSubscriptionQuery(tx)
		if err != nil {
			return err
		}
		premiumNewsletterSyncRequests, err = billing.GetPremiumNewsletterSyncRequests(tx)
		return err
	}); err != nil {
		c.Errorf("Error getting billing items to sync: %s", err.Error())
		return
	}
	userIDsWithExpirationToSkip := make(map[users.UserID]bool)
	for premiumSubscriptionID, updateType := range premiumNewsletterSyncRequests {
		var userID *users.UserID
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			premiumNewsletterSubscription, err := billing.GetPremiumNewsletterSubscriptionByID(c, tx, premiumSubscriptionID)
			if err != nil {
				return err
			}
			switch updateType {
			case billing.PremiumNewsletterSubscriptionUpdateTypeTransitionToActive:
				switch premiumNewsletterSubscription.PaymentState {
				case billing.PaymentStateCreatedUnpaid:
					c.Infof("Subscription %s is still unpaid, skipping.", premiumSubscriptionID)
					return billing.MarkPremiumNewsletterSyncRequestForRetry(tx, premiumSubscriptionID)
				case billing.PaymentStateTrialNoPaymentMethod,
					billing.PaymentStateTrialPaymentMethodAdded,
					billing.PaymentStateActive:
					c.Infof("Subscription %s is now active, updating", premiumSubscriptionID)
					userID, err = premiumNewsletterSubscription.GetUserID()
					if err != nil {
						return err
					}
					if err := useraccounts.AddSubscriptionLevelForUser(tx, useraccounts.AddSubscriptionLevelForUserInput{
						UserID:            *userID,
						SubscriptionLevel: useraccounts.SubscriptionLevelPremium,
						ShouldStartActive: true,
						ExpirationTime:    premiumNewsletterSubscription.CurrentPeriodEnd,
					}); err != nil {
						return err
					}
					return billing.MarkPremiumNewsletterSyncRequestDone(tx, premiumSubscriptionID)
				case billing.PaymentStateErrored,
					billing.PaymentStateTerminated:
					c.Infof("Subscription %s is terminated, marking done", premiumSubscriptionID)
					return billing.MarkPremiumNewsletterSyncRequestDone(tx, premiumSubscriptionID)
				default:
					return fmt.Errorf("Unrecognized payment state for subscription ID %s: %d", premiumSubscriptionID, premiumNewsletterSubscription.PaymentState)
				}
			case billing.PremiumNewsletterSubscriptionUpdateTypeCanceled:
				switch premiumNewsletterSubscription.PaymentState {
				case billing.PaymentStateCreatedUnpaid,
					billing.PaymentStateTrialNoPaymentMethod,
					billing.PaymentStateTrialPaymentMethodAdded,
					billing.PaymentStateActive,
					billing.PaymentStateErrored:
					return fmt.Errorf("Subscription %s is in the wrong state", premiumSubscriptionID)
				case billing.PaymentStateTerminated:
					userID, err = premiumNewsletterSubscription.GetUserID()
					if err != nil {
						return err
					}
					if err := useraccounts.ExpireSubscriptionForUser(tx, *userID); err != nil {
						return err
					}
					if _, err := useraccountsnotifications.EnqueueNotificationRequest(tx, *userID, useraccountsnotifications.NotificationTypePremiumSubscriptionCanceled, time.Now().Add(5*time.Minute)); err != nil {
						return err
					}
					return billing.MarkPremiumNewsletterSyncRequestDone(tx, premiumSubscriptionID)
				default:
					return fmt.Errorf("Unrecognized payment state for subscription ID %s: %d", premiumSubscriptionID, premiumNewsletterSubscription.PaymentState)
				}
			case billing.PremiumNewsletterSubscriptionUpdateTypeRemoteUpdated:
				userID, err = premiumNewsletterSubscription.GetUserID()
				if err != nil {
					return err
				}
				return useraccounts.SyncUserAccountWithPremiumNewsletterSubscription(tx, *userID, premiumNewsletterSubscription)
			default:
				return fmt.Errorf("Unrecognized update type: %s", updateType)
			}
		}); err != nil {
			c.Errorf("Error syncing stripe subscription with ID %s: %s", premiumSubscriptionID, err.Error())
			continue
		}
		if userID != nil {
			userIDsWithExpirationToSkip[*userID] = true
		}
	}
	for _, userID := range userIDsWithExpiredSubscription {
		if _, ok := userIDsWithExpirationToSkip[userID]; ok {
			c.Infof("User ID %s had an update, skipping", userID)
			continue
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			premiumNewsletterSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(c, tx, userID)
			switch {
			case err != nil:
				return err
			case premiumNewsletterSubscription == nil:
				// no-op
			case time.Now().Before(premiumNewsletterSubscription.CurrentPeriodEnd):
				c.Infof("User ID has an active subscription. Updating...")
				return useraccounts.UpdateSubscriptionExpirationTime(tx, userID, premiumNewsletterSubscription.CurrentPeriodEnd)
			default:
				// no-op
			}
			if err := useraccounts.ExpireSubscriptionForUser(tx, userID); err != nil {
				return err
			}
			if _, err := useraccountsnotifications.EnqueueNotificationRequest(tx, userID, useraccountsnotifications.NotificationTypePremiumSubscriptionCanceled, time.Now().Add(5*time.Minute)); err != nil {
				return err
			}
			return billing.CancelPremiumNewsletterSubscriptionForUser(c, tx, userID)
		}); err != nil {
			c.Errorf("Error canceling subscription for User %s: %s", userID, err.Error())
		}
	}
}
