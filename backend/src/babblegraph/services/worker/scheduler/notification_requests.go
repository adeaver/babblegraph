package scheduler

import (
	"babblegraph/model/billing"
	"babblegraph/model/email"
	"babblegraph/model/emailtemplates"
	"babblegraph/model/routes"
	"babblegraph/model/useraccountsnotifications"
	"babblegraph/model/users"
	"babblegraph/util/async"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"babblegraph/util/ses"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func handlePendingUserAccountNotificationRequests(c async.Context) {
	emailClient := ses.NewClient(ses.NewClientInput{
		AWSAccessKey:       env.MustEnvironmentVariable("AWS_SES_ACCESS_KEY"),
		AWSSecretAccessKey: env.MustEnvironmentVariable("AWS_SES_SECRET_KEY"),
		AWSRegion:          "us-east-1",
		FromAddress:        env.MustEnvironmentVariable("EMAIL_ADDRESS"),
	})
	var notificationRequests []useraccountsnotifications.NotificationRequest
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		notificationRequests, err = useraccountsnotifications.GetNotificationsToFulfill(tx)
		return err
	}); err != nil {
		c.Errorf("Error fetching notifications: %s", err.Error())
		return
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
			if user.Status != users.UserStatusVerified {
				c.Infof("User is not verified, not sending")
				return nil
			}
			var subject, emailHTML *string
			var emailType *email.EmailType
			emailRecordID := email.NewEmailRecordID()
			switch req.Type {
			case useraccountsnotifications.NotificationTypeTrialEndingSoon:
				subject = ptr.String("Attention! Your Babblegraph trial is ending soon!")
				emailHTML, emailType, err = handleTrialEndingSoonNotification(c, tx, emailRecordID, user)
			case useraccountsnotifications.NotificationTypePremiumSubscriptionCanceled:
				subject = ptr.String("Your Babblegraph subscription has ended")
				emailHTML, emailType, err = handleSubscriptionCanceledNotification(c, tx, emailRecordID, user)
			case useraccountsnotifications.NotificationTypeNeedPaymentMethodWarning,
				useraccountsnotifications.NotificationTypeNeedPaymentMethodWarningVeryUrgent:
				subject = ptr.String("Attention! Add a payment method to keep using Babblegraph")
				emailHTML, emailType, err = handleNeedPaymentMethodWarningNotification(c, tx, emailRecordID, user)
			default:
				return fmt.Errorf("Unknown notification type %s", req.Type)
			}
			if emailHTML == nil {
				return nil
			}
			if err := email.InsertEmailRecord(tx, emailRecordID, user.ID, *emailType); err != nil {
				return err
			}
			return email.SendEmailWithHTMLBody(tx, emailClient, email.SendEmailWithHTMLBodyInput{
				ID:           emailRecordID,
				EmailAddress: user.EmailAddress,
				Subject:      *subject,
				Body:         *emailHTML,
			})
		}); err != nil {
			c.Errorf("Error fulfilling request %s: %s", req.ID, err.Error())
		}
	}
}

func handleTrialEndingSoonNotification(c ctx.LogContext, tx *sqlx.Tx, emailRecordID email.ID, user *users.User) (*string, *email.EmailType, error) {
	premiumSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(c, tx, user.ID)
	switch {
	case err != nil:
		return nil, nil, err
	case premiumSubscription == nil:
		c.Warnf("Need to send trial ending soon, but there is no subscription for user %s", user.ID)
		return nil, nil, nil
	default:
		userAccessor, err := emailtemplates.GetDefaultUserAccessor(tx, user.ID)
		if err != nil {
			return nil, nil, err
		}
		beforeParagraphs := []string{
			"Hello!",
			"We hope you’ve been enjoying Babblegraph. If not, your feedback is always appreciated!",
			"Just respond to this email with any ideas or comments you have about what Babblegraph could be doing better.",
			"This email is to let you know that your trial is almost over.",
		}
		var action *emailtemplates.GenericEmailAction
		switch premiumSubscription.PaymentState {
		case billing.PaymentStateTrialNoPaymentMethod:
			beforeParagraphs = append(beforeParagraphs, "It looks like you haven’t added a payment method yet. If you’d like to continue using Babblegraph, then you’ll need to add a payment method. You can do that at this link")
			checkoutLink, err := routes.MakePremiumSubscriptionCheckoutLink(user.ID)
			if err != nil {
				return nil, nil, err
			}
			action = &emailtemplates.GenericEmailAction{
				Link:       *checkoutLink,
				ButtonText: "Add a payment method to your account",
			}
		case billing.PaymentStateTrialPaymentMethodAdded:
			paymentSettingsRoute, err := routes.MakePaymentSettingsRouteForUserID(user.ID)
			if err != nil {
				return nil, nil, err
			}
			if !premiumSubscription.IsAutoRenewEnabled {
				beforeParagraphs = append(beforeParagraphs, "You have a payment method added to your account, but your subscription will is not set to automatically renew. If you would like to continue to use Babblegraph, you’ll need to turn auto-renew on.")
				action = &emailtemplates.GenericEmailAction{
					Link:       *paymentSettingsRoute,
					ButtonText: "Manage your subscription settings here",
				}
			} else {
				beforeParagraphs = append(beforeParagraphs, "You have a payment method added to your account, so your subscription will automatically renew in a few days, and you’ll be charged $29 then.")
			}
		case billing.PaymentStateCreatedUnpaid,
			billing.PaymentStateTerminated,
			billing.PaymentStateErrored,
			billing.PaymentStateActive:
			return nil, nil, fmt.Errorf("Invalid payment state for premium subscription %s: %d", *premiumSubscription.ID, premiumSubscription.PaymentState)
		}

		emailHTML, err := emailtemplates.MakeGenericUserEmailHTML(emailtemplates.MakeGenericUserEmailHTMLInput{
			EmailRecordID:      emailRecordID,
			UserAccessor:       userAccessor,
			EmailTitle:         "Your Babblegraph trial is ending soon",
			PreheaderText:      "Your trial of Babblegraph is set to expire in the next few days",
			BeforeParagraphs:   beforeParagraphs,
			GenericEmailAction: action,
			AfterParagraphs: []string{
				"If you have any questions or believe there is an error in this email, just respond to this email.",
				"Thank you so much for trying out Babblegraph!",
			},
		})
		if err != nil {
			return nil, nil, err
		}
		return emailHTML, email.EmailTypeTrialEndingSoon.Ptr(), nil
	}
}

