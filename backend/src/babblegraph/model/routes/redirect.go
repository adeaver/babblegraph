package routes

import (
	"babblegraph/model/users"
	"fmt"
)

type LoginRedirectKey string

const (
	LoginRedirectKeySubscriptionManagement LoginRedirectKey = "sbmgmt"
	LoginRedirectKeyVocabulary             LoginRedirectKey = "vb"

	LoginRedirectKeyDefault = LoginRedirectKeySubscriptionManagement
)

func (l LoginRedirectKey) Str() string {
	return string(l)
}

func GetLoginRedirectKeyOrDefault(loginKey string) LoginRedirectKey {
	switch loginKey {
	case LoginRedirectKeySubscriptionManagement.Str():
		return LoginRedirectKeySubscriptionManagement
	case LoginRedirectKeyVocabulary.Str():
		return LoginRedirectKeyVocabulary
	default:
		return LoginRedirectKeyDefault
	}
}

func GetLoginRedirectRouteForKeyAndUser(loginRedirectKey LoginRedirectKey, userID users.UserID) (*string, error) {
	switch loginRedirectKey {
	case LoginRedirectKeySubscriptionManagement:
		return MakeSubscriptionManagementRouteForUserID(userID)
	case LoginRedirectKeyVocabulary:
		return MakeWordReinforcementLink(userID)
	default:
		return nil, fmt.Errorf("unimplemented")
	}
}
