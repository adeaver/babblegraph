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

func GetDailyEmailJob(emailClient *email.Client) func() error {
	return func() error {
		var activeUsers []users.User
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			activeUsers, err = users.GetAllActiveUsers(tx)
			return err
		}); err != nil {
			return err
		}
		for _, u := range activeUsers {
			if err := sendDailyEmailToUser(emailClient, u); err != nil {
				log.Println(fmt.Sprintf("Error sending daily email to %s: %s", u.EmailAddress, err.Error()))
			}
		}
		return nil
	}
}

func sendDailyEmailToUser(emailClient *email.Client, user users.User) error {
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
