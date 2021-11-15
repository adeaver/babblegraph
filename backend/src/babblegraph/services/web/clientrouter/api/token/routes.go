package token

import (
	"babblegraph/model/routes"
	"babblegraph/services/web/clientrouter/api"
	"babblegraph/services/web/clientrouter/util/routetoken"
	"encoding/json"
)

func RegisterRouteGroups() error {
	return api.RegisterRouteGroup(api.RouteGroup{
		Prefix: "token",
		Routes: []api.Route{
			{
				Path:    "get_reinforcement_token_1",
				Handler: handleGetReinforcementToken,
			}, {
				Path:    "get_manage_token_for_reinforcement_token_1",
				Handler: handleGetManageTokenForReinforcementToken,
			}, {
				Path:    "get_create_user_token_1",
				Handler: handleGetCreateUserToken,
			}, {
				Path:    "get_premium_checkout_token_1",
				Handler: handleGetPremiumCheckoutToken,
			}, {
				Path:    "get_manage_token_for_premium_checkout_token_1",
				Handler: handleGetManageTokenForPremiumCheckoutToken,
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

type getManageTokenForReinforcementTokenRequest struct {
	Token string `json:"token"`
}

type getManageTokenForReinforcementTokenResponse struct {
	Token string `json:"token"`
}

func handleGetManageTokenForReinforcementToken(body []byte) (interface{}, error) {
	var req getManageTokenForReinforcementTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.Token, routes.WordReinforcementKey)
	if err != nil {
		return nil, err
	}
	newToken, err := routes.MakeSubscriptionManagementToken(*userID)
	if err != nil {
		return nil, err
	}
	return getManageTokenForReinforcementTokenResponse{
		Token: *newToken,
	}, nil
}

type getCreateUserTokenRequest struct {
	Token string `json:"token"`
}

type getCreateUserTokenResponse struct {
	Token string `json:"token"`
}

func handleGetCreateUserToken(body []byte) (interface{}, error) {
	var req getCreateUserTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.Token, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		return nil, err
	}
	newToken, err := routes.MakeCreateUserToken(*userID)
	if err != nil {
		return nil, err
	}
	return getCreateUserTokenResponse{
		Token: *newToken,
	}, nil
}

type getPremiumCheckoutTokenRequest struct {
	Token string `json:"token"`
}

type getPremiumCheckoutTokenResponse struct {
	Token string `json:"token"`
}

func handleGetPremiumCheckoutToken(body []byte) (interface{}, error) {
	var req getPremiumCheckoutTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.Token, routes.SubscriptionManagementRouteEncryptionKey)
	if err != nil {
		return nil, err
	}
	newToken, err := routes.MakePremiumSubscriptionCheckoutToken(*userID)
	if err != nil {
		return nil, err
	}
	return getPremiumCheckoutTokenResponse{
		Token: *newToken,
	}, nil
}

type getManageTokenForPremiumCheckoutTokenRequest struct {
	Token string `json:"token"`
}

type getManageTokenForPremiumCheckoutTokenResponse struct {
	Token string `json:"token"`
}

func handleGetManageTokenForPremiumCheckoutToken(body []byte) (interface{}, error) {
	var req getManageTokenForPremiumCheckoutTokenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	userID, err := routetoken.ValidateTokenAndGetUserID(req.Token, routes.PremiumSubscriptionCheckoutKey)
	if err != nil {
		return nil, err
	}
	newToken, err := routes.MakeSubscriptionManagementToken(*userID)
	if err != nil {
		return nil, err
	}
	return getManageTokenForPremiumCheckoutTokenResponse{
		Token: *newToken,
	}, nil
}
