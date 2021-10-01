package tasks

import (
	"babblegraph/jobs/dailyemail"
	"babblegraph/util/env"
	"babblegraph/util/ses"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func SendDailyEmail() error {
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
	log.Println("Sending emails")
	today := time.Now()
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("email-job", fmt.Sprintf("email-job-%s-%d-%d", today.Month().String(), today.Day(), today.Year()))
	})
	dailyEmailFn := dailyemail.GetDailyEmailJob(localHub, emailClient)
	if err := dailyEmailFn(); err != nil {
		return err
	}
	return nil
}

func SendDailyEmailForEmailAddresses(emailAddresses []string) error {
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
	log.Println("Sending emails")
	today := time.Now()
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("email-job", fmt.Sprintf("email-job-%s-%d-%d", today.Month().String(), today.Day(), today.Year()))
	})
	dailyEmailFn := dailyemail.SendDailyEmailToUsersByEmailAddress(localHub, emailClient)
	if err := dailyEmailFn(emailAddresses); err != nil {
		return err
	}
	return nil
}
