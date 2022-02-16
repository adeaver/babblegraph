package billing

import (
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
)

const (
	PremiumNewsletterSubscriptionStripeProductIDProd = "price_1JIMqNJscBSiX47SxOGRUX1p"
	PremiumNewsletterSubscriptionStripeProductIDTest = "price_1JP8vOJscBSiX47S7sUI1z49"
)

func getStripeProductIDForEnvironment() (*string, error) {
	currentEnv := env.MustEnvironmentName()
	switch currentEnv {
	case env.EnvironmentProd:
		return ptr.String(PremiumNewsletterSubscriptionStripeProductIDProd), nil
	case env.EnvironmentStage,
		env.EnvironmentLocal,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentLocalTestEmail:
		return ptr.String(PremiumNewsletterSubscriptionStripeProductIDTest), nil
	default:
		return nil, fmt.Errorf("unsupported environment: %s", currentEnv)
	}
}
