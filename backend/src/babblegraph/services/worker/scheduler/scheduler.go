package scheduler

import (
	"babblegraph/jobs/dailyemail"
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
		c.AddFunc("30 2 * * *", makeRefetchSeedDomainJob(linkProcessor, errs))
		c.AddFunc("30 5 * * *", makeEmailJob(errs))
		c.AddFunc("*/1 * * * *", makeVerificationJob(errs))
	case "local":
		c.AddFunc("*/1 * * * *", makeVerificationJob(errs))
		c.AddFunc("*/30 * * * *", makeRefetchSeedDomainJob(linkProcessor, errs))
		makeEmailJob(errs)()
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
				err := fmt.Errorf("Refetch Panic: %s: %d: %v\n", fn, line, x)
				localHub.CaptureException(err)
				debug.PrintStack()
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

func makeEmailJob(errs chan error) func() {
	return func() {
		today := time.Now()
		localHub := sentry.CurrentHub().Clone()
		localHub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("email-job", fmt.Sprintf("email-job-%s-%d-%d", today.Month().String(), today.Day(), today.Year()))
		})
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(1)
				err := fmt.Errorf("Email Panic: %s: %d: %v\n", fn, line, x)
				localHub.CaptureException(err)
				debug.PrintStack()
				errs <- err
			}
		}()
		log.Println("Initializing email client...")
		emailClient := ses.NewClient(ses.NewClientInput{
			AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
			AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
			AWSRegion:          "us-east-1",
			FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
		})
		log.Println("Starting email job...")
		dailyEmailFn := dailyemail.GetDailyEmailJob(localHub, emailClient)
		if err := dailyEmailFn(); err != nil {
			localHub.CaptureException(err)
			errs <- err
		}
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
				err := fmt.Errorf("Verification Panic: %s: %d: %v\n", fn, line, x)
				localHub.CaptureException(err)
				debug.PrintStack()
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
