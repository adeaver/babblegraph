package routes

import (
	"babblegraph/model/users"
	"babblegraph/util/env"
	"fmt"
)

const RedirectKeyParameter string = "d"

type LoginRedirectKey string

const (
	LoginRedirectKeySubscriptionManagement LoginRedirectKey = "sbmgmt"
	LoginRedirectKeyVocabulary             LoginRedirectKey = "vb"
	LoginRedirectKeyContentTopics          LoginRedirectKey = "cts"

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
	case LoginRedirectKeyContentTopics.Str():
		return LoginRedirectKeyContentTopics
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
	case LoginRedirectKeyContentTopics:
		return MakeSetTopicsLink(userID)
	default:
		return nil, fmt.Errorf("unimplemented")
	}
}

func MakeLoginLinkWithContentTopicsRedirect() string {
	return makeLoginLinkForLoginRedirectKey(LoginRedirectKeyContentTopics)
}

func MakeLoginLinkWithReinforcementRedirect() string {
	return makeLoginLinkForLoginRedirectKey(LoginRedirectKeyVocabulary)
}

func makeLoginLinkForLoginRedirectKey(key LoginRedirectKey) string {
	return env.GetAbsoluteURLForEnvironment(fmt.Sprintf("login?%s=%s", RedirectKeyParameter, key.Str()))
}
