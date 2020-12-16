package dailyemail

import (
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/userdocuments"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/ses"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func GetDailyEmailJob(emailClient *ses.Client) func() error {
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

func sendDailyEmailToUser(emailClient *ses.Client, user users.User) error {
	var docs []documents.Document
	return database.WithTx(func(tx *sqlx.Tx) error {
		userPreferences, err := getPreferencesForUser(tx, user)
		if err != nil {
			return err
		}
		docs, err = getDocumentsForUser(tx, *userPreferences)
		if err != nil {
			return err
		}
		if len(docs) == 0 {
			log.Println(fmt.Sprintf("No documents for user %s", user.EmailAddress))
			return nil
		}
		recipient := email.Recipient{
			UserID:       user.ID,
			EmailAddress: user.EmailAddress,
		}
		emailRecordID, err := email.SendDailyEmailForDocuments(tx, emailClient, recipient, docs)
		if err != nil {
			return err
		}
		var docIDs []documents.DocumentID
		for _, doc := range docs {
			docIDs = append(docIDs, doc.ID)
		}
		return userdocuments.InsertDocumentIDsForUser(tx, user.ID, *emailRecordID, docIDs)
	})
}
