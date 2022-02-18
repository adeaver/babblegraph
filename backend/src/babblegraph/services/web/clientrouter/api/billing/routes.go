package billing

import (
	"babblegraph/model/billing"
	"babblegraph/model/routes"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"net/http"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "billing",
	Routes: []router.Route{
		{
			Path: "get_or_create_billing_information_1",
			Handler: routermiddleware.WithRequestBodyLogger(
				routermiddleware.WithAuthentication(getOrCreateBillingInformation),
			),
		}, {
			Path: "get_or_create_premium_newsletter_subscription_1",
			Handler: routermiddleware.WithRequestBodyLogger(
				routermiddleware.WithAuthentication(getOrCreatePremiumNewsletterSubscription),
			),
		}, {
			Path: "stripe_begin_payment_method_setup_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.WithAuthentication(stripeBeginPaymentMethodSetup),
			),
		},
	},
}

type getOrCreateBillingInformationRequest struct {
	PremiumSubscriptionCheckoutToken string `json:"premium_subscription_checkout_token"`
}

type getOrCreateBillingInformationResponse struct {
	StripeCustomerID string `json:"stripe_customer_id,omitempty"`
}

func getOrCreateBillingInformation(userAuth routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req getOrCreateBillingInformationRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.PremiumSubscriptionCheckoutToken, routes.PremiumSubscriptionCheckoutKey)
	if err != nil || *userID != userAuth.UserID {
		r.RespondWithStatus(http.StatusForbidden)
		return nil, nil
	}
	var billingInformation *billing.BillingInformation
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		billingInformation, err = billing.GetOrCreateBillingInformationForUser(r, tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	return getOrCreateBillingInformationResponse{
		StripeCustomerID: *billingInformation.StripeCustomerID,
	}, nil
}

type getOrCreatePremiumNewsletterSubscriptionRequest struct {
	PremiumSubscriptionCheckoutToken string `json:"premium_subscription_checkout_token"`
}

type getOrCreatePremiumNewsletterSubscriptionResponse struct {
	PremiumNewsletterSubscription billing.PremiumNewsletterSubscription `json:"premium_newsletter_subscription"`
}

func getOrCreatePremiumNewsletterSubscription(userAuth routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req getOrCreatePremiumNewsletterSubscriptionRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.PremiumSubscriptionCheckoutToken, routes.PremiumSubscriptionCheckoutKey)
	if err != nil || *userID != userAuth.UserID {
		r.RespondWithStatus(http.StatusForbidden)
		return nil, nil
	}
	var premiumNewsletterSubscription *billing.PremiumNewsletterSubscription
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		premiumNewsletterSubscription, err = billing.LookupPremiumNewsletterSubscriptionForUser(r, tx, *userID)
		switch {
		case err != nil:
			return err
		case premiumNewsletterSubscription != nil:
			return nil
		}
		premiumNewsletterSubscriptionID := billing.NewPremiumNewsletterSubscriptionID()
		if err := billing.InsertPremiumNewsletterSyncRequest(tx, premiumNewsletterSubscriptionID, billing.PremiumNewsletterSubscriptionUpdateTypeTransitionToActive); err != nil {
			return err
		}
		premiumNewsletterSubscription, err = billing.CreatePremiumNewsletterSubscriptionForUserWithID(r, tx, *userID, premiumNewsletterSubscriptionID)
		return err
	}); err != nil {
		return nil, err
	}
	return getOrCreatePremiumNewsletterSubscriptionResponse{
		PremiumNewsletterSubscription: *premiumNewsletterSubscription,
	}, nil
}
