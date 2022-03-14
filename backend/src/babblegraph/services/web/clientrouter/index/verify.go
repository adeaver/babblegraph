package index

import (
	"babblegraph/actions/verification"
	"babblegraph/model/routes"
	"babblegraph/services/web/router"
)

func handleVerification(r *router.Request) (interface{}, error) {
	token, err := r.GetRouteVar("token")
	if err != nil {
		return nil, err
	}
	userID, err := verification.VerifyUserByToken(*token)
	if err != nil {
		return nil, err
	}
	return routes.MakeSubscriptionManagementRouteForUserID(*userID)
}
