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
	c.Infof("Starting billing sync")
	var expiringSubscriptions []useraccounts.ExpiringSubscriptionInfo
	var premiumNewsletterSyncRequests map[billing.PremiumNewsletterSubscriptionID]billing.PremiumNewsletterSubscriptionUpdateType
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		expiringSubscriptions, err = useraccounts.GetExpiringSubscriptions(tx)
		if err != nil {
			return err
		}
		premiumNewsletterSyncRequests, err = billing.GetPremiumNewsletterSyncRequests(tx)
		return err
	}); err != nil {
		c.Errorf("Error getting billing items to sync: %s", err.Error())
		return
	}
	c.Infof("Billing sync got %d sync requests and %d expiring subscriptions", len(premiumNewsletterSyncRequests), len(expiringSubscriptions))
	userIDsWithExpirationToSkip := make(map[users.UserID]bool)
	for premiumSubscriptionID, updateType := range premiumNewsletterSyncRequests {
		var userID *users.UserID
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			premiumNewsletterSubscription, err := billing.GetPremiumNewsletterSubscriptionByID(c, tx, premiumSubscriptionID)
			if err != nil {
				return err
			}
			userID, err = premiumNewsletterSubscription.GetUserID()
			if err != nil {
				return err
			}
			switch updateType {
			case billing.PremiumNewsletterSubscriptionUpdateTypeTransitionToActive:
				switch premiumNewsletterSubscription.PaymentState {
				case billing.PaymentStateCreatedUnpaid,
					billing.PaymentStatePaymentPending:
					c.Infof("Subscription %s is still unpaid, skipping.", premiumSubscriptionID)
					return billing.MarkPremiumNewsletterSyncRequestForRetry(tx, premiumSubscriptionID)
				case billing.PaymentStateTrialNoPaymentMethod,
					billing.PaymentStateTrialPaymentMethodAdded,
					billing.PaymentStateActive,
					billing.PaymentStateErrored,
					billing.PaymentStateTerminated:
					if err := billing.SyncUserAccountWithPremiumNewsletterSubscription(tx, *userID, premiumNewsletterSubscription); err != nil {
						return err
					}
					c.Infof("Subscription %s is can be synced, marking done", premiumSubscriptionID)
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
					billing.PaymentStatePaymentPending,
					billing.PaymentStateErrored:
					return fmt.Errorf("Subscription %s is in the wrong state", premiumSubscriptionID)
				case billing.PaymentStateTerminated:
					if err := billing.SyncUserAccountWithPremiumNewsletterSubscription(tx, *userID, premiumNewsletterSubscription); err != nil {
						return err
					}
					return billing.MarkPremiumNewsletterSyncRequestDone(tx, premiumSubscriptionID)
				default:
					return fmt.Errorf("Unrecognized payment state for subscription ID %s: %d", premiumSubscriptionID, premiumNewsletterSubscription.PaymentState)
				}
			case billing.PremiumNewsletterSubscriptionUpdateTypeRemoteUpdated,
				billing.PremiumNewsletterSubscriptionUpdateTypePaymentMethodAdded:
				if err := billing.SyncUserAccountWithPremiumNewsletterSubscription(tx, *userID, premiumNewsletterSubscription); err != nil {
					return err
				}
				return billing.MarkPremiumNewsletterSyncRequestDone(tx, premiumSubscriptionID)
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
	for _, sub := range expiringSubscriptions {
		if _, ok := userIDsWithExpirationToSkip[sub.UserID]; ok {
			c.Infof("User ID %s had an update, skipping", sub.UserID)
			continue
		}
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			switch {
			case sub.SubscriptionLevel == useraccounts.SubscriptionLevelLegacy:
				c.Infof("Skipping subscription because it's in legacy")
				return nil
			case sub.SubscriptionLevel == useraccounts.SubscriptionLevelPremium:
				premiumNewsletterSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(c, tx, sub.UserID)
				if err != nil {
					return err
				}
				switch {
				case !time.Now().Before(sub.ExpiringAt):
					// We think the subscription should be expired.
					switch {
					case premiumNewsletterSubscription == nil:
						// Stripe also thinks that the subscription is expired, no-op
					case time.Now().Before(premiumNewsletterSubscription.CurrentPeriodEnd):
						// Stripe does not think that the subscription is over
						if premiumNewsletterSubscription.PaymentState != billing.PaymentStateErrored {
							// If the state is not errored, we need to update
							c.Infof("User ID has an active subscription. Updating...")
							return useraccounts.UpdateSubscriptionExpirationTime(tx, sub.UserID, premiumNewsletterSubscription.CurrentPeriodEnd)
						}
						// If the state is errored, we fall through here
					case premiumNewsletterSubscription.PaymentState == billing.PaymentStateActive,
						premiumNewsletterSubscription.PaymentState == billing.PaymentStateTrialPaymentMethodAdded:
						c.Warnf("Subscription %s is in state %d, and has expired", *premiumNewsletterSubscription.ID, premiumNewsletterSubscription.PaymentState)
					}
					if err := useraccounts.ExpireSubscriptionForUser(tx, sub.UserID); err != nil {
						return err
					}
					if _, err := useraccountsnotifications.EnqueueNotificationRequest(tx, sub.UserID, useraccountsnotifications.NotificationTypePremiumSubscriptionCanceled, time.Now()); err != nil {
						return err
					}
					return billing.CancelPremiumNewsletterSubscriptionForUser(c, tx, sub.UserID)
				default:
					// We don't think we should be expired
					if premiumNewsletterSubscription == nil {
						// Stripe thinks we should be expired - at this point, we just wait it out.
						return nil
					}
					numberOfDaysUntilSubscriptionExpires := int64(sub.ExpiringAt.Sub(time.Now().Add(24*time.Hour*useraccounts.DefaultSubscriptionBufferInDays)) / (time.Duration(24) * time.Hour))
					switch numberOfDaysUntilSubscriptionExpires {
					case 7:
						if premiumNewsletterSubscription.PaymentState == billing.PaymentStateTrialNoPaymentMethod {
							_, err := useraccountsnotifications.EnqueueNotificationRequest(tx, sub.UserID, useraccountsnotifications.NotificationTypeNeedPaymentMethodWarning, time.Now())
							return err
						}
					case 3:
						switch premiumNewsletterSubscription.PaymentState {
						case billing.PaymentStateTrialNoPaymentMethod,
							billing.PaymentStateTrialPaymentMethodAdded:
							_, err := useraccountsnotifications.EnqueueNotificationRequest(tx, sub.UserID, useraccountsnotifications.NotificationTypeTrialEndingSoon, time.Now())
							return err
						case billing.PaymentStateActive,
							billing.PaymentStateTerminated,
							billing.PaymentStateCreatedUnpaid,
							billing.PaymentStatePaymentPending,
							billing.PaymentStateErrored:
							// no-op
						default:
							c.Warnf("Unrecognized payment state for subscription %s: %d", *premiumNewsletterSubscription.ID, premiumNewsletterSubscription.PaymentState)
						}
					case 2:
						if premiumNewsletterSubscription == nil {
							return fmt.Errorf("User %s should have a stripe subscription, but doesn't", sub.UserID)
						}
						if premiumNewsletterSubscription.PaymentState == billing.PaymentStateTrialNoPaymentMethod {
							_, err := useraccountsnotifications.EnqueueNotificationRequest(tx, sub.UserID, useraccountsnotifications.NotificationTypeNeedPaymentMethodWarningUrgent, time.Now())
							return err
						}
					case 1:
						if premiumNewsletterSubscription == nil {
							return fmt.Errorf("User %s should have a stripe subscription, but doesn't", sub.UserID)
						}
						if premiumNewsletterSubscription.PaymentState == billing.PaymentStateTrialNoPaymentMethod {
							_, err := useraccountsnotifications.EnqueueNotificationRequest(tx, sub.UserID, useraccountsnotifications.NotificationTypeNeedPaymentMethodWarningVeryUrgent, time.Now())
							return err
						}
					case -1:
						if premiumNewsletterSubscription == nil {
							return nil
						}
						// This happens when we've entered the buffer phase
						switch premiumNewsletterSubscription.PaymentState {
						case billing.PaymentStateTrialNoPaymentMethod,
							billing.PaymentStateTrialPaymentMethodAdded,
							billing.PaymentStateActive,
							billing.PaymentStateTerminated,
							billing.PaymentStatePaymentPending,
							billing.PaymentStateCreatedUnpaid:
							// no-op
						case billing.PaymentStateErrored:
							_, err := useraccountsnotifications.EnqueueNotificationRequest(tx, sub.UserID, useraccountsnotifications.NotificationTypePaymentError, time.Now())
							return err
						default:
							c.Warnf("Unrecognized payment state for subscription %s: %d", *premiumNewsletterSubscription.ID, premiumNewsletterSubscription.PaymentState)
						}
					}
				}
			}
			return nil
		}); err != nil {
			c.Errorf("Error canceling subscription for User %s: %s", sub.UserID, err.Error())
		}
	}
}
