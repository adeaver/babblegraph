package stripe

import (
	"babblegraph/services/web/router"
)

const routeGroupPrefix = "stripe"

func RegisterRouteGroups() error {
	if err := router.RegisterRouteGroup(router.RouteGroup{
		Prefix: routeGroupPrefix,
		Routes: []router.Route{},
		AuthenticatedRoutes: []router.AuthenticatedRoute{
			{
				Path:    "create_user_subscription_1",
				Handler: createUserSubscription,
			}, {
				Path:    "get_user_nonterm_stripe_subscription_1",
				Handler: getUserNonTerminatedStripeSubscription,
			}, {
				Path:    "delete_stripe_subscription_for_user_1",
				Handler: deleteStripeSubscriptionForUser,
			}, {
				Path:    "update_stripe_subscription_for_user_1",
				Handler: updateStripeSubscriptionFrequencyForUser,
			}, {
				Path:    "get_setup_intent_for_user_1",
				Handler: getSetupIntentForUser,
			}, {
				Path:    "insert_payment_method_for_user_1",
				Handler: insertNewPaymentMethodForUser,
			}, {
				Path:    "set_default_payment_method_for_user_1",
				Handler: setDefaultPaymentMethodForUser,
			}, {
				Path:    "get_payment_methods_for_user_1",
				Handler: getPaymentMethodsForUser,
			}, {
				Path:    "get_payment_method_by_id_1",
				Handler: getPaymentMethodByID,
			},
		},
	}); err != nil {
		return err
	}
	return router.RegisterRouteWithoutWrapper(routeGroupPrefix, "handle_stripe_event_1", handleStripeWebhook)
}
