package stripe

import (
	"babblegraph/services/web/router"
)

func RegisterRouteGroups() error {
	return router.RegisterRouteGroup(router.RouteGroup{
		Prefix: "stripe",
		AuthenticatedRoutes: []router.AuthenticatedRoute{
			{
				Path:    "get_or_create_user_subscription_1",
				Handler: getOrCreateUserSubscription,
			}, {
				Path:    "get_user_nonterm_stripe_subscription_1",
				Handler: getUserNonTerminatedStripeSubscription,
			}, {
				Path:    "delete_stripe_subscription_for_user_1",
				Handler: deleteStripeSubscriptionForUser,
			},
		},
	})
}