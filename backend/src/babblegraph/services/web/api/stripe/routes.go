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
			},
		},
	}); err != nil {
		return err
	}
	return router.RegisterRouteWithoutWrapper(routeGroupPrefix, "handle_stripe_event_1", handleStripeWebhook)
}
