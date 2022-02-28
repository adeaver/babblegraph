package useraccounts

import (
	"babblegraph/model/billing"
	"babblegraph/model/users"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func SyncUserAccountWithPremiumNewsletterSubscription(tx *sqlx.Tx, userID users.UserID, premiumNewsletterSubscription *billing.PremiumNewsletterSubscription) error {
	if premiumNewsletterSubscription == nil {
		return ExpireSubscriptionForUser(tx, userID)
	}
	switch premiumNewsletterSubscription.PaymentState {
	case billing.PaymentStateCreatedUnpaid:
		return nil
	case billing.PaymentStateTrialNoPaymentMethod,
		billing.PaymentStateTrialPaymentMethodAdded,
		billing.PaymentStateActive:
		subscriptionLevel, err := LookupSubscriptionLevelForUser(tx, userID)
		switch {
		case err != nil:
			return err
		case subscriptionLevel == nil:
			return AddSubscriptionLevelForUser(tx, AddSubscriptionLevelForUserInput{
				UserID:            userID,
				SubscriptionLevel: SubscriptionLevelPremium,
				ShouldStartActive: true,
				ExpirationTime:    premiumNewsletterSubscription.CurrentPeriodEnd,
			})
		case subscriptionLevel != nil:
			return UpdateSubscriptionExpirationTime(tx, userID, premiumNewsletterSubscription.CurrentPeriodEnd)
		}
	case billing.PaymentStateErrored,
		billing.PaymentStateTerminated:
		return ExpireSubscriptionForUser(tx, userID)
	default:
		return fmt.Errorf("Unsupported Payment State %d", premiumNewsletterSubscription.PaymentState)
	}
	panic("Unreachable")
}
