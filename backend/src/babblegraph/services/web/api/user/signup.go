package user

import (
	"babblegraph/model/users"
	"babblegraph/model/userverificationattempt"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// TODO: Should add a timestamp to this (i.e. 3 per day, with 6 total)
const maxVerificationAttemptsForUser = 3

type signupError string

const (
	signupErrorInvalidEmailAddress signupError = "invalid-email"
	signupErrorIncorrectStatus     signupError = "invalid-account-status"
	signupErrorRateLimited         signupError = "rate-limited"
)

func (s signupError) Ptr() *signupError {
	return &s
}

type signupUserRequest struct {
	EmailAddress string `json:"email_address"`
}

type signupUserResponse struct {
	Success      bool         `json:"success"`
	ErrorMessage *signupError `json:"error_message,omitempty"`
}

func handleSignupUserRequest(body []byte) (interface{}, error) {
	var req signupUserRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	// TODO: insert CAPTCHA verification
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	if err := email.ValidateEmailAddress(formattedEmailAddress); err != nil {
		log.Println(fmt.Sprintf("Error validating email address %s: %s", formattedEmailAddress, err.Error()))
		return signupUserResponse{
			ErrorMessage: signupErrorInvalidEmailAddress.Ptr(),
		}, nil
	}
	var sErr *signupError
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		userID, status, err := users.InsertNewUnverifiedUser(tx, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case userID == nil || status == nil:
			return fmt.Errorf("No error, but no status or ID returned")
		case *status != users.UserStatusUnverified:
			sErr = signupErrorIncorrectStatus.Ptr()
			return fmt.Errorf("Invalid account status")
		}
		numAttempts, err := userverificationattempt.GetNumberOfFulfilledVerificationAttemptsForUser(tx, *userID)
		switch {
		case err != nil:
			return err
		case numAttempts != nil && *numAttempts < maxVerificationAttemptsForUser:
			sErr = signupErrorRateLimited.Ptr()
			return fmt.Errorf("Rate limited")
		}
		return userverificationattempt.InsertVerificationAttemptForUser(tx, *userID)
	}); err != nil {
		log.Println(fmt.Sprintf("Got error on email address %s: %s", formattedEmailAddress, err.Error()))
		if sErr != nil {
			return signupUserResponse{
				ErrorMessage: sErr,
			}, nil
		}
		return nil, err
	}
	return signupUserResponse{
		Success: true,
	}, nil
}
