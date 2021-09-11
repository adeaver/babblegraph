package email

import (
	"babblegraph/util/ses"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type SendEmailWithHTMLBodyInput struct {
	ID              ID
	EmailAddress    string
	Subject         string
	EmailSenderName *string
	Body            string
}

func SendEmailWithHTMLBody(tx *sqlx.Tx, cl *ses.Client, input SendEmailWithHTMLBodyInput) error {
	if err := SetEmailRecordSentAtTime(tx, input.ID); err != nil {
		return err
	}
	sesMessageID, err := cl.SendEmail(ses.SendEmailInput{
		Recipient:       input.EmailAddress,
		HTMLBody:        input.Body,
		Subject:         input.Subject,
		EmailSenderName: input.EmailSenderName,
	})
	if err != nil {
		return err
	}
	if err := UpdateEmailRecordIDWithSESMessageID(tx, input.ID, *sesMessageID); err != nil {
		// **** VERY IMPORTANT HERE ****
		// This *cannot* return an error if it fails
		// since SES has already successfully sent the email and returning an error
		// causes the transaction to abort and rollback. However, SES has a side effect.
		// The only thing that this does is update the email record
		// to have the SES message ID. It is not super important - but it is useful.
		log.Println(fmt.Sprintf("Error updating email record %s with SES Message ID %s: %s", input.ID, *sesMessageID, err.Error()))
	}
	return nil
}
