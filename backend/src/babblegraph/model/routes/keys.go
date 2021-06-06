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
	ArticleLinkKeyForUserDocumentID          RouteEncryptionKey = "article-link-user-document"
	PaywallReportKeyForUserDocumentID        RouteEncryptionKey = "paywall-report-user-document"

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
