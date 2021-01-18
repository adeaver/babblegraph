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

func handleSignupUser(body []byte) (interface{}, error) {
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
		if err := users.InsertNewUnverifiedUser(tx, formattedEmailAddress); err != nil {
			return err
		}
		// This is necessary because the above ignores duplicates if the user is already
		// present in the database. We must now query the database for the user we just
		// created or the already existing user.
		user, err := users.LookupUserByEmailAddress(tx, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case user == nil:
			return fmt.Errorf("No error, but no user returned")
		case user.Status != users.UserStatusUnverified:
			sErr = signupErrorIncorrectStatus.Ptr()
			return fmt.Errorf("Invalid account status")
		}
		numAttempts, err := userverificationattempt.GetNumberOfFulfilledVerificationAttemptsForUser(tx, user.ID)
		switch {
		case err != nil:
			return err
		case numAttempts != nil && *numAttempts >= maxVerificationAttemptsForUser:
			sErr = signupErrorRateLimited.Ptr()
			return fmt.Errorf("Rate limited")
		}
		return userverificationattempt.InsertVerificationAttemptForUser(tx, user.ID)
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
