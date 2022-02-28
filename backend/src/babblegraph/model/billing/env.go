package billing

import (
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
)

const (
	PremiumNewsletterSubscriptionStripeProductIDProd = "price_1KWWQoJscBSiX47S6Jz4vd8O"
	PremiumNewsletterSubscriptionStripeProductIDTest = "price_1KWWReJscBSiX47SQiyCPb5X"
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
