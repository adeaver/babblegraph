package initialdata

import (
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"encoding/json"
	"net/http"
)

type InitialFrontendData struct {
	Data string `json:"data"`
}

type initialData struct {
	StripePublicKey *string `json:"stripe_public_key,omitempty"`
}

type InitialFrontendDataOptions struct {
	IncludeStripePublicKey bool
}

func GetInitialFrontendData(r *http.Request, options InitialFrontendDataOptions) (*InitialFrontendData, error) {
	initialData := getInitialDataForRequest(r, options.IncludeStripePublicKey)
	bytes, err := json.Marshal(&initialData)
	if err != nil {
		return nil, err
	}
	return &InitialFrontendData{
		Data: string(bytes),
	}, nil
}

func getInitialDataForRequest(r *http.Request, includeStripePublicKey bool) initialData {
	var stripePublicKey *string
	if includeStripePublicKey {
		stripePublicKey = ptr.String(env.MustEnvironmentVariable("STRIPE_PUBLIC_KEY"))
	}
	return initialData{
		StripePublicKey: stripePublicKey,
	}
}
