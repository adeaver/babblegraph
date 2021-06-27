package routes

import (
	"babblegraph/model/email"
	"babblegraph/model/useraccounts"
	"babblegraph/model/userdocuments"
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

func MakeUserCreationLink(userID users.UserID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   CreateUserKey.Str(),
		Value: userID,
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("signup/%s", *token))), nil
}

func MakeArticleLink(userDocumentID userdocuments.UserDocumentID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   ArticleLinkKeyForUserDocumentID.Str(),
		Value: userDocumentID,
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("article/%s", *token))), nil
}

func MakePaywallReportLink(userDocumentID userdocuments.UserDocumentID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   PaywallReportKeyForUserDocumentID.Str(),
		Value: userDocumentID,
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("paywall-report/%s", *token))), nil
}

func MakeForgotPasswordLink(forgotPasswordAttemptID useraccounts.ForgotPasswordAttemptID) (*string, error) {
	token, err := encrypt.GetToken(encrypt.TokenPair{
		Key:   ForgotPasswordKey.Str(),
		Value: forgotPasswordAttemptID,
	})
	if err != nil {
		return nil, err
	}
	return ptr.String(env.GetAbsoluteURLForEnvironment(fmt.Sprintf("password-reset/%s", *token))), nil
}
