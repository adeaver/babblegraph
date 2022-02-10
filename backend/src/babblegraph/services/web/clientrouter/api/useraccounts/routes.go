package useraccounts

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "useraccounts",
	Routes: []router.Route{
		{
			Path: "get_user_profile_information_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.MaybeWithAuthentication(getUserProfileInformation),
			),
		},
	},
}

type getUserProfileInformationRequest struct {
	Token string                     `json:"token"`
	Key   routes.RouteEncryptionKey  `json:"key"`
	ToKey *routes.RouteEncryptionKey `json:"to_key,omitempty"`
}

type getUserProfileInformationResponse struct {
	UserProfile *userProfileInformation      `json:"user_profile,omitempty"`
	Error       *userProfileInformationError `json:"error,omitempty"`
}

type userProfileInformation struct {
	HasAccount        bool                            `json:"has_account"`
	IsLoggedIn        bool                            `json:"is_logged_in"`
	SubscriptionLevel *useraccounts.SubscriptionLevel `json:"subscription_level,omitempty"`
	NextToken         *string                         `json:"next_token"`
}

type userProfileInformationError string

const (
	userProfileInformationErrorInvalidKey     userProfileInformationError = "invalid-key"
	userProfileInformationErrorInvalidAccount userProfileInformationError = "invalid-token" // This is deliberately ambiguous since it's exposed to the client
)

func (u userProfileInformationError) Ptr() *userProfileInformationError {
	return &u
}

func getUserProfileInformation(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	var req getUserProfileInformationRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.Token, req.Key)
	if err != nil {
		return getUserProfileInformationResponse{
			Error: userProfileInformationErrorInvalidKey.Ptr(),
		}, nil
	}
	var nextToken *string
	if req.ToKey != nil {
		nextToken, err = routes.EncryptUserIDWithKey(*userID, *req.ToKey)
		if err != nil {
			return getUserProfileInformationResponse{
				Error: userProfileInformationErrorInvalidKey.Ptr(),
			}, nil
		}
	}
	if userAuth == nil {
		var doesUserHaveAccount bool
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			doesUserHaveAccount, err = useraccounts.DoesUserAlreadyHaveAccount(tx, *userID)
			return err
		}); err != nil {
			return nil, err
		}
		userProfile := &userProfileInformation{
			HasAccount: doesUserHaveAccount,
		}
		if !doesUserHaveAccount {
			userProfile.NextToken = nextToken
		}
		return getUserProfileInformationResponse{
			UserProfile: userProfile,
		}, nil
	}
	if userAuth.UserID != *userID {
		// TODO: clear token
		return getUserProfileInformationResponse{
			Error: userProfileInformationErrorInvalidAccount.Ptr(),
		}, nil
	}
	var subscriptionLevel *useraccounts.SubscriptionLevel
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		subscriptionLevel, err = useraccounts.LookupSubscriptionLevelForUser(tx, *userID)
		return err
	}); err != nil {
		return nil, err
	}
	return getUserProfileInformationResponse{
		UserProfile: &userProfileInformation{
			HasAccount:        true,
			IsLoggedIn:        true,
			SubscriptionLevel: subscriptionLevel,
			NextToken:         nextToken,
		},
	}, nil
}
