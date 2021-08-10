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
	SubscriptionCreationToken string `json:"subscription_creation_token"`
	IsYearlySubscription      bool   `json:"is_yearly_subscription"`
}

type createUserSubscriptionResponse struct {
	StripeSubscriptionID bgstripe.SubscriptionID `json:"stripe_subscription_id"`
	StripeClientSecret   string                  `json:"stripe_client_secret"`
	StripePaymentState   bgstripe.PaymentState   `json:"stripe_payment_state"`
}

func createUserSubscription(userID users.UserID, body []byte) (interface{}, error) {
	var req createUserSubscriptionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var stripeSubscriptionOutput *bgstripe.StripeCustomerSubscriptionOutput
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		if err := useraccounts.AddSubscriptionLevelForUser(tx, useraccounts.AddSubscriptionLevelForUserInput{
			UserID:            userID,
			SubscriptionLevel: useraccounts.SubscriptionLevelPremium,
			ShouldStartActive: false,
			// Create an expired subscription, the webhook will update this
			ExpirationTime: time.Now().Add(-10 * 24 * time.Hour),
		}); err != nil {
			return err
		}
		stripeSubscriptionOutput, err = bgstripe.CreateUnpaidStripeCustomerSubscriptionForUser(tx, userID, req.IsYearlySubscription)
		return err
	}); err != nil {
		return nil, err
	}
	return createUserSubscriptionResponse{
		StripeSubscriptionID: stripeSubscriptionOutput.SubscriptionID,
		StripeClientSecret:   stripeSubscriptionOutput.ClientSecret,
		StripePaymentState:   stripeSubscriptionOutput.PaymentState,
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

type deleteStripeSubscriptionForUserRequest struct{}

type deleteStripeSubscriptionForUserResponse struct{}

func deleteStripeSubscriptionForUser(userID users.UserID, body []byte) (interface{}, error) {
	var req deleteStripeSubscriptionForUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return bgstripe.CancelStripeSubscription(tx, userID)
	}); err != nil {
		return nil, err
	}
	return deleteStripeSubscriptionForUserResponse{}, nil
}

type updateStripeSubscriptionFrequencyForUserRequest struct {
	StripeSubscriptionID bgstripe.SubscriptionID `json:"stripe_subscription_id"`
	IsYearlySubscription bool                    `json:"is_yearly_subscription"`
}

type updateStripeSubscriptionFrequencyForUserResponse struct {
	Success bool `json:"success"`
}

func updateStripeSubscriptionFrequencyForUser(userID users.UserID, body []byte) (interface{}, error) {
	var req updateStripeSubscriptionFrequencyForUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var success bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		success, err = bgstripe.UpdateStripeSubscriptionChargeFrequency(tx, userID, req.StripeSubscriptionID, req.IsYearlySubscription)
		return err
	}); err != nil {
		return nil, err
	}
	return updateStripeSubscriptionFrequencyForUserResponse{
		Success: success,
	}, nil
}
