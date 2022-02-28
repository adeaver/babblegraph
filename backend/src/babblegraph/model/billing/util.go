package billing

import "github.com/google/uuid"

func NewPremiumNewsletterSubscriptionID() PremiumNewsletterSubscriptionID {
	return PremiumNewsletterSubscriptionID(uuid.New().String())
}
