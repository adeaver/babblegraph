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

	"github.com/jmoiron/sqlx"
)

func main() {
	if err := initializeDatabases(); err != nil {
		log.Fatal(err.Error())
	}
	var allUserEmailInfo []userprefs.UserEmailInfo
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		allUserEmailInfo, err = userprefs.GetActiveUserEmailInfo(tx)
		return err
	}); err != nil {
		log.Fatal(fmt.Sprintf("Error getting email info %s", err.Error()))
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
		log.Fatal(err.Error())
	}
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
