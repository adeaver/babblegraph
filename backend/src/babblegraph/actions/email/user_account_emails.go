package email

import (
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/util/env"
	"babblegraph/util/ses"

	"github.com/jmoiron/sqlx"
)

func SendUserCreationEmailForRecipient(tx *sqlx.Tx, cl *ses.Client, recipient email.Recipient) (*email.ID, error) {
	signupLink, err := routes.MakeUserCreationLink(recipient.UserID)
	if err != nil {
		return nil, err
	}
	return SendGenericEmailWithOptionalActionForRecipient(tx, cl, SendGenericEmailWithOptionalActionForRecipientInput{
		EmailType:     email.EmailTypeUserCreation,
		Recipient:     recipient,
		Subject:       "Create your Babblegraph Account to Access Premium Features",
		EmailTitle:    "Babblegraph Premium Subscription",
		PreheaderText: "Create an account to use premium features on Babblegraph.",
		BeforeParagraphs: []string{
			"Hola!",
			"Thank you so much for supporting Babblegraph. I’m giving access to not-yet-released features as well as some premium features to everyone who supports Babblegraph. The first step is to create an account, which you can do with the link below.",
		},
		GenericEmailAction: &GenericEmailAction{
			ButtonText: "Click here to create your premium account",
			Link:       *signupLink,
		},
		AfterParagraphs: []string{
			"Once you’ve signed up, you’ll be able to access the new Babblegraph features! If you have any questions or ideas for features that you’d like to see, you can reply to this email or any of your daily emails.",
			"Thanks again for supporting Babblegraph!",
		},
	})
}

func SendAccountReactivationEmailForRecipient(tx *sqlx.Tx, cl *ses.Client, recipient email.Recipient) (*email.ID, error) {
	return SendGenericEmailWithOptionalActionForRecipient(tx, cl, SendGenericEmailWithOptionalActionForRecipientInput{
		EmailType:     email.EmailTypeUserReactivation,
		Recipient:     recipient,
		Subject:       "Thank you for supporting Babblegraph!",
		EmailTitle:    "Babblegraph Premium Subscription",
		PreheaderText: "Your premium subscription has been reactivated as a thank you for supporting Babblegraph",
		BeforeParagraphs: []string{
			"Hola!",
			"Thank you so much for supporting Babblegraph again. I’m still giving access to not-yet-released features as well as some premium features to everyone who supports Babblegraph. As you already had an account with Babblegraph, it has been given access to premium and beta features once again. You can still login with the same password you’ve been using. You can login with the link below. If you forgot your password, there is a link to reset it at the page below.",
		},
		GenericEmailAction: &GenericEmailAction{
			ButtonText: "Click here to login to your account",
			Link:       env.GetAbsoluteURLForEnvironment("login"),
		},
		AfterParagraphs: []string{
			"Thanks again for supporting Babblegraph!",
		},
	})
}

func SendAccountExpirationEmail(tx *sqlx.Tx, cl *ses.Client, recipient email.Recipient) (*email.ID, error) {
	return SendGenericEmailWithOptionalActionForRecipient(tx, cl, SendGenericEmailWithOptionalActionForRecipientInput{
		EmailType:     email.EmailTypeUserExpiration,
		Recipient:     recipient,
		Subject:       "Sad to see you go!",
		EmailTitle:    "Babblegraph Premium Subscription",
		PreheaderText: "Your access to Babblegraph Premium features has expired",
		BeforeParagraphs: []string{
			"Hola!",
			"Thank you so much for supporting Babblegraph! It’s support from people like you that enable me to continue to work on and grow Babblegraph.",
			"This email is to let you know that your access to premium features has expired. You will still receive the daily email. If you would like to unsubscribe, click on the link at the bottom of this email.",
			"If you think this is an error, just reply to this email to let me know! You can still log into your account; however, you will not be able to access premium features. If you change your mind and want access again, head on over to Babblegraph’s Buy Me A Coffee page!",
			"Thanks again for supporting Babblegraph!",
		},
	})
}