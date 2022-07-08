package user

import (
	"babblegraph/model/billing"
	"babblegraph/services/web/router"
)

// TODO: Should add a timestamp to this (i.e. 3 per day, with 6 total)
const maxVerificationAttemptsForUser = 3

type signupError string

const (
	signupErrorInvalidEmailAddress signupError = "invalid-email"
	signupErrorIncorrectStatus     signupError = "invalid-account-status"
	signupErrorRateLimited         signupError = "rate-limited"
	signupErrorLowScore            signupError = "low-score"
)

func (s signupError) Ptr() *signupError {
	return &s
}

type signupUserRequest struct {
	EmailAddress string `json:"email_address"`
	CaptchaToken string `json:"captcha_token"`
}

type signupUserResponse struct {
	Success      bool         `json:"success"`
	ErrorMessage *signupError `json:"error_message,omitempty"`
}

func handleSignupUser(promotionCode *billing.PromotionCode, r *router.Request) (interface{}, error) {
	var req signupUserRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	return signupUserResponse{
		Success: false,
	}, nil
}
