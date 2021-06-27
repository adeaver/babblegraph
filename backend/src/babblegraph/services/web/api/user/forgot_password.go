package user

import (
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"babblegraph/util/recaptcha"
	"encoding/json"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

type requestPasswordResetLinkRequest struct {
	EmailAddress string `json:"email_address"`
	CaptchaToken string `json:"captcha_token"`
}

type requestPasswordResetLinkResponse struct {
	Success bool `json:"success"`
}

func requestPasswordResetLink(body []byte) (interface{}, error) {
	var req requestPasswordResetLinkRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	isValid, err := recaptcha.VerifyRecaptchaToken("forgotpassword", req.CaptchaToken)
	switch {
	case err != nil:
		sentry.CaptureException(err)
		return nil, err
	case !isValid:
		return requestPasswordResetLinkResponse{
			Success: false,
		}, nil
	default:
		// no-op
		log.Println("Successfully cleared captcha")
	}
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	shouldSendSentryOnError := true
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		user, err := users.LookupUserByEmailAddress(tx, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case user == nil:
			log.Println(fmt.Sprintf("No user found for email address %s, continuing", formattedEmailAddress))
			return nil
		case user.Status != users.UserStatusVerified:
			log.Println(fmt.Sprintf("User found for email address %s does not have verified status, continuing", formattedEmailAddress))
			return nil
		}
		alreadyHasAccount, err := useraccounts.DoesUserAlreadyHaveAccount(tx, user.ID)
		switch {
		case err != nil:
			return err
		case !alreadyHasAccount:
			log.Println(fmt.Sprintf("User found for email address %s does not have an acccount, continuing", formattedEmailAddress))
			return nil
		}
		hasTooManyAttempts, err := useraccounts.AddForgotPasswordAttemptForUserID(tx, user.ID)
		switch {
		case err != nil:
			return err
		case hasTooManyAttempts:
			shouldSendSentryOnError = false
			return fmt.Errorf("User found for email address %s has too many forgot password attempts, continuing", formattedEmailAddress)
		}
		return nil
	}); err != nil {
		if shouldSendSentryOnError {
			sentry.CaptureException(fmt.Errorf("Error adding new forgot password request for user %s: %s", formattedEmailAddress, err.Error()))
		}
		log.Println(fmt.Sprintf("Error adding new forgot password request for user %s: %s", formattedEmailAddress, err.Error()))
		return requestPasswordResetLinkResponse{
			Success: false,
		}, nil
	}
	return requestPasswordResetLinkResponse{
		Success: true,
	}, nil
}
