package scheduler

import (
	"babblegraph/model/useraccountsnotifications"
	"babblegraph/util/database"
	"babblegraph/util/ses"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func handlePendingUserAccountNotificatioRequests(localSentryHub *sentry.Hub, emailClient *ses.Client) error {
	var notificationRequests []useraccountsnotifications.NotificationRequest
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		notificationRequests, err = useraccountsnotifications.GetNotificationsToFulfill(tx)
		return err
	}); err != nil {
		localSentryHub.CaptureException(err)
		return err
	}
	for _, req := range notificationRequests {
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			if err := useraccountsnotifications.FulfillNotificationRequest(tx, req.ID); err != nil {
				return err
			}
			switch req.Type {
			case useraccountsnotifications.NotificationTypeTrialEndingSoon:
				// Send Trial Ending Soon email
			case useraccountsnotifications.NotificationTypeAccountCreated:
				// Send account creation notificaiton
			case useraccountsnotifications.NotificationTypePaymentError:
				// Send payment error notification
			case useraccountsnotifications.NotificationTypePremiumSubscriptionCanceled:
				// Send subscription canceled notification
			default:
				return fmt.Errorf("Unknown notification type %s", req.Type)
			}
			return nil
		}); err != nil {
			localSentryHub.CaptureException(err)
		}
	}
	return nil
}
