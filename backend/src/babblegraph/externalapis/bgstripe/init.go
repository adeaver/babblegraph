package bgstripe

import (
	"babblegraph/util/env"

	"github.com/stripe/stripe-go"
)

func InitStripe() {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
}
