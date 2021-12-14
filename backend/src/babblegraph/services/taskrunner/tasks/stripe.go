package tasks

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/util/ctx"
)

func ForceSyncStripeEvents(c ctx.LogContext) {
	bgstripe.ForceSyncStripeEvents(c)
}
