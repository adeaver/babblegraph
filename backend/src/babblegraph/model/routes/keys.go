package routes

import (
	"babblegraph/model/users"
	"babblegraph/util/encrypt"
)

type RouteEncryptionKey string

const (
	SubscriptionManagementRouteEncryptionKey RouteEncryptionKey = "subscription-management"
	UnsubscribeRouteEncryptionKey            RouteEncryptionKey = "unsubscribe"
	EmailOpenedKey                           RouteEncryptionKey = "email-opened"
	UserVerificationKey                      RouteEncryptionKey = "user-verification"
	WordReinforcementKey                     RouteEncryptionKey = "word-reinforcement"
)

func (r RouteEncryptionKey) Str() string {
	return string(r)
}

func MakeWordReinforcementToken(userID users.UserID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   WordReinforcementKey.Str(),
		Value: userID,
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func MakeSubscriptionManagementToken(userID users.UserID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   SubscriptionManagementRouteEncryptionKey.Str(),
		Value: userID,
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
