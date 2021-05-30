package routes

import (
	"babblegraph/model/email"
	"babblegraph/model/users"
	"babblegraph/util/encrypt"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
)

func MustGetHomePageURL() string {
	return env.GetAbsoluteURLForEnvironment("")
}

func MakeSubscriptionManagementRouteForUserID(userID users.UserID) (*string, error) {
	token, err := MakeSubscriptionManagementToken(userID)
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("manage/%s", *token))), nil
}

func MakeSetTopicsLink(userID users.UserID) (*string, error) {
	managementLink, err := MakeSubscriptionManagementRouteForUserID(userID)
	if err != nil {
		return nil, err
	}
	return ptr.String(fmt.Sprintf("%s/interests", *managementLink)), nil
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

func MakeLogoURLForEmailRecordID(emailRecordID email.ID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   EmailOpenedKey.Str(),
		Value: emailRecordID,
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("dist/%s/logo.png", *token))), nil
}

func MakeUserVerificationLink(userID users.UserID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   UserVerificationKey.Str(),
		Value: userID,
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("verify/%s", *token))), nil
}

func MakeWordReinforcementLink(userID users.UserID) (*string, error) {
	token, err := MakeWordReinforcementToken(userID)
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("manage/%s/vocabulary", *token))), nil
}

func MakeArticleLink(userID users.UserID, emailRecordID email.ID, u string) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key: ArticleLinkKey.Str(),
		Value: ArticleLinkBody{
			UserID:        userID,
			EmailRecordID: emailRecordID,
			URL:           u,
		},
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("article/%s", *token))), nil
}

func MakePaywallReportLink(userID users.UserID, emailRecordID email.ID, u string) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key: PaywallReportKey.Str(),
		Value: PaywallReportBody{
			UserID:        userID,
			EmailRecordID: emailRecordID,
			URL:           u,
		},
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("paywall-report/%s", *token))), nil
}
