package main

import (
	"babblegraph/services/taskrunner/tasks"
	"babblegraph/util/database"
	"babblegraph/util/elastic"
	"babblegraph/util/env"
	"babblegraph/util/ses"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	if err := setupDatabases(); err != nil {
		log.Fatal(err.Error())
	}
	taskName := flag.String("task", "none", "Name of task to run [daily-email] [create-user]")
	userEmail := flag.String("user-email", "none", "Email address of user to create")
	flag.Parse()
	if taskName == nil {
		log.Fatal("No task specified")
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         env.MustEnvironmentVariable("SENTRY_DSN"),
		Environment: env.MustEnvironmentName().Str(),
	}); err != nil {
		log.Fatal(err.Error())
	}
	defer sentry.Flush(2 * time.Second)
	switch *taskName {
	case "daily-email":
		if err := tasks.SendDailyEmail(); err != nil {
			log.Fatal(err.Error())
		}
	case "create-user":
		// Creates a user with Beta Premium Subscription
		if userEmail == nil {
			log.Fatal("no email specified")
		}
		emailClient := ses.NewClient(ses.NewClientInput{
			AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
			AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
			AWSRegion:          "us-east-1",
			FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
		})
		if err := tasks.CreateUserWithBetaPremiumSubscription(emailClient, *userEmail); err != nil {
			log.Fatal(err.Error())
		}
	default:
		log.Fatal(fmt.Sprintf("Invalid task specified %s", *taskName))
	}
}

func setupDatabases() error {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		return fmt.Errorf("Error setting up main-db: %s", err.Error())
	}
	if err := elastic.InitializeElasticsearchClientForEnvironment(); err != nil {
		return fmt.Errorf("Error setting up elasticsearch: %s", err.Error())
	}
	return nil
}
