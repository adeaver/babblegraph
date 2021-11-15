package initialdata

import (
	"babblegraph/util/env"
	"encoding/json"
	"net/http"
)

type InitialFrontendData struct {
	Data string `json:"data"`
}

type initialData struct {
	StripePublicKey string `json:"stripe_public_key,omitempty"`
}

func GetInitialFrontendData(r *http.Request) (*InitialFrontendData, error) {
	initialData := getInitialDataForRequest(r)
	bytes, err := json.Marshal(&initialData)
	if err != nil {
		return nil, err
	}
	return &InitialFrontendData{
		Data: string(bytes),
	}, nil
}

func getInitialDataForRequest(r *http.Request) initialData {
	return initialData{
		StripePublicKey: env.MustEnvironmentVariable("STRIPE_PUBLIC_KEY"),
	}
}
