package scheduler

import (
	"babblegraph/util/async"
	"babblegraph/util/env"
	"fmt"
	"time"

	cron "github.com/robfig/cron/v3"
)

func StartScheduler(errs chan error) error {
	usEastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}
	c := cron.New(cron.WithLocation(usEastern))
	switch env.MustEnvironmentName() {
	case env.EnvironmentProd,
		env.EnvironmentStage:
		c.AddFunc("30 0 * * *", async.WithContext(errs, "archive-forgot-passwords", handleArchiveForgotPasswordAttempts).Func())
		c.AddFunc("30 2 * * *", async.WithContext(errs, "refetch", fetchNewLinksForSeedURLs).Func())
		c.AddFunc("30 3 * * *", async.WithContext(errs, "admin-2fa-cleanup", handleCleanUpAdminTwoFactorCodesAndAccessTokens).Func())
		c.AddFunc("30 4 * * *", async.WithContext(errs, "cleanup-newsletters", handleCleanupOldNewsletter).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "pending-verifications", handlePendingVerifications).Func())
		c.AddFunc("*/3 * * * *", async.WithContext(errs, "forgot-passwords", handlePendingForgotPasswordAttempts).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "send-2fa-codes", handleSendAdminTwoFactorAuthenticationCode).Func())
		c.AddFunc("*/10 * * * *", async.WithContext(errs, "sync-billing", handleSyncBilling).Func())
		c.AddFunc("*/10 * * * *", async.WithContext(errs, "user-account-notifications", handlePendingUserAccountNotificationRequests).Func())
	case env.EnvironmentLocal,
		env.EnvironmentLocalTestEmail:
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "cleanup-newsletters", handleCleanupOldNewsletter).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "admin-2fa-cleanup", handleCleanUpAdminTwoFactorCodesAndAccessTokens).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "pending-verifications", handlePendingVerifications).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "forgot-passwords", handlePendingForgotPasswordAttempts).Func())
		c.AddFunc("*/5 * * * *", async.WithContext(errs, "archive-forgot-passwords", handleArchiveForgotPasswordAttempts).Func())
		c.AddFunc("*/30 * * * *", async.WithContext(errs, "refetch", fetchNewLinksForSeedURLs).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "send-2fa-codes", handleSendAdminTwoFactorAuthenticationCode).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "sync-billing", handleSyncBilling).Func())
		c.AddFunc("*/1 * * * *", async.WithContext(errs, "user-account-notifications", handlePendingUserAccountNotificationRequests).Func())
	case env.EnvironmentLocalNoEmail:
		async.WithContext(errs, "sync-billing", handleSyncBilling).Func()()
		async.WithContext(errs, "refetch", fetchNewLinksForSeedURLs).Func()()
		async.WithContext(errs, "cleanup-newsletters", handleCleanupOldNewsletter).Func()()
	case env.EnvironmentTest:
		// no-op
	default:
		panic(fmt.Sprintf("Unsupported environment name: %s", env.MustEnvironmentName()))
	}
	c.Start()
	return nil
}
