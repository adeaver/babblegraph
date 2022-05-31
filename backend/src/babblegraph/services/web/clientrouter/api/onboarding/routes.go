package onboarding

import (
	"babblegraph/model/onboarding"
	"babblegraph/model/routes"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/ptr"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "onboarding",
	Routes: []router.Route{
		{
			Path: "get_onboarding_status_for_user_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				getOnboardingStatusForUser,
			),
		},
	},
}

type getOnboardingStatusForUserRequest struct {
	OnboardingToken string `json:"onboarding_token"`
}

type getOnboardingStatusForUserResponse struct {
	EmailAddress                *string                    `json:"email_address,omitempty"`
	OnboardingStatus            onboarding.OnboardingStage `json:"onboarding_status"`
	SubscriptionManagementToken *string                    `json:"subscription_management_token,omitempty"`
}

func getOnboardingStatusForUser(r *router.Request) (interface{}, error) {
	var req getOnboardingStatusForUserRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.OnboardingToken, routes.OnboardingKey)
	if err != nil {
		return nil, err
	}
	var emailAddress *string
	var onboardingStatus *onboarding.OnboardingStage
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		user, err := users.GetUser(tx, *userID)
		switch {
		case err != nil:
			return err
		case user.Status != users.UserStatusVerified:
			// no-op
		default:
			emailAddress = ptr.String(user.EmailAddress)
		}
		onboardingStatus, err = onboarding.GetOnboardingStageForUser(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	var subscriptionManagementToken *string
	if emailAddress != nil {
		subscriptionManagementToken, err = routes.MakeSubscriptionManagementToken(*userID)
		if err != nil {
			return nil, err
		}
	}
	return getOnboardingStatusForUserResponse{
		EmailAddress:                emailAddress,
		OnboardingStatus:            *onboardingStatus,
		SubscriptionManagementToken: subscriptionManagementToken,
	}, nil
}
