package index

import (
	"babblegraph/actions/verification"
	"babblegraph/model/routes"
	"babblegraph/services/web/router"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleVerificationForToken(w http.ResponseWriter, r *http.Request) {
	router.LogRequestWithoutBody(r)
	routeVars := mux.Vars(r)
	token, ok := routeVars["token"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := verification.VerifyUserByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	subscriptionManagementLink, err := routes.MakeSubscriptionManagementRouteForUserID(*userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, *subscriptionManagementLink, http.StatusMovedPermanently)
}
