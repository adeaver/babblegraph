package main

import (
	"babblegraph/services/taskrunner/tasks"
	"babblegraph/util/database"
	"babblegraph/util/elastic"
	"babblegraph/util/env"
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
	taskName := flag.String("task", "none", "Name of task to run [daily-email, privacy-policy]")
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
	case "privacy-policy":
		tasks.SendPrivacyPolicyUpdate()
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
