package main

import (
	"babblegraph/services/taskrunner/tasks"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/elastic"
	"babblegraph/util/env"
	"babblegraph/util/ses"
	"babblegraph/wordsmith"
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
	taskName := flag.String("task", "none", `Name of task to run
        sample-email: send sample
        privacy-policy: send privacy policy
        create-user: create beta-premium user
        expire-user: expire user
        create-elastic-indexes: create new indices in ElasticSearch
        sync-stripe: sync failed stripe events
        product-updates: send product updates
        admin-content-backfill: backfill all admin content
        admin-content-values-backfill: backfill all admin content values
        admin-content-document-backfill: backfill all documents with new content
        create-admin: create admin`)
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
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
	defer sentry.Flush(2 * time.Second)
	switch *taskName {
	case "sample-email":
		if userEmail == nil {
			log.Fatal("no email specified")
		}
		if err := tasks.SendSampleNewsletter(emailClient, *userEmail); err != nil {
			log.Fatal(err.Error())
		}
	case "create-user":
		// Creates a user with Beta Premium Subscription
		if userEmail == nil {
			log.Fatal("no email specified")
		}
		if err := tasks.CreateUserWithBetaPremiumSubscription(emailClient, *userEmail); err != nil {
			log.Fatal(err.Error())
		}
	case "expire-user":
		// Creates a user with Beta Premium Subscription
		if userEmail == nil {
			log.Fatal("no email specified")
		}
		if err := tasks.DeactivateUserSubscriptionForUser(emailClient, *userEmail); err != nil {
			log.Fatal(err.Error())
		}
	case "privacy-policy":
		tasks.SendPrivacyPolicyUpdate()
	case "create-elastic-indexes":
		if err := tasks.CreateElasticIndexes(); err != nil {
			log.Fatal(err.Error())
		}
	case "sync-stripe":
		tasks.ForceSyncStripeEvents(ctx.GetDefaultLogContext())
	case "product-updates":
		tasks.SendProductUpdates()
	case "create-admin":
		if userEmail == nil {
			log.Fatal("no email specified")
		}
		if err := tasks.CreateAdminAndEmitToken(*userEmail); err != nil {
			log.Fatal(err.Error())
		}
	case "admin-content-backfill":
		if err := tasks.SyncHardcodedContent(); err != nil {
			log.Fatal(err.Error())
		}
	case "admin-content-values-backfill":
		if err := tasks.BackfillAdminContentValues(); err != nil {
			log.Fatal(err.Error())
		}
	case "admin-content-document-backfill":
		if err := tasks.ReindexDocuments(); err != nil {
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
	if err := wordsmith.MustSetupWordsmithForEnvironment(); err != nil {
		return fmt.Errorf("Error setting up wordsmith: %s", err.Error())
	}
	if err := elastic.InitializeElasticsearchClientForEnvironment(); err != nil {
		return fmt.Errorf("Error setting up elasticsearch: %s", err.Error())
	}
	return nil
}
