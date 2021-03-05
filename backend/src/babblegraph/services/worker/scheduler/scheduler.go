package scheduler

import (
	"babblegraph/jobs/dailyemail"
	"babblegraph/services/worker/linkprocessing"
	"babblegraph/util/env"
	"babblegraph/util/ses"
	"fmt"
	"log"
	"runtime/debug"
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
		c.AddFunc("30 2 * * *", makeRefetchSeedDomainJob(linkProcessor, errs))
		c.AddFunc("30 5 * * *", makeEmailJob(errs))
		c.AddFunc("*/3 * * * *", makeVerificationJob(errs))
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
		defer func() {
			x := recover()
			if err, ok := x.(error); ok {
				errs <- err
				debug.PrintStack()
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
		defer func() {
			x := recover()
			if err, ok := x.(error); ok {
				log.Println(fmt.Sprintf("Got error on email job: %s", err.Error()))
				errs <- err
				debug.PrintStack()
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
		dailyEmailFn := dailyemail.GetDailyEmailJob(emailClient)
		if err := dailyEmailFn(); err != nil {
			errs <- err
		}
	}
}

func makeVerificationJob(errs chan error) func() {
	return func() {
		defer func() {
			x := recover()
			if err, ok := x.(error); ok {
				errs <- err
				debug.PrintStack()
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
		if err := handlePendingVerifications(emailClient); err != nil {
			errs <- err
		}
	}
}
