package index

import (
	"babblegraph/model/routes"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/initialdata"
	"babblegraph/services/web/middleware"
	"log"
	"net/http"
)

func HandleCheckoutPage(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.WithAuthorizationCheck(w, r, middleware.WithAuthorizationCheckInput{
			HandleFoundUser: func(userID users.UserID, subscriptionLevel *useraccounts.SubscriptionLevel, w http.ResponseWriter, r *http.Request) {
				log.Println("Here 3")
				if subscriptionLevel != nil {
					subscriptionLink, err := routes.MakeSubscriptionManagementRouteForUserID(userID)
					if err != nil {
						http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
						return
					}
					http.Redirect(w, r, *subscriptionLink, http.StatusTemporaryRedirect)
					return
				}
				HandleServeIndexPage(staticFileDirName, initialdata.InitialFrontendDataOptions{
					IncludeStripePublicKey: true,
				})(w, r)
			},
			HandleNoUserFound: func(w http.ResponseWriter, r *http.Request) {
				log.Println("Here")
				http.Redirect(w, r, routes.MakeLoginLinkWithPremiumSubscriptionCheckoutRedirect(), http.StatusTemporaryRedirect)
			},
			HandleInvalidAuthenticationToken: func(w http.ResponseWriter, r *http.Request) {
				log.Println("Here 2")
				http.Redirect(w, r, routes.MakeLoginLinkWithPremiumSubscriptionCheckoutRedirect(), http.StatusTemporaryRedirect)
			},
			HandleError: middleware.HandleAuthorizationError,
		})
	}
}
