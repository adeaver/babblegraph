package main

import (
	"babblegraph/services/email/sendutil"
	"babblegraph/services/email/userprefs"
	"babblegraph/services/email/userquery"
	"babblegraph/util/database"
	"babblegraph/util/elastic"
	"babblegraph/wordsmith"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	cron "github.com/robfig/cron/v3"
)

func main() {
	if err := initializeDatabases(); err != nil {
		log.Fatal(err.Error())
	}
	errs := make(chan error, 1)
	usEastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err.Error())
	}
	c := cron.New(cron.WithLocation(usEastern))
	c.AddFunc("30 5 * * *", func() {
		log.Println("Starting email job...")
		var allUserEmailInfo []userprefs.UserEmailInfo
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			allUserEmailInfo, err = userprefs.GetActiveUserEmailInfo(tx)
			return err
		}); err != nil {
			errs <- err
			return
		}
		log.Println(fmt.Sprintf("Sending emails to %d address", len(allUserEmailInfo)))
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			emailAddressesToDocuments, err := userquery.GetDocumentsForUser(tx, allUserEmailInfo)
			if err != nil {
				return fmt.Errorf("Error getting documents for users %s", err.Error())
			}
			err = sendutil.SendEmailsToUser(emailAddressesToDocuments)
			if err != nil {
				return fmt.Errorf("Error sending emails to users %s", err.Error())
			}
			return err
		}); err != nil {
			errs <- err
			return
		}
		return
	})
	c.Start()
	log.Println(c.Entries())
	err = <-errs
	log.Println("Error detected")
	if err != nil {
		log.Fatal(err.Error())
	}
	c.Stop()
	log.Println("Ending email service")
}

func initializeDatabases() error {
	if err := database.GetDatabaseForEnvironmentRetrying(); err != nil {
		return fmt.Errorf("error connecting to main-db: %s", err.Error())
	}
	log.Println("successfully connected to main db")
	if err := wordsmith.MustSetupWordsmithForEnvironment(); err != nil {
		return fmt.Errorf("error connecting to wordsmith: %s", err.Error())
	}
	log.Println("successfully connected to wordsmith db")
	if err := elastic.InitializeElasticsearchClientForEnvironment(); err != nil {
		return fmt.Errorf("error connecting to elasticsearch: %s", err.Error())
	}
	log.Println("successfully connected to elasticsearch")
	if err := sendutil.InitializeEmailClient(); err != nil {
		return fmt.Errorf("error setting up email client: %s", err.Error())
	}
	log.Println("successfully setup email client")
	return nil
}
