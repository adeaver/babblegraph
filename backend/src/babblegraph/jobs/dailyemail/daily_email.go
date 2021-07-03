package dailyemail

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/model/email"
	"babblegraph/model/useraccounts"
	"babblegraph/model/usercontenttopics"
	"babblegraph/model/usernewsletterschedule"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/ses"
	"fmt"
	"log"
	"time"

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

func SendDailyEmailToUsersByEmailAddress(localSentryHub *sentry.Hub, emailClient *ses.Client) func([]string) error {
	return func(emailAddresses []string) error {
		return database.WithTx(func(tx *sqlx.Tx) error {
			for _, emailAddress := range emailAddresses {
				u, err := users.LookupUserByEmailAddress(tx, emailAddress)
				if err != nil {
					nErr := fmt.Errorf("Error finding user by email address for daily email to %s: %s", emailAddress, err.Error())
					log.Println(nErr.Error())
					localSentryHub.CaptureException(nErr)
					continue
				}
				if u == nil {
					log.Println(fmt.Sprintf("No user for email address %s, continuing", emailAddress))
					continue
				}
				if u.Status != users.UserStatusVerified {
					log.Println(fmt.Sprintf("User %s is not verified, continuing", emailAddress))
					continue
				}
				if err := sendDailyEmailToUser(emailClient, *u); err != nil {
					nErr := fmt.Errorf("Error sending daily email to %s: %s", u.EmailAddress, err.Error())
					log.Println(nErr.Error())
					localSentryHub.CaptureException(nErr)
					continue
				}
			}
			return nil
		})
	}
}

func sendDailyEmailToUser(emailClient *ses.Client, user users.User) error {
	var docs []email_actions.CategorizedDocuments
	return database.WithTx(func(tx *sqlx.Tx) error {
		var userScheduleForDay *usernewsletterschedule.UserNewsletterScheduleDayMetadata
		subscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, user.ID)
		if err != nil {
			return err
		}
		if subscriptionLevel != nil && *subscriptionLevel == useraccounts.SubscriptionLevelBetaPremium {
			var err error
			userScheduleForDay, err = usernewsletterschedule.LookupNewsletterDayMetadataForUserAndDay(tx, user.ID, int(time.Now().UTC().Weekday()))
			if err != nil {
				return err
			}
			if userScheduleForDay != nil && !userScheduleForDay.IsActive {
				log.Println(fmt.Sprintf("Schedule is inactive for current day for user %s. Skipping...", user.ID))
				return nil
			}
		}
		userPreferences, err := getPreferencesForUser(tx, user)
		if err != nil {
			return err
		}
		docs, err = getDocumentsForUser(tx, *userPreferences, userScheduleForDay)
		if err != nil {
			return err
		}
		if len(docs) == 0 {
			return fmt.Errorf("No documents for user %s", user.EmailAddress)
		}
		recipient := email.Recipient{
			UserID:       user.ID,
			EmailAddress: user.EmailAddress,
		}
		contentTopics, err := usercontenttopics.GetContentTopicsForUser(tx, user.ID)
		if err != nil {
			return err
		}
		return email_actions.SendDailyEmailForDocuments(tx, emailClient, recipient, email_actions.DailyEmailInput{
			CategorizedDocuments: docs,
			HasSetTopics:         len(contentTopics) != 0,
		})
	})
}
