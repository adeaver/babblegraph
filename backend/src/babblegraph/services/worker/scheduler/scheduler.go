package scheduler

import (
	"babblegraph/externalapis/bgstripe"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/util/env"
	"babblegraph/util/ses"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
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
		c.AddFunc("30 0 * * *", makeArchiveForgotPasswordAttemptsJob(errs))
		c.AddFunc("30 2 * * *", makeRefetchSeedDomainJob(linkProcessor, errs))
		c.AddFunc("30 3 * * *", makeHandleAdminExpirationJob(errs))
		c.AddFunc("30 12 * * *", makeUserFeedbackJob(errs))
		c.AddFunc("11 */3 * * *", makeExpireUserAccountsJob(errs))
		c.AddFunc("*/1 * * * *", makeVerificationJob(errs))
		c.AddFunc("*/3 * * * *", makeForgotPasswordJob(errs))
		c.AddFunc("*/5 * * * *", makeUserAccountNotificationsJob(errs))
		c.AddFunc("14 */1 * * *", makeSyncStripeEventsJob(errs))
		c.AddFunc("*/1 * * * *", makeHandleTwoFactorAuthenticationCode(errs))
	case "local-test-emails",
		"local":
		c.AddFunc("*/1 * * * *", makeHandleAdminExpirationJob(errs))
		c.AddFunc("*/1 * * * *", makeUserAccountNotificationsJob(errs))
		c.AddFunc("*/1 * * * *", makeVerificationJob(errs))
		c.AddFunc("*/1 * * * *", makeForgotPasswordJob(errs))
		c.AddFunc("*/3 * * * *", makeExpireUserAccountsJob(errs))
		c.AddFunc("*/5 * * * *", makeArchiveForgotPasswordAttemptsJob(errs))
		c.AddFunc("*/30 * * * *", makeRefetchSeedDomainJob(linkProcessor, errs))
		c.AddFunc("*/1 * * * *", makeSyncStripeEventsJob(errs))
		c.AddFunc("*/1 * * * *", makeHandleTwoFactorAuthenticationCode(errs))
		makeUserFeedbackJob(errs)()
	case "local-no-email":
		makeRefetchSeedDomainJob(linkProcessor, errs)()
	}
	c.Start()
	return nil
}

func makeRefetchSeedDomainJob(linkProcessor *linkprocessing.LinkProcessor, errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("reseed-job", fmt.Sprintf("reseed-job-%s-%d-%d", today.Month().String(), today.Day(), today.Year()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Refetch Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		if err := refetchSeedDomainsForNewContent(); err != nil {
			errs <- err
		}
		log.Println(fmt.Sprintf("Finished refetch. Reseeding link processor"))
		linkProcessor.ReseedDomains()
	}
}

func makeVerificationJob(errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("verification-job", fmt.Sprintf("verification-job-%s-%d-%d-%d-%d", today.Month().String(), today.Day(), today.Year(), today.Hour(), today.Minute()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Verification Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		log.Println("Initializing verification job...")
		emailClient := ses.NewClient(ses.NewClientInput{
			AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
			AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
			AWSRegion:          "us-east-1",
			FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
		})
		log.Println("Starting verification job...")
		if err := handlePendingVerifications(localHub, emailClient); err != nil {
			localHub.CaptureException(err)
			errs <- err
		}
	}
}

func makeUserFeedbackJob(errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("user-feedback-job", fmt.Sprintf("user-feedback-job-%s-%d-%d-%d-%d", today.Month().String(), today.Day(), today.Year(), today.Hour(), today.Minute()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("User Feedback Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		emailClient := ses.NewClient(ses.NewClientInput{
			AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
			AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
			AWSRegion:          "us-east-1",
			FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
		})
		if err := sendUserFeedbackEmails(localHub, emailClient); err != nil {
			localHub.CaptureException(err)
			errs <- err
		}
	}
}

func makeForgotPasswordJob(errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("forgot-password-job", fmt.Sprintf("forgot-password-job-%s-%d-%d-%d-%d", today.Month().String(), today.Day(), today.Year(), today.Hour(), today.Minute()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Forgotten Password Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		log.Println("Initializing forgot password job...")
		emailClient := ses.NewClient(ses.NewClientInput{
			AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
			AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
			AWSRegion:          "us-east-1",
			FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
		})
		log.Println("Starting forgot password job...")
		if err := handlePendingForgotPasswordAttempts(localHub, emailClient); err != nil {
			localHub.CaptureException(err)
			errs <- err
		}
	}
}

func makeArchiveForgotPasswordAttemptsJob(errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("archive-forgot-password-job", fmt.Sprintf("archive-forgot-password-job-%s-%d-%d-%d-%d", today.Month().String(), today.Day(), today.Year(), today.Hour(), today.Minute()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Forgot Password Archive Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		log.Println("Starting forgot password archive job...")
		if err := handleArchiveForgotPasswordAttempts(); err != nil {
			localHub.CaptureException(err)
			errs <- err
		}
	}
}

func makeUserAccountNotificationsJob(errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("user-accounts-notifications-job", fmt.Sprintf("user-accounts-notifications-job-%s-%d-%d-%d-%d", today.Month().String(), today.Day(), today.Year(), today.Hour(), today.Minute()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("User Accounts Notifications Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		log.Println("Starting user accounts notifications job...")
		emailClient := ses.NewClient(ses.NewClientInput{
			AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
			AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
			AWSRegion:          "us-east-1",
			FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
		})
		if err := handlePendingUserAccountNotificationRequests(localHub, emailClient); err != nil {
			localHub.CaptureException(err)
			errs <- err
		}
	}
}

func makeExpireUserAccountsJob(errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("user-accounts-expiration-job", fmt.Sprintf("user-accounts-expiration-job-%s-%d-%d-%d-%d", today.Month().String(), today.Day(), today.Year(), today.Hour(), today.Minute()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("User Accounts Expiration Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		log.Println("Starting user accounts expiration job...")
		if err := expireUserAccounts(); err != nil {
			localHub.CaptureException(err)
			errs <- err
		}
	}
}

func makeSyncStripeEventsJob(errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("sync-stripe-events-job", fmt.Sprintf("sync-stripe-events-job-%s-%d-%d-%d-%d", today.Month().String(), today.Day(), today.Year(), today.Hour(), today.Minute()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Sync Stripe Events Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		log.Println("Starting sync stripe events job...")
		bgstripe.ForceSyncStripeEvents()
	}
}

func makeHandleTwoFactorAuthenticationCode(errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("admin-2fa-codes", fmt.Sprintf("admin-2fa-codes-%s-%d-%d-%d-%d", today.Month().String(), today.Day(), today.Year(), today.Hour(), today.Minute()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Admin 2FA Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		log.Println("Starting Admin 2FA code job...")
		emailClient := ses.NewClient(ses.NewClientInput{
			AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
			AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
			AWSRegion:          "us-east-1",
			FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
		})
		if err := handleSendAdminTwoFactorAuthenticationCode(localHub, emailClient); err != nil {
			localHub.CaptureException(err)
			errs <- err
		}
	}
}

func makeHandleAdminExpirationJob(errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("admin-codes-expiration", fmt.Sprintf("admin-codes-expiration-%s-%d-%d-%d-%d", today.Month().String(), today.Day(), today.Year(), today.Hour(), today.Minute()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Admin Codes Expiration: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				localHub.CaptureException(err)
				errs <- err
			}
		}()
		log.Println("Starting admin codes expiration job...")
		if err := handleCleanUpAdminTwoFactorCodesAndAccessTokens(); err != nil {
			localHub.CaptureException(err)
			errs <- err
		}
	}
}
