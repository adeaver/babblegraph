package middleware

import (
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/util/auth"
	"babblegraph/util/database"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

const authTokenCookieName = "session_token"

func AssignAuthToken(w http.ResponseWriter, userID users.UserID) error {
	token, err := auth.CreateJWTForUser(userID)
	if err != nil {
		log.Println(fmt.Sprintf("Error creating token for user ID %s: %s", userID, err.Error()))
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     authTokenCookieName,
		Value:    *token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(auth.SessionExpirationTime),
	})
	return nil
}

type WithAuthorizationCheckInput struct {
	HandleFoundSubscribedUser        func(users.UserID, useraccounts.SubscriptionLevel, http.ResponseWriter, *http.Request)
	HandleNoUserFound                func(http.ResponseWriter, *http.Request)
	HandleInvalidAuthenticationToken func(http.ResponseWriter, *http.Request)
	HandleError                      func(error, http.ResponseWriter, *http.Request)
}

func WithAuthorizationCheck(w http.ResponseWriter, r *http.Request, input WithAuthorizationCheckInput) {
	var userSubscriptionLevel *useraccounts.SubscriptionLevel
	var userID *users.UserID
	for _, cookie := range r.Cookies() {
		if cookie.Name == authTokenCookieName {
			var err error
			var isValid bool
			token := cookie.Value
			userID, isValid, err = auth.VerifyJWTAndGetUserID(token)
			switch {
			case err != nil:
				input.HandleError(err, w, r)
				return
			case !isValid:
				input.HandleInvalidAuthenticationToken(w, r)
				return
			case userID == nil:
				input.HandleNoUserFound(w, r)
				return
			default:
				var userSubscriptionLevel *useraccounts.SubscriptionLevel
				if err := database.WithTx(func(tx *sqlx.Tx) error {
					var err error
					userSubscriptionLevel, err = useraccounts.LookupSubscriptionLevelForUser(tx, *userID)
					return err
				}); err != nil {
					input.HandleError(err, w, r)
					return
				}
			}
		}
	}
	if userSubscriptionLevel == nil {
		input.HandleNoUserFound(w, r)
		return
	}
	input.HandleFoundSubscribedUser(*userID, *userSubscriptionLevel, w, r)
}

func WithAuthorizationLevelVerification(validAuthorizationLevels []useraccounts.SubscriptionLevel, fn func(userID users.UserID) func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var userSubscriptionLevel *useraccounts.SubscriptionLevel
		var userID *users.UserID
		for _, cookie := range r.Cookies() {
			if cookie.Name == authTokenCookieName {
				var err error
				var isValid bool
				token := cookie.Value
				userID, isValid, err = auth.VerifyJWTAndGetUserID(token)
				switch {
				case err != nil:
					log.Println(fmt.Sprintf("Error authenticating user: %s", err.Error()))
					w.WriteHeader(http.StatusBadRequest)
					return
				case !isValid:
					// Redirect to login page
				case userID == nil:
					log.Println(fmt.Sprintf("Error authenticating user: No user"))
					w.WriteHeader(http.StatusBadRequest)
					return
				default:
					if err := database.WithTx(func(tx *sqlx.Tx) error {
						var err error
						userSubscriptionLevel, err = useraccounts.LookupSubscriptionLevelForUser(tx, *userID)
						return err
					}); err != nil {
						log.Println(fmt.Sprintf("Error authenticating user: %s", err.Error()))
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					break
				}
			}
		}
		if userSubscriptionLevel == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		for _, validSubscriptionLevel := range validAuthorizationLevels {
			if *userSubscriptionLevel == validSubscriptionLevel {
				fn(*userID)(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
	}
}
