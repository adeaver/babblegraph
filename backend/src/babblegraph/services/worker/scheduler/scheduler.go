package scheduler

import (
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/util/async"
	"babblegraph/util/env"
	"time"

	cron "github.com/robfig/cron/v3"
)

func StartScheduler(linkProcessor *linkprocessing.LinkProcessor, errs chan error) error {
	usEastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}
	c := cron.New(cron.WithLocation(usEastern))
	switch env.GetEnvironmentVariableOrDefault("ENV", "prod") {
	case "prod":
		c.AddFunc("30 0 * * *", async.WithContext(errs, "archive-forgot-passwords", handleArchiveForgotPasswordAttempts).Func())
		c.AddFunc("30 2 * * *", async.WithContext(errs, "refetch", makeRefetchSeedDomainJob(linkProcessor)).Func())
		c.AddFunc("30 3 * * *", async.WithContext(errs, "admin-2fa-cleanup", handleCleanUpAdminTwoFactorCodesAndAccessTokens).Func())
		c.AddFunc("30 4 * * *", async.WithContext(errs, "cleanup-newsletters", handleCleanupOldNewsletter).Func())
		c.AddFunc("30 12 * * *", async.WithContext(errs, "user-feedback", sendUserFeedbackEmails).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "pending-verifications", handlePendingVerifications).Func())
		c.AddFunc("*/3 * * * *", async.WithContext(errs, "forgot-passwords", handlePendingForgotPasswordAttempts).Func())
		c.AddFunc("*/5 * * * *", async.WithContext(errs, "account-notifications", handlePendingUserAccountNotificationRequests).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "send-2fa-codes", handleSendAdminTwoFactorAuthenticationCode).Func())
		c.AddFunc("*/10 * * * *", async.WithContext(errs, "sync-billing", handleSyncBilling).Func())
	case "local-test-emails",
		"local":
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "cleanup-newsletters", handleCleanupOldNewsletter).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "admin-2fa-cleanup", handleCleanUpAdminTwoFactorCodesAndAccessTokens).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "account-notifications", handlePendingUserAccountNotificationRequests).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "pending-verifications", handlePendingVerifications).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "forgot-passwords", handlePendingForgotPasswordAttempts).Func())
		c.AddFunc("*/5 * * * *", async.WithContext(errs, "archive-forgot-passwords", handleArchiveForgotPasswordAttempts).Func())
		c.AddFunc("*/30 * * * *", async.WithContext(errs, "refetch", makeRefetchSeedDomainJob(linkProcessor)).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "send-2fa-codes", handleSendAdminTwoFactorAuthenticationCode).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "sync-billing", handleSyncBilling).Func())
		async.WithContext(errs, "user-feedback", sendUserFeedbackEmails).Func()()
	case "local-no-email":
		async.WithContext(errs, "refetch", makeRefetchSeedDomainJob(linkProcessor)).Func()()
		async.WithContext(errs, "cleanup-newsletters", handleCleanupOldNewsletter).Func()()
	}
	c.Start()
	return nil
}

func makeRefetchSeedDomainJob(linkProcessor *linkprocessing.LinkProcessor) func(c async.Context) {
	return func(c async.Context) {
		if err := refetchSeedDomainsForNewContent(); err != nil {
			c.Errorf("Error refetching seed domains: %s", err.Error())
			return
		}
		c.Infof("Finished refetch. Reseeding link processor")
		linkProcessor.ReseedDomains()
	}
}
