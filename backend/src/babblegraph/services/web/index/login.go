package index

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/middleware"
	"babblegraph/util/env"
	"net/http"
)

func HandleLoginPage(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.WithAuthorizationCheck(w, r, middleware.WithAuthorizationCheckInput{
			HandleFoundSubscribedUser: func(userID users.UserID, subscriptionLevel useraccounts.SubscriptionLevel, w http.ResponseWriter, r *http.Request) {
				subscriptionManagementRoute, err := routes.MakeSubscriptionManagementRouteForUserID(userID)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				http.Redirect(w, r, *subscriptionManagementRoute, http.StatusTemporaryRedirect)
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
			HandleFoundSubscribedUser: func(userID users.UserID, subscriptionLevel useraccounts.SubscriptionLevel, w http.ResponseWriter, r *http.Request) {
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
