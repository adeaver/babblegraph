package scheduler

import (
	"babblegraph/util/ses"

	"github.com/getsentry/sentry-go"
)

func handleSendAdminTwoFactorAuthenticationCode(localSentryHub *sentry.Hub, emailClient *ses.Client) error {

}

func handleCleanUpAdminTwoFactorCodesAndAccessTokens(localSentryHub *sentry.Hub) error {

}
