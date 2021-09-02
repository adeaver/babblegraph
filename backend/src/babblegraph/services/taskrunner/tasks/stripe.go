package tasks

import "babblegraph/externalapis/bgstripe"

func ForceSyncStripeEvents() {
	bgstripe.ForceSyncStripeEvents()
}
