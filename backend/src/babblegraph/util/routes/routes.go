package routes

import (
	"babblegraph/model/users"
	"babblegraph/util/env"
	"fmt"
)

func MakeUnsubscribeRouteForUserID(userID users.UserID) string {
	return env.GetAbsoluteURLForEnvironment(fmt.Sprintf("unsubscribe/%s", userID))
}
