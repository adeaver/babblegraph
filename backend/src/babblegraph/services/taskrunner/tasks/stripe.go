package tasks

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/util/bglog"
)

func ForceSyncStripeEvents(c bglog.LogContext) {
	bgstripe.ForceSyncStripeEvents(c)
}
