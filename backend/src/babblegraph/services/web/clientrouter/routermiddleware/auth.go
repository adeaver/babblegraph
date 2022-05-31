package routermiddleware

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/clienterror"
	"babblegraph/services/web/clientrouter/util/auth"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

const AuthTokenCookieName = "session_token"

type UserAuthentication struct {
	UserID            users.UserID
	SubscriptionLevel *useraccounts.SubscriptionLevel
}

type MaybeWithAuthenticationHandler func(auth *UserAuthentication, r *router.Request) (interface{}, error)

func MaybeWithAuthentication(handler MaybeWithAuthenticationHandler) router.RequestHandler {
	return func(r *router.Request) (interface{}, error) {
		var userAuth *UserAuthentication
		for _, cookie := range r.GetCookies() {
			if cookie.Name == AuthTokenCookieName {
				token := cookie.Value
				userID, isValid, err := auth.VerifyJWTAndGetUserID(token)
				switch {
				case err != nil:
					return nil, err
				case !isValid,
					userID == nil:
					// no-op
				default:
					var userStatus users.UserStatus
					var userSubscriptionLevel *useraccounts.SubscriptionLevel
					if err := database.WithTx(func(tx *sqlx.Tx) error {
						user, err := users.GetUser(tx, *userID)
						if err != nil {
							return err
						}
						userStatus = user.Status
						userSubscriptionLevel, err = useraccounts.LookupSubscriptionLevelForUser(tx, *userID)
						return err
					}); err != nil {
						return nil, err
					}
					switch userStatus {
					case users.UserStatusVerified:
						userAuth = &UserAuthentication{
							UserID:            *userID,
							SubscriptionLevel: userSubscriptionLevel,
						}
					case users.UserStatusUnverified,
						users.UserStatusUnsubscribed,
						users.UserStatusBlocklistBounced,
						users.UserStatusBlocklistComplaint:
						// no-op
					default:
						return nil, fmt.Errorf("Invalid user state: %s", userStatus)
					}
				}
			}
		}
		return handler(userAuth, r)
	}
}

type WithAuthenticationHandler func(auth UserAuthentication, r *router.Request) (interface{}, error)

func WithAuthentication(handler WithAuthenticationHandler) router.RequestHandler {
	return MaybeWithAuthentication(func(auth *UserAuthentication, r *router.Request) (interface{}, error) {
		if auth == nil {
			r.RespondWithStatus(http.StatusForbidden)
			return nil, nil
		}
		return handler(*auth, r)
	})
}

type WithRequiredSubscriptionHandler func(userID users.UserID, r *router.Request) (interface{}, error)

func WithRequiredSubscription(validSubscriptionLevels []useraccounts.SubscriptionLevel, handler WithRequiredSubscriptionHandler) router.RequestHandler {
	return MaybeWithAuthentication(func(auth *UserAuthentication, r *router.Request) (interface{}, error) {
		switch {
		case auth == nil,
			auth.SubscriptionLevel == nil:
			// no-op
		default:
			for _, subscriptionLevel := range validSubscriptionLevels {
				if subscriptionLevel == *auth.SubscriptionLevel {
					return handler(auth.UserID, r)
				}
			}
		}
		r.RespondWithStatus(http.StatusForbidden)
		return nil, nil
	})
}

type ValidateUserAuthWithTokenInput struct {
	RequireEmailAddress bool
	EmailAddress        *string
	Token               string
	KeyType             routes.RouteEncryptionKey
}

func ValidateUserAuthWithToken(userAuthentication *UserAuthentication, input ValidateUserAuthWithTokenInput) (*users.UserID, *clienterror.Error, error) {
	var tokenUserID *users.UserID
	var err error
	if input.EmailAddress != nil {
		formattedEmailAddress := email.FormatEmailAddress(*input.EmailAddress)
		tokenUserID, err = routetoken.ValidateTokenAndEmailAndGetUserID(input.Token, input.KeyType, formattedEmailAddress)
		if err != nil {
			return nil, nil, err
		}
	} else {
		tokenUserID, err = routetoken.ValidateTokenAndGetUserID(input.Token, input.KeyType)
		if err != nil {
			return nil, nil, err
		}
	}
	var doesUserHaveAccount bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		doesUserHaveAccount, err = useraccounts.DoesUserAlreadyHaveAccount(tx, *tokenUserID)
		return err
	}); err != nil {
		return nil, nil, err
	}
	switch {
	case userAuthentication == nil && !doesUserHaveAccount:
		if input.EmailAddress == nil && input.RequireEmailAddress {
			return nil, clienterror.ErrorInvalidEmailAddress.Ptr(), nil
		}
		return tokenUserID, nil, nil
	case userAuthentication == nil && doesUserHaveAccount:
		return nil, clienterror.ErrorNoAuth.Ptr(), nil
	case userAuthentication != nil && !doesUserHaveAccount:
		return nil, clienterror.ErrorIncorrectKey.Ptr(), nil
	case userAuthentication != nil && doesUserHaveAccount:
		if userAuthentication.UserID != *tokenUserID {
			return nil, clienterror.ErrorIncorrectKey.Ptr(), nil
		}
		return tokenUserID, nil, nil
	default:
		return nil, nil, fmt.Errorf("Unreachable")
	}
}
