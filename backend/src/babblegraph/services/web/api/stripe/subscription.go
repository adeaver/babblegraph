package stripe

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

type createUserSubscriptionRequest struct {
	SubscriptionType bgstripe.SubscriptionType `json:"subscription_type"`
}

type createUserSubscriptionResponse struct {
	Subscription *bgstripe.Subscription `json:"subscription,omitempty"`
}

func createUserSubscription(userID users.UserID, body []byte) (interface{}, error) {
	var req createUserSubscriptionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var subscription *bgstripe.Subscription
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		// Create an expired subscription, the webhook will update this
		if err := useraccounts.AddSubscriptionLevelForUser(tx, useraccounts.AddSubscriptionLevelForUserInput{
			UserID:            userID,
			SubscriptionLevel: useraccounts.SubscriptionLevelPremium,
			ShouldStartActive: false,
			ExpirationTime:    time.Now().Add(-10 * 24 * time.Hour),
		}); err != nil {
			return err
		}
		subscription, err = bgstripe.CreateSubscriptionForUser(tx, userID, req.SubscriptionType)
		return err
	}); err != nil {
		return nil, err
	}
	return createUserSubscriptionResponse{
		Subscription: subscription,
	}, nil
}

type getActiveSubscriptionForUserRequest struct{}

type getActiveSubscriptionForUserResponse struct {
	Subscription *bgstripe.Subscription `json:"subscription,omitempty"`
}

func getActiveSubscriptionForUser(userID users.UserID, body []byte) (interface{}, error) {
	var subscription *bgstripe.Subscription
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		subscription, err = bgstripe.LookupActiveSubscriptionForUser(tx, userID)
		return err
	}); err != nil {
		return nil, err
	}
	return getActiveSubscriptionForUserResponse{
		Subscription: subscription,
	}, nil
}

type getSubscriptionTrialInfoForUserRequest struct{}

type getSubscriptionTrialInfoForUserResponse struct {
	SubscriptionTrialInfo *bgstripe.SubscriptionTrialInfo `json:"subscription_trial_info"`
}

func getSubscriptionTrialInfoForUser(userID users.UserID, body []byte) (interface{}, error) {
	var subscriptionTrialInfo *bgstripe.SubscriptionTrialInfo
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		subscriptionTrialInfo, err = bgstripe.LookupSubscriptionTrialInfoForUser(tx, userID)
		return err
	}); err != nil {
		return nil, err
	}
	return getSubscriptionTrialInfoForUserResponse{
		SubscriptionTrialInfo: subscriptionTrialInfo,
	}, nil
}

type updateStripeSubscriptionForUserRequest struct {
	Options bgstripe.UpdateSubscriptionOptions `json:"options"`
}

type updateStripeSubscriptionForUserResponse struct {
	Subscription *bgstripe.Subscription `json:"subscription"`
}

func updateStripeSubscriptionForUser(userID users.UserID, body []byte) (interface{}, error) {
	var req updateStripeSubscriptionForUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return bgstripe.UpdateSubscription(tx, userID, req.Options)
	}); err != nil {
		return nil, err
	}
	// This is intentionally done in two separate transactions
	// so that we don't roll back database updates if the second
	// request to Stripe fails
	var subscription *bgstripe.Subscription
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		subscription, err = bgstripe.LookupActiveSubscriptionForUser(tx, userID)
		return err
	}); err != nil {
		return nil, err
	}
	return updateStripeSubscriptionForUserResponse{
		Subscription: subscription,
	}, nil
}
