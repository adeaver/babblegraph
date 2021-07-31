package index

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/middleware"
	"babblegraph/util/env"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
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
