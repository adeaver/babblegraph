package tasks

import (
	"babblegraph/model/documents"
	"babblegraph/model/email"
	"babblegraph/model/emailtemplates"
	"babblegraph/model/newsletter"
	"babblegraph/model/userlemma"
	"babblegraph/model/usernewsletterpreferences"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"babblegraph/util/random"
	"babblegraph/util/ses"
	"babblegraph/util/timeutils"
	"babblegraph/wordsmith"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func SendSampleNewsletter(cl *ses.Client, emailAddress string) error {
	c := ctx.GetDefaultLogContext()
	switch env.MustEnvironmentName() {
	case env.EnvironmentLocal,
		env.EnvironmentLocalTestEmail:
		return database.WithTx(func(tx *sqlx.Tx) error {
			user, err := users.LookupUserByEmailAddress(tx, emailAddress)
			switch {
			case err != nil:
				return err
			case user == nil:
				return fmt.Errorf("No user found for that email address")
			case user.Status != users.UserStatusVerified:
				return fmt.Errorf("User found was not verified")
			}
			emailRecordID := email.NewEmailRecordID()
			newsletter, err := createNewsletter(c, tx, user.ID, emailRecordID)
			if err != nil {
				return err
			}
			if newsletter == nil {
				c.Infof("Newsletter is null")
				return nil
			}
			return createNewsletterHTMLAndSend(cl, tx, emailAddress, user.ID, emailRecordID, newsletter.Body)
		})
	case env.EnvironmentStage,
		env.EnvironmentLocalNoEmail,
		env.EnvironmentProd,
		env.EnvironmentTest:
		return fmt.Errorf("Can't send sample email in non-local environment")
	default:
		return fmt.Errorf("Unrecognized environment")
	}
}

func createNewsletter(c ctx.LogContext, tx *sqlx.Tx, userID users.UserID, emailRecordID email.ID) (*newsletter.Newsletter, error) {
	emailAccessor := newsletter.GetDefaultEmailAccessor(tx)
	documentAccessor := newsletter.GetDefaultDocumentsAccessor()
	utcMidnight := timeutils.ConvertToMidnight(time.Now().UTC())
	userNewsletterPreferences, err := usernewsletterpreferences.GetUserNewsletterPrefrencesForLanguage(c, tx, userID, wordsmith.LanguageCodeSpanish, ptr.Time(utcMidnight))
	if err != nil {
		return nil, err
	}
	userAccessor, err := newsletter.GetSampleNewsletterUserAccessor(c, tx, newsletter.GetSampleNewsletterUserAccessorInput{
		UserID:                    userID,
		LanguageCode:              wordsmith.LanguageCodeSpanish,
		SentDocumentIDs:           []documents.DocumentID{},
		SpotlightRecords:          []userlemma.UserLemmaReinforcementSpotlightRecord{},
		UserNewsletterPreferences: userNewsletterPreferences,
	})
	if err != nil {
		return nil, err
	}
	contentAccessor, err := newsletter.GetDefaultContentAccessor(tx, wordsmith.LanguageCodeSpanish)
	if err != nil {
		return nil, err
	}
	podcastAccessor, err := newsletter.GetDefaultPodcastAccessor(c, tx, wordsmith.LanguageCodeSpanish, userID)
	if err != nil {
		return nil, err
	}
	advertisementAccessor, err := newsletter.GetDefaultAdvertisementAccessor(tx, userID, wordsmith.LanguageCodeSpanish)
	if err != nil {
		return nil, err
	}
	return newsletter.CreateNewsletter(c, newsletter.CreateNewsletterInput{
		DateOfSendMidnightUTC: utcMidnight,
		WordsmithAccessor:     newsletter.GetDefaultWordsmithAccessor(),
		EmailAccessor:         emailAccessor,
		UserAccessor:          userAccessor,
		DocsAccessor:          documentAccessor,
		ContentAccessor:       contentAccessor,
		PodcastAccessor:       podcastAccessor,
		AdvertisementAccessor: advertisementAccessor,
	})
}

func createNewsletterHTMLAndSend(emailClient *ses.Client, tx *sqlx.Tx, emailAddress string, userID users.UserID, emailRecordID email.ID, newsletterBody newsletter.NewsletterBody) error {
	userAccessor, err := emailtemplates.GetDefaultUserAccessor(tx, userID)
	if err != nil {
		return err
	}
	newsletterHTML, err := emailtemplates.MakeNewsletterHTML(emailtemplates.MakeNewsletterHTMLInput{
		EmailRecordID: emailRecordID,
		UserAccessor:  userAccessor,
		Body:          newsletterBody,
	})
	if err != nil {
		return err
	}
	return email.SendEmailWithHTMLBody(tx, emailClient, email.SendEmailWithHTMLBodyInput{
		ID:              emailRecordID,
		EmailAddress:    emailAddress,
		Subject:         fmt.Sprintf("Sample Newsletter - %s", random.MustMakeRandomString(5)),
		EmailSenderName: ptr.String("Babblegraph Sample Email"),
		Body:            *newsletterHTML,
	})
}
