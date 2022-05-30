package tasks

import (
	"babblegraph/model/billing"
	"babblegraph/model/useraccounts"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func SubscriptionExpirationDryRun(c ctx.LogContext) error {
	var expiringSubscriptions []useraccounts.ExpiringSubscriptionInfo
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		expiringSubscriptions, err = useraccounts.GetExpiringSubscriptions(tx)
		return err
	}); err != nil {
		return fmt.Errorf("Error getting expiring subscriptions for dry run: %s", err.Error())

	}
	for _, sub := range expiringSubscriptions {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			switch {
			case sub.SubscriptionLevel == useraccounts.SubscriptionLevelLegacy:
				c.Infof("Will skip sub for user %s because it's in legacy", sub.UserID)
				return nil
			case sub.SubscriptionLevel == useraccounts.SubscriptionLevelPremium:
				c.Infof("Subscription for user %s has premium level", sub.UserID)
				premiumNewsletterSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(c, tx, sub.UserID)
				if err != nil {
					return err
				}
				switch {
				case !time.Now().Before(sub.ExpiringAt):
					c.Infof("Subscription for user %s is expired", sub.UserID)
					// We think the subscription should be expired.
					switch {
					case premiumNewsletterSubscription == nil:
						// Stripe also thinks that the subscription is expired, no-op
						c.Infof("Subscription for user %s requires no action", sub.UserID)
					case time.Now().Before(premiumNewsletterSubscription.CurrentPeriodEnd):
						// Stripe does not think that the subscription is over, we update
						c.Infof("Subscription for user %s requires an update", sub.UserID)
						return nil
					case premiumNewsletterSubscription.PaymentState == billing.PaymentStateActive,
						premiumNewsletterSubscription.PaymentState == billing.PaymentStateTrialPaymentMethodAdded:
						c.Infof("Subscription %s is in state %d, and has expired", *premiumNewsletterSubscription.ID, premiumNewsletterSubscription.PaymentState)
					}
				default:
					c.Infof("Subscription for user %s is not expired", sub.UserID)
					// We don't think we should be expired
					if premiumNewsletterSubscription == nil {
						// Stripe thinks we should be expired - at this point, we just wait it out.
						return nil
					}
					numberOfDaysUntilSubscriptionExpires := int64(sub.ExpiringAt.Sub(time.Now().Add(24*time.Hour*useraccounts.DefaultSubscriptionBufferInDays)) / (time.Duration(24) * time.Hour))
					c.Infof("Subscription for user %s will expire in %d days", sub.UserID, numberOfDaysUntilSubscriptionExpires)
					switch numberOfDaysUntilSubscriptionExpires {
					case 7:
						if premiumNewsletterSubscription.PaymentState == billing.PaymentStateTrialNoPaymentMethod {
							c.Infof("Subscription for user %s does not have payment method, will enqueue 7 day notification", sub.UserID)
						} else {
							c.Infof("Subscription for user %s does have payment method, will not enqueue 7 day notification", sub.UserID)
						}
					case 3:
						switch premiumNewsletterSubscription.PaymentState {
						case billing.PaymentStateTrialNoPaymentMethod,
							billing.PaymentStateTrialPaymentMethodAdded:
							c.Infof("Subscription for user %s will have a trial ending soon notification", sub.UserID)
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
							c.Infof("Subscription for user %s does not have payment method, will enqueue 2 day notification", sub.UserID)
						}
					case 1:
						if premiumNewsletterSubscription == nil {
							return fmt.Errorf("User %s should have a stripe subscription, but doesn't", sub.UserID)
						}
						c.Infof("Subscription for user %s does not have payment method, will enqueue 1 day notification", sub.UserID)
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
							c.Infof("Subscription for user %s will enqueue error notification", sub.UserID)
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
	return nil
}
