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
		userInfo, err = userprefs.GetActiveUserEmailInfo(tx)
		return err
	}); err != nil {
		log.Fatal(fmt.Sprintf("Error getting email info %s", err.Error()))
	}
	emailAddressesToDocuments, err := userquery.GetDocumentsForUser(allUserEmailInfo)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error getting documents for users %s", err.Error()))
	}
	if err := sendutil.SendEmailsToUser(emailAddressesToDocuments); err != nil {
		log.Fatal(fmt.Sprintf("Error sending emails to users %s", err.Error()))
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
