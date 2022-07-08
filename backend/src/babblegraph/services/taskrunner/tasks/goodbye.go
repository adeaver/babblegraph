package tasks

import (
	"babblegraph/model/email"
	"babblegraph/model/emailtemplates"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"

	"github.com/jmoiron/sqlx"
)

func SendGoodbyeEmail(cl *ses.Client) error {
	var u []users.User
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		u, err = users.GetAllActiveUsers(tx)
		return err
	}); err != nil {
		return err
	}
	defaultInput := emailtemplates.MakeGenericUserEmailHTMLInput{
		EmailTitle:    "Babblegraph is shutting down at the end of the month",
		PreheaderText: "This email is to let you know that we’re shutting down.",
		BeforeParagraphs: []string{
			"Hello!",
			"At the end of this month, I’ll be shutting down Babblegraph.",
			"When this happens, all of the data will be completely destroyed and no longer be accessible by anyone, including you and me. If you have any questions, let me know by responding to this email.",
			"Babblegraph started as a side project of mine when I wanted a better way to incorporate Spanish into my daily routine. About a year and a half ago, I decided to see how much I could grow it! Since then it’s been a whirlwind adventure that has seen thousands of people sign up and give it a try. I’m truly humbled by how much interest it has gotten. I truly could have never foreseen it.",
			"I really loved building Babblegraph, and for months, I really strived to turn Babblegraph from a side project into a full-fledged company. Along the way, I learned a lot about what makes a great product and what makes a great company.",
			"Unfortunately, I also couldn't forsee what the challenges of that transition were. As a solo developer, it really put a lot of stress on me. In the end, it really ended up draining the love I had for building Babblegraph when it first got started.",
			"As with any product or service that you use, you should really expect that the person providing it to you has their heart in it, and for Babblegraph, that just isn't as true as it was a year and a half ago. Therefore, I’ve made the tough decision to let it go and stop working on it.",
			"And thank you so much for giving it a try!",
		},
		ExcludeFooter: true,
	}
	c := ctx.GetDefaultLogContext()
	for _, user := range u {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			sub, err := useraccounts.LookupSubscriptionLevelForUser(tx, user.ID)
			switch {
			case err != nil:
				return err
			case sub == nil:
				return nil
			}
			emailRecordID := email.NewEmailRecordID()
			if err := email.InsertEmailRecord(tx, emailRecordID, user.ID, email.EmailTypeGoodbye); err != nil {
				return err
			}
			emailAccessor, err := emailtemplates.GetDefaultUserAccessor(tx, user.ID)
			if err != nil {
				return err
			}
			defaultInput.UserAccessor = emailAccessor
			defaultInput.EmailRecordID = emailRecordID
			body, err := emailtemplates.MakeGenericUserEmailHTML(defaultInput)
			if err != nil {
				return err
			}
			return email.SendEmailWithHTMLBody(tx, cl, email.SendEmailWithHTMLBodyInput{
				ID:              emailRecordID,
				EmailAddress:    user.EmailAddress,
				Subject:         "So long!",
				Body:            *body,
				EmailSenderName: ptr.String("Andrew from Babblegraph"),
			})
		}); err != nil {
			c.Errorf("Error ending email to %s: %s", user.ID, err.Error())
		}
	}
	return nil
}
