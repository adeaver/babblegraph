package recaptcha

import (
	"babblegraph/util/env"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const verificationURL = "https://www.google.com/recaptcha/api/siteverify"

type verificationRequest struct {
	Secret string `json:"secret"`
	Token  string `json:"response"`
}

type verificationResponse struct {
	Success                     bool     `json:"success"`
	ChallengeTimestampISOFormat string   `json:"challenge_ts"`
	Hostname                    string   `json:"hostname"`
	ErrorCodes                  []string `json:"error-codes,omitempty"`
}

func VerifyRecaptchaToken(token string) error {
	reqBytes, err := json.Marshal(verificationRequest{
		Secret: env.MustEnvironmentVariable("CAPTCHA_SECRET"),
		Token:  token,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(verificationURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var unmarshalled verificationResponse
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return err
	}
	if !unmarshalled.Success {
		log.Println(fmt.Sprintf("Got error codes %+v", unmarshalled.ErrorCodes))
		return fmt.Errorf("reCAPTCHA verification failed")
	}
	return nil
}
