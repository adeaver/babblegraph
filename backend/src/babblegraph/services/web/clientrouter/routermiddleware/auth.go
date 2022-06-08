package routermiddleware

import (
	"babblegraph/model/admin"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/util/auth"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"fmt"
	"net/http"
	"time"

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

var safeCookies = map[string]bool{
	AuthTokenCookieName:         true,
	"__stripe_mid":              true,
	"__stripe_sid":              true,
	utmTrackingIDCookieName:     true,
	admin.AccessTokenCookieName: true,
}

func RemoveUnsafeCookies(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if _, ok := safeCookies[cookie.Name]; !ok {
			http.SetCookie(w, &http.Cookie{
				Name:     cookie.Name,
				Value:    "",
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Now().Add(-5 * time.Minute),
			})
		}
	}
}

type WithSafeCookieHandler func(r *router.Request) (interface{}, error)

func WithSafeCookie(handler WithSafeCookieHandler) router.RequestHandler {
	return func(r *router.Request) (interface{}, error) {
		for _, cookie := range r.GetCookies() {
			if _, ok := safeCookies[cookie.Name]; !ok {
				r.RemoveCookieByName(cookie.Name)
			}
		}
		return handler(r)
	}
}
