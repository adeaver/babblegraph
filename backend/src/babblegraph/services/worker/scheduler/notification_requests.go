package scheduler

import (
	email_actions "babblegraph/actions/email"
	"babblegraph/externalapis/bgstripe"
	"babblegraph/model/email"
	"babblegraph/model/routes"
	"babblegraph/model/useraccountsnotifications"
	"babblegraph/model/users"
	"babblegraph/util/database"
	"babblegraph/util/ses"
	"fmt"
	"log"

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
			user, err := users.GetUser(tx, req.UserID)
			if err != nil {
				return err
			}
			switch req.Type {
			case useraccountsnotifications.NotificationTypeTrialEndingSoon:
				return handleTrialEndingSoonNotification(tx, emailClient, *user, req)
			case useraccountsnotifications.NotificationTypeAccountCreated:
				return handleAccountCreationNotification(tx, emailClient, *user)
			case useraccountsnotifications.NotificationTypePaymentError:
				return handlePaymentFailureNotification(tx, emailClient, *user)
			case useraccountsnotifications.NotificationTypePremiumSubscriptionCanceled:
				return handlePremiumSubscriptionCanceledNotification(tx, emailClient, *user)
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

func handleTrialEndingSoonNotification(tx *sqlx.Tx, cl *ses.Client, user users.User, req useraccountsnotifications.NotificationRequest) error {
	subscriptionDetails, _, err := bgstripe.LookupNonterminatedStripeSubscriptionForUser(tx, req.UserID)
	switch {
	case err != nil:
		return err
	case subscriptionDetails == nil:
		// User does not have an active subscription, log and eat the message
		log.Println(fmt.Sprintf("User %s does not have an active subscription", req.UserID))
		return nil
	default:
		paymentSettingsLink := routes.MakeLoginLinkWithPaymentSettingsRedirectKey()
		var emailInput *email_actions.SendGenericEmailWithOptionalActionForRecipientInput
		switch subscriptionDetails.PaymentState {
		case bgstripe.PaymentStateTrialNoPaymentMethod:
			emailInput = &email_actions.SendGenericEmailWithOptionalActionForRecipientInput{
				EmailType: email.EmailTypeTrialEndingSoonActionRequired,
				Recipient: email.Recipient{
					UserID:       user.ID,
					EmailAddress: user.EmailAddress,
				},
				Subject:       "ACTION REQUIRED: Trial Ending Soon",
				EmailTitle:    "ACTION REQUIRED: Trial Ending Soon",
				PreheaderText: "Your trial is ending in a few days and there is action required!",
				BeforeParagraphs: []string{
					"Hello!",
					"Thank you so much for trying out Babblegraph Premium. I hope you’ve enjoyed it!",
					"Your trial is ending in a few days, and it looks like you don’t have a payment method attached to your account.",
					"Without a payment method attached, your access to premium will expire when your trial ends. If you do not wish to continue with premium, then no action is needed on your part. If you do, please add a payment method below.",
					"You will still continue to receive the Babblegraph newsletter",
				},
				GenericEmailAction: &email_actions.GenericEmailAction{
					Link:       paymentSettingsLink,
					ButtonText: "Add a payment method",
				},
				AfterParagraphs: []string{
					"Thank you again for trying out Babblegraph Premium. If you have any questions or need any help, just reply to this email!",
				},
			}
		case bgstripe.PaymentStateTrialPaymentMethodAdded:
			emailInput = &email_actions.SendGenericEmailWithOptionalActionForRecipientInput{
				EmailType: email.EmailTypeTrialEndingSoon,
				Recipient: email.Recipient{
					UserID:       user.ID,
					EmailAddress: user.EmailAddress,
				},
				Subject:       "Your Trial is Ending Soon",
				EmailTitle:    "Your Trial is Ending Soon",
				PreheaderText: "Your trial is ending in a few days!",
				BeforeParagraphs: []string{
					"Hello!",
					"Thank you so much for trying out Babblegraph Premium. I hope you’ve enjoyed it!",
					"Your trial is ending in a few days. In the next few days, you’ll see a charge from Babblegraph on your card statement!",
					"If you don’t wish to be charged, you can cancel your subscription with the link below. You will retain access to Babblegraph Premium until your trial period is over.",
				},
				GenericEmailAction: &email_actions.GenericEmailAction{
					Link:       paymentSettingsLink,
					ButtonText: "Cancel your subscription",
				},
				AfterParagraphs: []string{
					"If you want to continue using Babblegraph Premium, then there’s no action on your part! Thank you for supporting Babblegraph!",
					"If you have any questions or need any help, just reply to this email.",
				},
			}
		case bgstripe.PaymentStateActive:
			return nil
		case bgstripe.PaymentStateCreatedUnpaid,
			bgstripe.PaymentStateErrored,
			bgstripe.PaymentStateTerminated:
			// Log because this error is not retryable
			log.Println(fmt.Sprintf("User %s has a subscription in state %d, but expected either 1 or 2", req.UserID, subscriptionDetails.PaymentState))
			return nil
		default:
			log.Println(fmt.Sprintf("User %s has a subscription in an unrecognized state %d, but expected either 1 or 2", req.UserID, subscriptionDetails.PaymentState))
			return nil

		}
		if emailInput == nil {
			return fmt.Errorf("unreachable")
		}
		if _, err := email_actions.SendGenericEmailWithOptionalActionForRecipient(tx, cl, *emailInput); err != nil {
			return err
		}
		return nil
	}
}

func handlePremiumSubscriptionCanceledNotification(tx *sqlx.Tx, cl *ses.Client, user users.User) error {
	emailInput := &email_actions.SendGenericEmailWithOptionalActionForRecipientInput{
		EmailType: email.EmailTypePremiumSubscriptionCanceled,
		Recipient: email.Recipient{
			UserID:       user.ID,
			EmailAddress: user.EmailAddress,
		},
		Subject:       "Your Premium subscription has ended.",
		EmailTitle:    "Your Premium subscription has ended.",
		PreheaderText: "Your Premium subscription has ended.",
		BeforeParagraphs: []string{
			"Hello!",
			"Thank you so much for trying out Babblegraph Premium! This email is a confirmation that your subscription has ended.",
			"This means that you no longer have access to Premium features.",
			"However, you will continue to receive the daily newsletter, unless you unsubscribed!",
		},
		AfterParagraphs: []string{
			"If you think that you received this email by mistake or if you have any other questions or concerns, then just reply to this email.",
		},
	}
	if _, err := email_actions.SendGenericEmailWithOptionalActionForRecipient(tx, cl, *emailInput); err != nil {
		return err
	}
	return nil
}

func handlePaymentFailureNotification(tx *sqlx.Tx, cl *ses.Client, user users.User) error {
	paymentSettingsLink := routes.MakeLoginLinkWithPaymentSettingsRedirectKey()
	emailInput := &email_actions.SendGenericEmailWithOptionalActionForRecipientInput{
		EmailType: email.EmailTypePaymentFailureNotification,
		Recipient: email.Recipient{
			UserID:       user.ID,
			EmailAddress: user.EmailAddress,
		},
		Subject:       "ACTION REQUIRED: Payment Failed",
		EmailTitle:    "ACTION REQUIRED: Payment Failed",
		PreheaderText: "A recent payment attempt failed.",
		BeforeParagraphs: []string{
			"Hello!",
			"There was a failure to charge the default payment method on your account.",
			"Double check to make sure that your payment information is correct with the link below, as well as making sure with your financial institution that the transaction is allowed!",
			"We will automatically retry with your current payment method. Sometimes, these payment failures happen and get resolved on their own.",
			"If it keeps failing, you will lose access to Babblegraph Premium",
		},
		GenericEmailAction: &email_actions.GenericEmailAction{
			Link:       paymentSettingsLink,
			ButtonText: "Check your payment settings",
		},
		AfterParagraphs: []string{
			"If you have any questions or need any help, just reply to this email!",
		},
	}
	if _, err := email_actions.SendGenericEmailWithOptionalActionForRecipient(tx, cl, *emailInput); err != nil {
		return err
	}
	return nil
}

func handleAccountCreationNotification(tx *sqlx.Tx, cl *ses.Client, user users.User) error {
	emailInput := &email_actions.SendGenericEmailWithOptionalActionForRecipientInput{
		EmailType: email.EmailTypeAccountCreationNotification,
		Recipient: email.Recipient{
			UserID:       user.ID,
			EmailAddress: user.EmailAddress,
		},
		Subject:       "Account Creation Confirmation",
		EmailTitle:    "Account Creation Confirmation",
		PreheaderText: "Your account was successfully created",
		BeforeParagraphs: []string{
			"Hello!",
			"This email is to let you know that your Babblegraph account was successfully created.",
			"If you did not initiate this, please respond to this email!",
			"If you did initiate this, then no further action is required on your part.",
		},
	}
	if _, err := email_actions.SendGenericEmailWithOptionalActionForRecipient(tx, cl, *emailInput); err != nil {
		return err
	}
	return nil
}
