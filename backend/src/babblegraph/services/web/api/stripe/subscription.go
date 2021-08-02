package stripe

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type getOrCreateUserSubscriptionRequest struct {
	SubscriptionCreationToken string `json:"subscription_creation_token"`
	IsYearlySubscription      bool   `json:"is_yearly_subscription"`
}

type getOrCreateUserSubscriptionResponse struct {
	StripeSubscriptionID bgstripe.SubscriptionID `json:"stripe_subscription_id"`
	StripeClientSecret   string                  `json:"stripe_client_secret"`
	StripePaymentState   bgstripe.PaymentState   `json:"stripe_payment_state"`
}

func getOrCreateUserSubscription(userID users.UserID, body []byte) (interface{}, error) {
	var req *getOrCreateUserSubscriptionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var stripeSubscriptionOutput *bgstripe.StripeCustomerSubscriptionOutput
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		// This method also creates the user subscription if applicable
		stripeSubscriptionOutput, err = bgstripe.GetOrCreateUnpaidStripeCustomerSubscriptionForUser(tx, userID, req.IsYearlySubscription)
		return err
	}); err != nil {
		return nil, err
	}
	return getOrCreateUserSubscriptionResponse{
		StripeSubscriptionID: stripeSubscriptionOutput.SubscriptionID,
		StripeClientSecret:   stripeSubscriptionOutput.ClientSecret,
		StripePaymentState:   stripeSubscriptionOutput.PaymentState,
	}, nil
}

type getUserNonTerminatedStripeSubscriptionRequest struct {
	SubscriptionCreationToken string `json:"subscription_creation_token"`
}

type getUserNonTerminatedStripeSubscriptionResponse struct {
	IsYearlySubscription *bool                    `json:"is_yearly_subscription,omitempty"`
	StripeSubscriptionID *bgstripe.SubscriptionID `json:"stripe_subscription_id,omitempty"`
	StripeClientSecret   *string                  `json:"stripe_client_secret,omitempty"`
	StripePaymentState   *bgstripe.PaymentState   `json:"stripe_payment_state,omitempty"`
}

func getUserNonTerminatedStripeSubscription(userID users.UserID, body []byte) (interface{}, error) {
	var stripeSubscriptionOutput *bgstripe.StripeCustomerSubscriptionOutput
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		stripeSubscriptionOutput, err = bgstripe.LookupNonterminatedStripeSubscriptionForUser(tx, userID)
		return err
	}); err != nil {
		return nil, err
	}
	if stripeSubscriptionOutput == nil {
		return getUserNonTerminatedStripeSubscriptionResponse{}, nil
	}
	return getUserNonTerminatedStripeSubscriptionResponse{
		StripeSubscriptionID: &stripeSubscriptionOutput.SubscriptionID,
		StripeClientSecret:   &stripeSubscriptionOutput.ClientSecret,
		StripePaymentState:   &stripeSubscriptionOutput.PaymentState,
		IsYearlySubscription: &stripeSubscriptionOutput.IsYearlySubscription,
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
