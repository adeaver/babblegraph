package routes

import (
	"babblegraph/model/users"
	"babblegraph/util/encrypt"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
)

func MakeUnsubscribeRouteForUserID(userID users.UserID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   "unsubscribe",
		Value: userID,
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("unsubscribe/%s", *token))), nil
}
