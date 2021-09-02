package index

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/middleware"
	"babblegraph/services/web/util/routetoken"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func HandleLoginPage(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.WithAuthorizationCheck(w, r, middleware.WithAuthorizationCheckInput{
			HandleFoundUser: func(userID users.UserID, subscriptionLevel *useraccounts.SubscriptionLevel, w http.ResponseWriter, r *http.Request) {
				params := r.URL.Query()
				var redirectLocation string
				if redirectLocationParams, _ := params[routes.RedirectKeyParameter]; len(redirectLocationParams) > 0 {
					redirectLocation = redirectLocationParams[0]
				}
				redirectKeyForLocation := routes.GetLoginRedirectKeyOrDefault(redirectLocation)
				redirectURL, err := routes.GetLoginRedirectRouteForKeyAndUser(redirectKeyForLocation, userID)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					sentry.CaptureException(fmt.Errorf("Error redirecting on login: %s", err.Error()))
					return
				}
				http.Redirect(w, r, *redirectURL, http.StatusTemporaryRedirect)
			},
			HandleNoUserFound:                HandleServeIndexPage(staticFileDirName),
			HandleInvalidAuthenticationToken: HandleServeIndexPage(staticFileDirName),
			HandleError:                      middleware.HandleAuthorizationError,
		})
	}
}

func HandleLogout() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.WithAuthorizationCheck(w, r, middleware.WithAuthorizationCheckInput{
			HandleFoundUser: func(userID users.UserID, subscriptionLevel *useraccounts.SubscriptionLevel, w http.ResponseWriter, r *http.Request) {
				middleware.RemoveAuthToken(w)
				http.Redirect(w, r, env.GetAbsoluteURLForEnvironment("login"), http.StatusTemporaryRedirect)
			},
			HandleNoUserFound: func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, env.GetAbsoluteURLForEnvironment("login"), http.StatusTemporaryRedirect)
			},
			HandleInvalidAuthenticationToken: func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, env.GetAbsoluteURLForEnvironment("login"), http.StatusTemporaryRedirect)
			},
			HandleError: middleware.HandleAuthorizationError,
		})
	}
}

func HandleCreateUserPage(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		routeVars := mux.Vars(r)
		token, ok := routeVars["token"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		middleware.WithAuthorizationCheck(w, r, middleware.WithAuthorizationCheckInput{
			HandleFoundUser: func(userID users.UserID, subscriptionLevel *useraccounts.SubscriptionLevel, w http.ResponseWriter, r *http.Request) {
				destinationURL, err := routes.MakeSubscriptionManagementRouteForUserID(userID)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					sentry.CaptureException(fmt.Errorf("Error redirecting on login: %s", err.Error()))
					return
				}
				if subscriptionLevel == nil {
					destinationURL, err = routes.MakePremiumSubscriptionCheckoutLink(userID)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						sentry.CaptureException(fmt.Errorf("Error redirecting on login: %s", err.Error()))
						return
					}
				}
				http.Redirect(w, r, *destinationURL, http.StatusTemporaryRedirect)
			},
			HandleNoUserFound: func(w http.ResponseWriter, r *http.Request) {
				userID, err := routetoken.ValidateTokenAndGetUserID(token, routes.CreateUserKey)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				var doesUserHaveAccount bool
				if err := database.WithTx(func(tx *sqlx.Tx) error {
					var err error
					doesUserHaveAccount, err = useraccounts.DoesUserAlreadyHaveAccount(tx, *userID)
					return err
				}); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if doesUserHaveAccount {
					http.Redirect(w, r, routes.MakeLoginLinkWithPremiumSubscriptionCheckoutRedirect(), http.StatusTemporaryRedirect)
					return
				}
				HandleServeIndexPage(staticFileDirName)
			},
			HandleInvalidAuthenticationToken: func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, routes.MakeLoginLinkWithPremiumSubscriptionCheckoutRedirect(), http.StatusTemporaryRedirect)
			},
			HandleError: middleware.HandleAuthorizationError,
		})
	}
}
