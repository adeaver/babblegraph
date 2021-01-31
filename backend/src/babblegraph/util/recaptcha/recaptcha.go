package recaptcha

import (
	"babblegraph/util/env"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const verificationURL = "https://www.google.com/recaptcha/api/siteverify"

type verificationResponse struct {
	Success                     bool     `json:"success"`
	ChallengeTimestampISOFormat string   `json:"challenge_ts"`
	Hostname                    string   `json:"hostname"`
	ErrorCodes                  []string `json:"error-codes,omitempty"`
}

func VerifyRecaptchaToken(token string) error {
	data := url.Values{}
	data.Set("secret", env.MustEnvironmentVariable("CAPTCHA_SECRET"))
	data.Set("response", token)
	resp, err := http.Post(verificationURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
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
