package index

import (
	"babblegraph/model/billing"
	"babblegraph/model/routes"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func handleVerification(r *router.Request) (interface{}, error) {
	token, err := r.GetRouteVar("token")
	if err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(*token, routes.UserVerificationKey)
	if err != nil {
		return nil, err
	}
	var premiumNewsletterSubscription *billing.PremiumNewsletterSubscription
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		if err := users.SetUserStatusToVerified(tx, *userID); err != nil {
			return err
		}
		if _, err := billing.GetOrCreateBillingInformationForUser(r, tx, *userID); err != nil {
			return err
		}
		premiumNewsletterSubscriptionID := billing.NewPremiumNewsletterSubscriptionID()
		premiumNewsletterSubscription, err = billing.LookupPremiumNewsletterSubscriptionForUser(r, tx, *userID)
		switch {
		case err != nil:
			return err
		case premiumNewsletterSubscription != nil:
			return nil
		}
		if err := billing.InsertPremiumNewsletterSyncRequest(tx, premiumNewsletterSubscriptionID, billing.PremiumNewsletterSubscriptionUpdateTypeTransitionToActive); err != nil {
			return err
		}
		premiumNewsletterSubscription, err = billing.CreatePremiumNewsletterSubscriptionForUserWithID(r, tx, *userID, premiumNewsletterSubscriptionID)
		return err
	}); err != nil {
		return nil, err
	}
	switch premiumNewsletterSubscription.PaymentState {
	case billing.PaymentStateCreatedUnpaid:
		return routes.MakePremiumSubscriptionCheckoutLink(*userID)
	case billing.PaymentStateTrialNoPaymentMethod,
		billing.PaymentStateTrialPaymentMethodAdded,
		billing.PaymentStateActive:
		return routes.MakeSubscriptionManagementRouteForUserID(*userID)
	case billing.PaymentStateErrored,
		billing.PaymentStateTerminated:
		return nil, fmt.Errorf("Got invalid state: %d", premiumNewsletterSubscription.PaymentState)
	default:
		return nil, fmt.Errorf("Unrecognized state: %d", premiumNewsletterSubscription.PaymentState)
	}
}
