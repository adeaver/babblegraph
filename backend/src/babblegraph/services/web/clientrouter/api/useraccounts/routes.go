package useraccounts

import (
	"babblegraph/model/billing"
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/clientrouter/util/auth"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"net/http"
	"time"

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
		}, {
			Path: "create_user_1",
			Handler: routermiddleware.WithNoBodyRequestLogger(
				routermiddleware.MaybeWithAuthentication(createUser),
			),
		},
	},
}

type getUserProfileInformationRequest struct {
	Token    string                      `json:"token"`
	Key      routes.RouteEncryptionKey   `json:"key"`
	NextKeys []routes.RouteEncryptionKey `json:"next_keys,omitempty"`
}

type getUserProfileInformationResponse struct {
	UserProfile *userProfileInformation      `json:"user_profile,omitempty"`
	Error       *userProfileInformationError `json:"error,omitempty"`
}

type userProfileInformation struct {
	HasAccount           bool                            `json:"has_account"`
	IsLoggedIn           bool                            `json:"is_logged_in"`
	SubscriptionLevel    *useraccounts.SubscriptionLevel `json:"subscription_level,omitempty"`
	TrialEligibilityDays *int64                          `json:"trial_eligibility_days,omitempty"`
	HasPaymentMethod     bool                            `json:"has_payment_method"`
	NextTokens           []string                        `json:"next_tokens,omitempty"`
}

type userProfileInformationError string

const (
	userProfileInformationErrorInvalidKey userProfileInformationError = "invalid-key"
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
	var nextTokens []string
	if req.NextKeys != nil {
		for _, key := range req.NextKeys {
			nextToken, err := routes.EncryptUserIDWithKey(*userID, key)
			if err != nil {
				return getUserProfileInformationResponse{
					Error: userProfileInformationErrorInvalidKey.Ptr(),
				}, nil
			}
			nextTokens = append(nextTokens, *nextToken)
		}
	}
	if userAuth != nil {
		if userAuth.UserID == *userID {
			var subscriptionLevel *useraccounts.SubscriptionLevel
			var hasPaymentMethod bool
			var trialEligibilityDays *int64
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				subscriptionLevel, err = useraccounts.LookupSubscriptionLevelForUser(tx, *userID)
				if err != nil {
					return err
				}
				trialEligibilityDays, err = billing.GetPremiumNewsletterSubscriptionTrialEligibilityForUser(tx, *userID)
				if err != nil {
					return err
				}
				premiumSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(r, tx, *userID)
				if err != nil {
					return err
				}
				hasPaymentMethod = premiumSubscription != nil && premiumSubscription.PaymentState != billing.PaymentStateTrialNoPaymentMethod
				return nil
			}); err != nil {
				return nil, err
			}
			return getUserProfileInformationResponse{
				UserProfile: &userProfileInformation{
					HasAccount:           true,
					IsLoggedIn:           true,
					HasPaymentMethod:     hasPaymentMethod,
					SubscriptionLevel:    subscriptionLevel,
					TrialEligibilityDays: trialEligibilityDays,
					NextTokens:           nextTokens,
				},
			}, nil
		}
		r.Infof("User ID does not match user that is logged in, removing cookie...")
		r.RespondWithCookie(&http.Cookie{
			Name:     routermiddleware.AuthTokenCookieName,
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(-5 * 60 * time.Second),
		})
	}
	var hasPaymentMethod bool
	var doesUserHaveAccount bool
	var trialEligibilityDays *int64
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		doesUserHaveAccount, err = useraccounts.DoesUserAlreadyHaveAccount(tx, *userID)
		if err != nil {
			return err
		}
		trialEligibilityDays, err = billing.GetPremiumNewsletterSubscriptionTrialEligibilityForUser(tx, *userID)
		if err != nil {
			return err
		}
		premiumSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(r, tx, *userID)
		if err != nil {
			return err
		}
		subscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, *userID)
		if err != nil {
			return err
		}
		if subscriptionLevel != nil && *subscriptionLevel == useraccounts.SubscriptionLevelPremium {
			hasPaymentMethod = premiumSubscription != nil && premiumSubscription.PaymentState != billing.PaymentStateTrialNoPaymentMethod
		}
		return nil
	}); err != nil {
		return nil, err
	}
	userProfile := &userProfileInformation{
		HasAccount:           doesUserHaveAccount,
		TrialEligibilityDays: trialEligibilityDays,
		HasPaymentMethod:     hasPaymentMethod,
	}
	if !doesUserHaveAccount {
		userProfile.NextTokens = nextTokens
	}
	return getUserProfileInformationResponse{
		UserProfile: userProfile,
	}, nil
}

type createUserRequest struct {
	CreateUserToken string `json:"create_user_token"`
	EmailAddress    string `json:"email_address"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type createUserResponse struct {
	CreateUserError *createUserError `json:"create_user_error"`
}

type createUserError string

const (
	createUserErrorAlreadyExists        createUserError = "already-exists"
	createUserErrorInvalidToken         createUserError = "invalid-token"
	createUserErrorPasswordRequirements createUserError = "pass-requirements"
	createUserErrorNoSubscription       createUserError = "no-subscription"
	createUserErrorPasswordsNoMatch     createUserError = "passwords-no-match"
	createUserErrorInvalidState         createUserError = "invalid-state"
)

func (c createUserError) Ptr() *createUserError {
	return &c
}

func createUser(userAuth *routermiddleware.UserAuthentication, r *router.Request) (interface{}, error) {
	if userAuth != nil {
		return createUserResponse{
			CreateUserError: createUserErrorAlreadyExists.Ptr(),
		}, nil
	}
	var req createUserRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	userID, err := routetoken.ValidateTokenAndEmailAndGetUserID(req.CreateUserToken, routes.CreateUserKey, formattedEmailAddress)
	switch {
	case err != nil:
		return createUserResponse{
			CreateUserError: createUserErrorInvalidToken.Ptr(),
		}, nil
	case req.Password != req.ConfirmPassword:
		return createUserResponse{
			CreateUserError: createUserErrorPasswordsNoMatch.Ptr(),
		}, nil
	case !useraccounts.ValidatePasswordMeetsRequirements(req.Password):
		return createUserResponse{
			CreateUserError: createUserErrorPasswordRequirements.Ptr(),
		}, nil
	}
	var cErr *createUserError
	err = database.WithTx(func(tx *sqlx.Tx) error {
		user, err := users.GetUser(tx, *userID)
		if err != nil {
			return err
		}
		if user.Status != users.UserStatusVerified {
			cErr = createUserErrorInvalidState.Ptr()
			return nil
		}
		alreadyHasAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, *userID)
		switch {
		case err != nil:
			return err
		case alreadyHasAccount:
			cErr = createUserErrorAlreadyExists.Ptr()
			return nil
		}
		return useraccounts.CreateUserPasswordForUser(tx, *userID, req.Password)
	})
	switch {
	case err != nil:
		return nil, err
	case cErr != nil:
		return createUserResponse{
			CreateUserError: cErr,
		}, nil
	}
	token, err := auth.CreateJWTForUser(*userID)
	if err != nil {
		return nil, err
	}
	r.RespondWithCookie(&http.Cookie{
		Name:     routermiddleware.AuthTokenCookieName,
		Value:    *token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(auth.SessionExpirationTime),
	})
	return createUserResponse{}, nil
}
