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

const (
	verificationURL       = "https://www.google.com/recaptcha/api/siteverify"
	verificationThreshold = 0.5
)

type verificationResponse struct {
	Success                     bool     `json:"success"`
	Score                       float64  `json:"score"`
	Action                      string   `json:"action"`
	ChallengeTimestampISOFormat string   `json:"challenge_ts"`
	Hostname                    string   `json:"hostname"`
	ErrorCodes                  []string `json:"error-codes,omitempty"`
}

func VerifyRecaptchaToken(action, token string) (bool, error) {
	data := url.Values{}
	data.Set("secret", env.MustEnvironmentVariable("CAPTCHA_SECRET"))
	data.Set("response", token)
	resp, err := http.Post(verificationURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var unmarshalled verificationResponse
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return false, err
	}
	if !unmarshalled.Success {
		log.Println(fmt.Sprintf("Got error codes %+v", unmarshalled.ErrorCodes))
		return false, fmt.Errorf("reCAPTCHA verification failed")
	}
	if unmarshalled.Action != action {
		log.Println(fmt.Sprintf("Action does not match"))
		return false, fmt.Errorf("reCAPTCHA action does not match action provided. Expected %s, got %s", action, unmarshalled.Action)
	}
	if unmarshalled.Score < verificationThreshold {
		return false, nil
	}
	return true, nil
}
