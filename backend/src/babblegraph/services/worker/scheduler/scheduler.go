package scheduler

import (
	"babblegraph/jobs/dailyemail"
	"babblegraph/util/env"
	"babblegraph/util/ses"
	"log"
	"runtime/debug"
	"time"

	cron "github.com/robfig/cron/v3"
)

func StartScheduler(errs chan error) error {
	usEastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}
	c := cron.New(cron.WithLocation(usEastern))
	switch env.GetEnvironmentVariableOrDefault("ENV", "prod") {
	case "prod":
		c.AddFunc("30 2 * * *", makeRefetchSeedDomainJob(errs))
		c.AddFunc("30 5 * * *", makeEmailJob(errs))
	case "local":
		makeEmailJob(errs)()
		makeRefetchSeedDomainJob(errs)()
	case "local-no-email":
		// no-op
		makeRefetchSeedDomainJob(errs)()
	}
	c.Start()
	return nil
}

func makeRefetchSeedDomainJob(errs chan error) func() {
	return func() {
		defer func() {
			x := recover()
			if err, ok := x.(error); ok {
				errs <- err
				debug.PrintStack()
			}
		}()
		if err := RefetchSeedDomainsForNewContent(); err != nil {
			errs <- err
		}
	}
}

func makeEmailJob(errs chan error) func() {
	return func() {
		defer func() {
			x := recover()
			if err, ok := x.(error); ok {
				errs <- err
				debug.PrintStack()
			}
		}()
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
