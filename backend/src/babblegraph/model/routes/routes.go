package routes

import (
	"babblegraph/model/users"
	"babblegraph/util/encrypt"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
)

func MakeSubscriptionManagementRouteForUserID(userID users.UserID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   SubscriptionManagementRouteEncryptionKey.Str(),
		Value: userID,
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("manage/%s", *token))), nil
}

func MakeUnsubscribeRouteForUserID(userID users.UserID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   UnsubscribeRouteEncryptionKey.Str(),
		Value: userID,
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("unsubscribe/%s", *token))), nil
}
