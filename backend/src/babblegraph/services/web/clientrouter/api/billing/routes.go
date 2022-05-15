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
			Path: "lookup_active_premium_newsletter_subscription_1",
			Handler: routermiddleware.WithRequestBodyLogger(
				routermiddleware.WithAuthentication(lookupActivePremiumNewsletterSubscription),
			),
		}, {
			Path: "set_premium_newsletter_subscription_auto_renew_1",
			Handler: routermiddleware.WithRequestBodyLogger(
				routermiddleware.WithAuthentication(setPremiumNewsletterSubscriptionAutoRenew),
			),
		}, {
			Path: "prepare_premium_newsletter_subscription_sync_1",
			Handler: routermiddleware.WithRequestBodyLogger(
				routermiddleware.WithAuthentication(preparePremiumNewsletterSubscriptionSync),
			),
		}, {
			Path: "get_payment_methods_for_user_1",
			Handler: routermiddleware.WithRequestBodyLogger(
				routermiddleware.WithAuthentication(getPaymentMethodsForUser),
			),
		}, {
			Path: "mark_payment_method_as_default_1",
			Handler: routermiddleware.WithRequestBodyLogger(
				routermiddleware.WithAuthentication(markPaymentMethodAsDefault),
			),
		}, {
			Path: "delete_payment_method_for_user_1",
			Handler: routermiddleware.WithRequestBodyLogger(
				routermiddleware.WithAuthentication(deletePaymentMethodForUser),
			),
		}, {
			Path: "stripe_begin_payment_method_setup_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.WithAuthentication(stripeBeginPaymentMethodSetup),
			),
		}, {
			Path: "stripe_handle_webhook_event_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				stripeHandleWebhookEvent,
			),
		}, {
			Path: "lookup_promotion_code_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.WithMaybePromotion(lookupPromotionCode),
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
	premiumNewsletterSubscriptionID := billing.NewPremiumNewsletterSubscriptionID()
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
		if err := billing.InsertPremiumNewsletterSyncRequest(tx, premiumNewsletterSubscriptionID, billing.PremiumNewsletterSubscriptionUpdateTypeTransitionToActive); err != nil {
			return err
		}
		premiumNewsletterSubscription, err = billing.CreatePremiumNewsletterSubscriptionForUserWithID(r, tx, *userID, premiumNewsletterSubscriptionID)
		// THIS IS A HACK
		premiumNewsletterSubscription.ID = &premiumNewsletterSubscriptionID
		return err
	}); err != nil {
		return nil, err
	}
	return getOrCreatePremiumNewsletterSubscriptionResponse{
		PremiumNewsletterSubscription: *premiumNewsletterSubscription,
	}, nil
}

type lookupActivePremiumNewsletterSubscriptionRequest struct {
	SubscriptionManagementToken string `json:"subscription_management_token"`
}

type lookupActivePremiumNewsletterSubscriptionResponse struct {
	PremiumNewsletterSubscription *billing.PremiumNewsletterSubscription `json:"premium_newsletter_subscription,omitempty"`
}

func lookupActivePremiumNewsletterSubscription(userAuth routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req lookupActivePremiumNewsletterSubscriptionRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil || *userID != userAuth.UserID {
		r.RespondWithStatus(http.StatusForbidden)
		return nil, nil
	}
	var premiumNewsletterSubscription *billing.PremiumNewsletterSubscription
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		premiumNewsletterSubscription, err = billing.LookupPremiumNewsletterSubscriptionForUser(r, tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	return lookupActivePremiumNewsletterSubscriptionResponse{
		PremiumNewsletterSubscription: premiumNewsletterSubscription,
	}, nil
}

type setPremiumNewsletterSubscriptionAutoRenewRequest struct {
	SubscriptionManagementToken string `json:"subscription_management_token"`
	IsAutoRenewEnabled          bool   `json:"is_auto_renew_enabled"`
}

type setPremiumNewsletterSubscriptionAutoRenewResponse struct {
	Success bool `json:"success"`
}

func setPremiumNewsletterSubscriptionAutoRenew(userAuth routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req setPremiumNewsletterSubscriptionAutoRenewRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.SubscriptionManagementToken, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil || *userID != userAuth.UserID {
		r.RespondWithStatus(http.StatusForbidden)
		return nil, nil
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return billing.UpdateSubscriptionAutoRenewForUser(tx, *userID, req.IsAutoRenewEnabled)
	}); err != nil {
		return nil, err
	}
	return setPremiumNewsletterSubscriptionAutoRenewResponse{
		Success: true,
	}, nil
}

type preparePremiumNewsletterSubscriptionSyncRequest struct {
	ID         billing.PremiumNewsletterSubscriptionID `json:"id"`
	UpdateType string                                  `json:"update_type"`
}

type preparePremiumNewsletterSubscriptionSyncResponse struct {
	Success bool `json:"success"`
}

func preparePremiumNewsletterSubscriptionSync(userAuth routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req preparePremiumNewsletterSubscriptionSyncRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	updateType, err := billing.GetPremiumNewsletterSubscriptionUpdateTypeFromString(req.UpdateType)
	if err != nil {
		return preparePremiumNewsletterSubscriptionSyncResponse{
			Success: false,
		}, nil
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return billing.InsertPremiumNewsletterSyncRequest(tx, req.ID, *updateType)
	}); err != nil {
		return nil, err
	}
	return preparePremiumNewsletterSubscriptionSyncResponse{
		Success: true,
	}, nil
}

type getPaymentMethodsForUserRequest struct{}

type getPaymentMethodsForUserResponse struct {
	PaymentMethods []billing.PaymentMethod `json:"payment_methods"`
}

func getPaymentMethodsForUser(userAuth routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var paymentMethods []billing.PaymentMethod
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		paymentMethods, err = billing.GetPaymentMethodsForUser(tx, userAuth.UserID)
		return err
	}); err != nil {
		return nil, err
	}
	return getPaymentMethodsForUserResponse{
		PaymentMethods: paymentMethods,
	}, nil
}

type markPaymentMethodAsDefaultRequest struct {
	PaymentMethodID string `json:"payment_method_id"`
}

type markPaymentMethodAsDefaultResponse struct {
	Success bool `json:"success"`
}

func markPaymentMethodAsDefault(userAuth routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req markPaymentMethodAsDefaultRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return billing.MarkPaymentMethodAsDefaultForUser(tx, userAuth.UserID, req.PaymentMethodID)
	}); err != nil {
		return nil, err
	}
	return markPaymentMethodAsDefaultResponse{
		Success: true,
	}, nil
}

type deletePaymentMethodForUserRequest struct {
	PaymentMethodID string `json:"payment_method_id"`
}

type deletePaymentMethodForUserResponse struct {
	Success bool `json:"success"`
}

func deletePaymentMethodForUser(userAuth routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req deletePaymentMethodForUserRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return billing.DeletePaymentMethodForUser(tx, userAuth.UserID, req.PaymentMethodID)
	}); err != nil {
		return nil, err
	}
	return deletePaymentMethodForUserResponse{
		Success: true,
	}, nil
}

type lookupPromotionCodeResponse struct {
	PromotionCode *billing.PromotionCode `json:"promotion_code,omitempty"`
}

func lookupPromotionCode(promotionCode *billing.PromotionCode, r *router.Request) (interface{}, error) {
	return lookupPromotionCodeResponse{
		PromotionCode: promotionCode,
	}, nil
}
