package token

import (
	"babblegraph/model/routes"
	"babblegraph/services/web/router"
	"babblegraph/services/web/util/routetoken"
	"encoding/json"
)

func RegisterRouteGroups() error {
	return router.RegisterRouteGroup(router.RouteGroup{
		Prefix: "token",
		Routes: []router.Route{
			{
				Path:    "get_reinforcement_token_1",
				Handler: handleGetReinforcementToken,
			},
		},
	})
}

type getReinforcementTokenRequest struct {
	Token string `json:"token"`
}

type getReinforcementTokenResponse struct {
	Token string `json:"token"`
}

func handleGetReinforcementToken(body []byte) (interface{}, error) {
	var req getReinforcementTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.Token, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		return nil, err
	}
	newToken, err := routes.MakeWordReinforcementToken(*userID)
	if err != nil {
		return nil, err
	}
	return getReinforcementTokenResponse{
		Token: *newToken,
	}, nil
}
