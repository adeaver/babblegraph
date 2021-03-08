package dailyemail

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/usercontenttopics"
	"babblegraph/model/userdocuments"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/ses"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func GetDailyEmailJob(localSentryHub *sentry.Hub, emailClient *ses.Client) func() error {
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
				nErr := fmt.Errorf("Error sending daily email to %s: %s", u.EmailAddress, err.Error())
				log.Println(nErr.Error())
				localSentryHub.CaptureException(nErr)
			}
		}
		return nil
	}
}

func sendDailyEmailToUser(emailClient *ses.Client, user users.User) error {
	var docs []email_actions.CategorizedDocuments
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
		contentTopics, err := usercontenttopics.GetContentTopicsForUser(tx, user.ID)
		if err != nil {
			return err
		}
		emailRecordID, err := email_actions.SendDailyEmailForDocuments(tx, emailClient, recipient, email_actions.DailyEmailInput{
			CategorizedDocuments: docs,
			HasSetTopics:         len(contentTopics) != 0,
		})
		if err != nil {
			return err
		}
		var docIDs []documents.DocumentID
		for _, categorizedDocs := range docs {
			for _, doc := range categorizedDocs.Documents {
				docIDs = append(docIDs, doc.ID)
			}
		}
		return userdocuments.InsertDocumentIDsForUser(tx, user.ID, *emailRecordID, docIDs)
	})
}
