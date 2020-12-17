package tasks

import (
	"babblegraph/jobs/dailyemail"
	"babblegraph/util/env"
	"babblegraph/util/ses"
	"log"
)

func SendDailyEmail() error {
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
	log.Println("Sending emails")
	dailyEmailFn := dailyemail.GetDailyEmailJob(emailClient)
	if err := dailyEmailFn(); err != nil {
		return err
	}
	return nil
}
