package routes

import (
	"babblegraph/model/admin"
	"babblegraph/model/users"
	"babblegraph/util/encrypt"
	"fmt"
)

type RouteEncryptionKey string

const (
	SubscriptionManagementRouteEncryptionKey RouteEncryptionKey = "subscription-management"
	UnsubscribeRouteEncryptionKey            RouteEncryptionKey = "unsubscribe"
	EmailOpenedKey                           RouteEncryptionKey = "email-opened"
	UserVerificationKey                      RouteEncryptionKey = "user-verification"
	WordReinforcementKey                     RouteEncryptionKey = "word-reinforcement"
	OnboardingKey                            RouteEncryptionKey = "onboarding"

	PremiumSubscriptionCheckoutKey RouteEncryptionKey = "premium-subscription-checkout"
	CreateUserKey                  RouteEncryptionKey = "create-user"
	ForgotPasswordKey              RouteEncryptionKey = "forgot-password"

	ArticleLinkKeyForUserDocumentID   RouteEncryptionKey = "article-link-user-document"
	PaywallReportKeyForUserDocumentID RouteEncryptionKey = "paywall-report-user-document"

	AdminRegistrationKey RouteEncryptionKey = "admin-registration"

	/*
			   A mistake to learn from:
			   These are deprecated because Babblegraph used to encrypt
			   a JSON object that contained UserID, EmailRecordID, and ArticleLink
			   This generated 300+ character urls. Most clients and browsers can handle this
			   but there are some that can't

		       These remain here for backwards compatibility
	*/
	ArticleLinkKeyDEPRECATED   RouteEncryptionKey = "article-link"
	PaywallReportKeyDEPRECATED RouteEncryptionKey = "paywall-report"
)

func (r RouteEncryptionKey) Str() string {
	return string(r)
}

func EncryptUserIDWithKey(userID users.UserID, key RouteEncryptionKey) (*string, error) {
	switch key {
	case SubscriptionManagementRouteEncryptionKey,
		UnsubscribeRouteEncryptionKey,
		UserVerificationKey,
		WordReinforcementKey,
		PremiumSubscriptionCheckoutKey,
		CreateUserKey,
		ForgotPasswordKey,
		ArticleLinkKeyForUserDocumentID,
		PaywallReportKeyForUserDocumentID:
		return encrypt.GetToken(encrypt.TokenPair{
			Key:   key.Str(),
			Value: userID,
		})
	default:
		return nil, fmt.Errorf("Invalid key type: %s", key.Str())
	}
}

func MakeWordReinforcementToken(userID users.UserID) (*string, error) {
	return encrypt.GetToken(encrypt.TokenPair{
		Key:   WordReinforcementKey.Str(),
		Value: userID,
	})
}

func MakeSubscriptionManagementToken(userID users.UserID) (*string, error) {
	return encrypt.GetToken(encrypt.TokenPair{
		Key:   SubscriptionManagementRouteEncryptionKey.Str(),
		Value: userID,
	})
}

func MakeCreateUserToken(userID users.UserID) (*string, error) {
	return encrypt.GetToken(encrypt.TokenPair{
		Key:   CreateUserKey.Str(),
		Value: userID,
	})
}

func MakePremiumSubscriptionCheckoutToken(userID users.UserID) (*string, error) {
	return encrypt.GetToken(encrypt.TokenPair{
		Key:   PremiumSubscriptionCheckoutKey.Str(),
		Value: userID,
	})
}

func MakeAdminRegistrationToken(adminID admin.ID) (*string, error) {
	return encrypt.GetToken(encrypt.TokenPair{
		Key:   AdminRegistrationKey.Str(),
		Value: adminID,
	})
}
