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
				routermiddleware.MaybeWithAuthentication(getOrCreateBillingInformation),
			),
		}, {
			Path: "get_or_create_premium_newsletter_subscription_1",
			Handler: routermiddleware.WithRequestBodyLogger(
				routermiddleware.MaybeWithAuthentication(getOrCreatePremiumNewsletterSubscription),
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

func getOrCreateBillingInformation(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	if userAuth == nil {
		r.RespondWithStatus(http.StatusForbidden)
		return nil, nil
	}
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

func getOrCreatePremiumNewsletterSubscription(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	if userAuth == nil {
		r.RespondWithStatus(http.StatusForbidden)
		return nil, nil
	}
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
		premiumNewsletterSubscription, err = billing.GetOrCreatePremiumNewsletterSubscriptionForUser(r, tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	return getOrCreatePremiumNewsletterSubscriptionResponse{
		PremiumNewsletterSubscription: *premiumNewsletterSubscription,
	}, nil
}
