package routermiddleware

import (
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/clienterror"
	"babblegraph/services/web/clientrouter/util/auth"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"babblegraph/util/ptr"
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

type ValidateUserAuthInput struct {
	DecodedUserID users.UserID
	EmailAddress  *string
}

func ValidateUserAuth(userAuth *UserAuthentication, input ValidateUserAuthInput) (*clienterror.Error, error) {
	if userAuth == nil {
		var formattedEmailAddress *string
		if input.EmailAddress != nil {
			formattedEmailAddress = ptr.String(email.FormatEmailAddress(*input.EmailAddress))
		}
		var cErr *clienterror.Error
		var doesUserHaveAccount bool
		err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			doesUserHaveAccount, err = useraccounts.DoesUserAlreadyHaveAccount(tx, input.DecodedUserID)
			if err != nil {
				return err
			}
			if formattedEmailAddress != nil {
				user, err := users.LookupUserForIDAndEmail(tx, input.DecodedUserID, *formattedEmailAddress)
				switch {
				case err != nil:
					return err
				case user == nil:
					cErr = clienterror.ErrorInvalidEmailAddress.Ptr()
					return fmt.Errorf("No user")
				}
			}
			return nil
		})
		switch {
		case err != nil:
			return nil, err
		case cErr != nil:
			return cErr, nil
		case doesUserHaveAccount:
			// User has an account but is not logged in
			return clienterror.ErrorNoAuth.Ptr(), nil
		}
		return nil, nil
	}
	if input.DecodedUserID != userAuth.UserID {
		return clienterror.ErrorNoAuth.Ptr(), nil
	}
	return nil, nil
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
