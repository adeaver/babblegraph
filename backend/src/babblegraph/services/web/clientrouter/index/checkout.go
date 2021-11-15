package index

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/clientrouter/middleware"
	"net/http"
)

func HandleCheckoutPage(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.WithAuthorizationCheck(w, r, middleware.WithAuthorizationCheckInput{
			HandleFoundUser: func(userID users.UserID, subscriptionLevel *useraccounts.SubscriptionLevel, w http.ResponseWriter, r *http.Request) {
				if subscriptionLevel != nil {
					subscriptionLink, err := routes.MakeSubscriptionManagementRouteForUserID(userID)
					if err != nil {
						http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
						return
					}
					http.Redirect(w, r, *subscriptionLink, http.StatusTemporaryRedirect)
					return
				}
				HandleServeIndexPage(staticFileDirName)(w, r)
			},
			HandleNoUserFound: func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, routes.MakeLoginLinkWithPremiumSubscriptionCheckoutRedirect(), http.StatusTemporaryRedirect)
			},
			HandleInvalidAuthenticationToken: func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, routes.MakeLoginLinkWithPremiumSubscriptionCheckoutRedirect(), http.StatusTemporaryRedirect)
			},
			HandleError: middleware.HandleAuthorizationError,
		})
	}
}
