package dailyemail

import (
	"babblegraph/model/documents"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/email"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func SendDailyEmailToUser(emailClient *email.Client, user users.User) error {
	var documents []documents.Document
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		userPreferences, err := getPreferencesForUser(tx, user)
		if err != nil {
			return err
		}
		documents, err = getDocumentsForUser(tx, *userPreferences)
		return err
	}); err != nil {
		return err
	}
	if len(documents) == 0 {
		log.Println(fmt.Sprintf("No documents for user %s", user.EmailAddress))
	}
	if err := sendDailyEmailsForDocuments(emailClient, user.EmailAddress, documents); err != nil {
		return err
	}
	return nil
}
