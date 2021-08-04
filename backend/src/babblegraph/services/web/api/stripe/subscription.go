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

type getUserNonTerminatedStripeSubscriptionRequest struct {
	SubscriptionCreationToken string `json:"subscription_creation_token"`
}

type getUserNonTerminatedStripeSubscriptionResponse struct {
	IsEligibleForTrial   bool                     `json:"is_eligible_for_trial"`
	IsYearlySubscription *bool                    `json:"is_yearly_subscription,omitempty"`
	StripeSubscriptionID *bgstripe.SubscriptionID `json:"stripe_subscription_id,omitempty"`
	StripeClientSecret   *string                  `json:"stripe_client_secret,omitempty"`
	StripePaymentState   *bgstripe.PaymentState   `json:"stripe_payment_state,omitempty"`
}

func getUserNonTerminatedStripeSubscription(userID users.UserID, body []byte) (interface{}, error) {
	var stripeSubscriptionOutput *bgstripe.StripeCustomerSubscriptionOutput
	var isEligibleForTrial bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		stripeSubscriptionOutput, isEligibleForTrial, err = bgstripe.LookupNonterminatedStripeSubscriptionForUser(tx, userID)
		return err
	}); err != nil {
		return nil, err
	}
	if stripeSubscriptionOutput == nil {
		return getUserNonTerminatedStripeSubscriptionResponse{
			IsEligibleForTrial: isEligibleForTrial,
		}, nil
	}
	return getUserNonTerminatedStripeSubscriptionResponse{
		StripeSubscriptionID: &stripeSubscriptionOutput.SubscriptionID,
		StripeClientSecret:   &stripeSubscriptionOutput.ClientSecret,
		StripePaymentState:   &stripeSubscriptionOutput.PaymentState,
		IsYearlySubscription: &stripeSubscriptionOutput.IsYearlySubscription,
		IsEligibleForTrial:   isEligibleForTrial,
	}, nil
}

type deleteStripeSubscriptionForUserRequest struct {
	StripeSubscriptionID bgstripe.SubscriptionID `json:"stripe_subscription_id"`
}

type deleteStripeSubscriptionForUserResponse struct {
	Success bool `json:"success"`
}

func deleteStripeSubscriptionForUser(userID users.UserID, body []byte) (interface{}, error) {
	var req deleteStripeSubscriptionForUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var success bool
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		success, err = bgstripe.CancelStripeSubscription(tx, userID, req.StripeSubscriptionID)
		return err
	}); err != nil {
		return nil, err
	}
	return deleteStripeSubscriptionForUserResponse{
		Success: success,
	}, nil
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