func handleSubscriptionCanceledNotification(c ctx.LogContext, tx *sqlx.Tx, emailRecordID email.ID, user *users.User) (*string, *email.EmailType, error) {
	premiumSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(c, tx, user.ID)
	switch {
	case err != nil:
		return nil, nil, err
	case premiumSubscription != nil:
		c.Warnf("Need to send subscription ended, but there is an active subscription for user %s", user.ID)
		return nil, nil, nil
	default:
		userAccessor, err := emailtemplates.GetDefaultUserAccessor(tx, user.ID)
		if err != nil {
			return nil, nil, err
		}
		emailHTML, err := emailtemplates.MakeGenericUserEmailHTML(emailtemplates.MakeGenericUserEmailHTMLInput{
			EmailRecordID: emailRecordID,
			UserAccessor:  userAccessor,
			EmailTitle:    "Your Babblegraph subscription has ended",
			PreheaderText: "Your Babblegraph subscription has expired",
			BeforeParagraphs: []string{
				"Hello!",
				"Thanks so much for trying Babblegraph!",
				"This email is to let you know that your subscription to Babblegraph has ended. You will no longer be charged for the subscription.",
				"If you see any new charges on your credit card statement from Babblegraph, just respond to this email and we’ll get it sorted out.",
				"Lastly, before you go, we’d love to know what we could do better! You can respond directly to this email to give feedback.",
				"If you have any questions or believe there is an error in this email, just respond to this email.",
				"Thanks again so much for trying out Babblegraph!",
			},
		})
		if err != nil {
			return nil, nil, err
		}
		return emailHTML, email.EmailTypePremiumSubscriptionCanceled.Ptr(), nil
	}
}

func handleNeedPaymentMethodWarningNotification(c ctx.LogContext, tx *sqlx.Tx, emailRecordID email.ID, user *users.User) (*string, *email.EmailType, error) {
	premiumSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(c, tx, user.ID)
	switch {
	case err != nil:
		return nil, nil, err
	case premiumSubscription == nil:
		c.Warnf("Need to send payment method warning email, but there is no subscription for user %s", user.ID)
		return nil, nil, nil
	default:
		userAccessor, err := emailtemplates.GetDefaultUserAccessor(tx, user.ID)
		if err != nil {
			return nil, nil, err
		}
		switch premiumSubscription.PaymentState {
		case billing.PaymentStateTrialNoPaymentMethod:
			checkoutLink, err := routes.MakePremiumSubscriptionCheckoutLink(user.ID)
			if err != nil {
				return nil, nil, err
			}
			emailHTML, err := emailtemplates.MakeGenericUserEmailHTML(emailtemplates.MakeGenericUserEmailHTMLInput{
				EmailRecordID: emailRecordID,
				UserAccessor:  userAccessor,
				EmailTitle:    "Your trial is about to expire with no payment method",
				PreheaderText: "Your trial is ending, but you don’t have a pyament method added.",
				BeforeParagraphs: []string{
					"Hello!",
					"This email is to let you know that your subscription is ending soon, but you have no payment method added to your account.",
					"If you would like to continue to use Babblegraph, you’ll need to add a payment method. You can do that at the link below.",
				},
				GenericEmailAction: &emailtemplates.GenericEmailAction{
					Link:       *checkoutLink,
					ButtonText: "Add a payment method to your account",
				},
				AfterParagraphs: []string{
					"If you have any questions or believe there is an error in this email, just respond to this email.",
					"Thanks again so much for trying out Babblegraph!",
				},
			})
			if err != nil {
				return nil, nil, err
			}
			return emailHTML, email.EmailTypeTrialEndingSoonActionRequired.Ptr(), nil
		case billing.PaymentStateTerminated,
			billing.PaymentStateErrored,
			billing.PaymentStateActive,
			billing.PaymentStateTrialPaymentMethodAdded:
			return nil, nil, nil
		}
	}
	return nil, nil, nil
}
