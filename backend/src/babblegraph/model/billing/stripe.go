package billing

import (
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/setupintent"
	"github.com/stripe/stripe-go/v72/webhook"
)

// This file is for any stripe specific methods

const (
	captureStripeEventQuery = "INSERT INTO billing_stripe_event (type, processed, data) VALUES ($1, $2, $3)"
)

func GetSetupIntentClientSecretForUser(tx *sqlx.Tx, userID users.UserID) (*string, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return nil, err
	case billingInformation == nil:
		return nil, fmt.Errorf("Expected billing information for user %s, but got none", userID)
	default:
		externalID, err := getExternalIDMapping(tx, billingInformation.ExternalIDMappingID)
		if err != nil {
			return nil, err
		}
		if externalID.IDType != externalIDTypeStripe {
			return nil, fmt.Errorf("User %s is not a stripe user, has type %s", userID, externalID.IDType)
		}
		params := &stripe.SetupIntentParams{
			Customer: ptr.String(externalID.ExternalID),
			PaymentMethodTypes: []*string{
				stripe.String("card"),
			},
			Usage: ptr.String("off_session"),
		}
		si, err := setupintent.New(params)
		if err != nil {
			return nil, err
		}
		return ptr.String(si.ClientSecret), nil
	}
}

func convertStripeSubscriptionToPremiumNewsletterSubscription(tx *sqlx.Tx, stripeSubscription *stripe.Subscription, dbNewsletterSubscription *dbPremiumNewsletterSubscription) (*PremiumNewsletterSubscription, error) {
	var paymentIntentID *string
	if stripeSubscription.LatestInvoice != nil && stripeSubscription.LatestInvoice.PaymentIntent != nil {
		paymentIntentID = ptr.String(stripeSubscription.LatestInvoice.PaymentIntent.ClientSecret)
	}
	var priceCents *int64
	var hasValidDiscount bool
	if stripeSubscription.Plan != nil {
		priceCents = ptr.Int64(stripeSubscription.Plan.Amount)
		hasValidDiscount = stripeSubscription.Discount != nil && stripeSubscription.Discount.Coupon != nil && stripeSubscription.Discount.Coupon.Valid
		if hasValidDiscount {
			priceCents = ptr.Int64(*priceCents - stripeSubscription.Discount.Coupon.AmountOff)
		}
	}
	premiumNewsletterSubscription := PremiumNewsletterSubscription{
		StripePaymentIntentID: paymentIntentID,
		CurrentPeriodEnd:      time.Unix(stripeSubscription.CurrentPeriodEnd, 0),
		IsAutoRenewEnabled:    !stripeSubscription.CancelAtPeriodEnd,
		PriceCents:            priceCents,
		HasValidDiscount:      hasValidDiscount,
	}
	var billingInformation *dbBillingInformation
	if dbNewsletterSubscription != nil {
		var err error
		billingInformation, err = getBillingInformation(tx, dbNewsletterSubscription.BillingInformationID)
		if err != nil {
			return nil, err
		}
		premiumNewsletterSubscription.userID = billingInformation.UserID
		premiumNewsletterSubscription.ID = &dbNewsletterSubscription.ID
	}
	switch stripeSubscription.Status {
	case stripe.SubscriptionStatusTrialing:
		premiumNewsletterSubscription.PaymentState = PaymentStateTrialNoPaymentMethod
		if billingInformation != nil && billingInformation.UserID != nil {
			paymentMethods, err := GetPaymentMethodsForUser(tx, *billingInformation.UserID)
			if err != nil {
				return nil, err
			}
			if len(paymentMethods) > 0 {
				premiumNewsletterSubscription.PaymentState = PaymentStateTrialPaymentMethodAdded
			}
		}
	case stripe.SubscriptionStatusIncomplete:
		premiumNewsletterSubscription.PaymentState = PaymentStateCreatedUnpaid
	case stripe.SubscriptionStatusActive:
		status := PaymentStateActive
		if stripeSubscription.LatestInvoice != nil && stripeSubscription.LatestInvoice.Status != stripe.InvoiceStatusPaid {
			status = PaymentStatePaymentPending
		}
		premiumNewsletterSubscription.PaymentState = status
	case stripe.SubscriptionStatusPastDue:
		premiumNewsletterSubscription.PaymentState = PaymentStatePaymentPending
	case stripe.SubscriptionStatusUnpaid:
		premiumNewsletterSubscription.PaymentState = PaymentStateErrored
	case stripe.SubscriptionStatusIncompleteExpired,
		stripe.SubscriptionStatusCanceled:
		premiumNewsletterSubscription.PaymentState = PaymentStateTerminated
	case stripe.SubscriptionStatusAll:
		return nil, fmt.Errorf("Unsupported payment status: all")
	default:
		return nil, fmt.Errorf("Unsupported payment status: %s", stripeSubscription.Status)
	}
	return &premiumNewsletterSubscription, nil
}

func HandleStripeEvent(c ctx.LogContext, tx *sqlx.Tx, stripeSignature string, eventBytes []byte) error {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	webhookSecret := env.MustEnvironmentVariable("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(eventBytes, stripeSignature, webhookSecret)
	if err != nil {
		return err
	}
	var wasProcessed bool
	switch event.Type {
	case "customer.subscription.created",
		"customer.subscription.deleted",
		"customer.subscription.pending_update_applied",
		"customer.subscription.updated":
		wasProcessed = true
		var stripeSubscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &stripeSubscription); err != nil {
			return err
		}
		premiumNewsletterSubscription, err := lookupPremiumNewsletterSubscriptionForStripeID(c, tx, stripeSubscription.ID)
		switch {
		case err != nil:
			return err
		case premiumNewsletterSubscription == nil:
			c.Warnf("No subscription found for stripe ID %s", stripeSubscription.ID)
		case premiumNewsletterSubscription != nil:
			if err := InsertPremiumNewsletterSyncRequest(tx, *premiumNewsletterSubscription.ID, PremiumNewsletterSubscriptionUpdateTypeRemoteUpdated); err != nil {
				return err
			}
		}
	case "invoice.paid",
		"invoice.payment_action_required",
		"invoice.payment_failed",
		"invoice.payment_succeeded":
		wasProcessed = true
		var stripeInvoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &stripeInvoice); err != nil {
			return err
		}
		if stripeInvoice.Subscription == nil {
			c.Warnf("Received invoiced %s with no subscription", stripeInvoice.ID)
		} else {
			premiumNewsletterSubscription, err := lookupPremiumNewsletterSubscriptionForStripeID(c, tx, stripeInvoice.Subscription.ID)
			switch {
			case err != nil:
				return err
			case premiumNewsletterSubscription == nil:
				c.Warnf("No subscription found for stripe ID %s", stripeInvoice.Subscription.ID)
			case premiumNewsletterSubscription != nil:
				if err := InsertPremiumNewsletterSyncRequest(tx, *premiumNewsletterSubscription.ID, PremiumNewsletterSubscriptionUpdateTypeRemoteUpdated); err != nil {
					return err
				}
			}
		}
	case "setup_intent.succeeded":
		wasProcessed = true
		var stripeSetupIntent stripe.SetupIntent
		if err := json.Unmarshal(event.Data.Raw, &stripeSetupIntent); err != nil {
			return err
		}
		if stripeSetupIntent.Customer != nil {
			billingInformation, err := LookupBillingInformationByExternalID(tx, stripeSetupIntent.Customer.ID)
			switch {
			case err != nil:
				return err
			case billingInformation == nil,
				billingInformation.UserID == nil:
				// no-op
			default:
				subscription, err := LookupPremiumNewsletterSubscriptionForUser(c, tx, *billingInformation.UserID)
				if err != nil {
					return err
				}
				if err := InsertPremiumNewsletterSyncRequest(tx, *subscription.ID, PremiumNewsletterSubscriptionUpdateTypePaymentMethodAdded); err != nil {
					return err
				}

			}
		}
	case "account.updated",
		"account.application.authorized",
		"account.application.deauthorized",
		"account.external_account.created",
		"account.external_account.deleted",
		"account.external_account.updated", "application_fee.created",
		"application_fee.refunded",
		"application_fee.refund.updated",
		"balance.available",
		"billing_portal.configuration.created",
		"billing_portal.configuration.updated",
		"capability.updated",
		"charge.captured",
		"charge.expired",
		"charge.failed",
		"charge.pending",
		"charge.refunded",
		"charge.succeeded",
		"charge.updated",
		"charge.dispute.closed",
		"charge.dispute.created",
		"charge.dispute.funds_reinstated",
		"charge.dispute.funds_withdrawn",
		"charge.dispute.updated",
		"charge.refund.updated",
		"checkout.session.async_payment_failed",
		"checkout.session.async_payment_succeeded",
		"checkout.session.completed",
		"checkout.session.expired",
		"coupon.created",
		"coupon.deleted",
		"coupon.updated",
		"credit_note.created",
		"credit_note.updated",
		"credit_note.voided",
		"customer.created",
		"customer.deleted",
		"customer.updated",
		"customer.discount.created",
		"customer.discount.deleted",
		"customer.discount.updated",
		"customer.source.created",
		"customer.source.deleted",
		"customer.source.expiring",
		"customer.source.updated",
		"customer.subscription.pending_update_expired",
		"customer.subscription.trial_will_end",
		"customer.tax_id.created",
		"customer.tax_id.deleted",
		"customer.tax_id.updated",
		"file.created",
		"identity.verification_session.canceled",
		"identity.verification_session.created",
		"identity.verification_session.processing",
		"identity.verification_session.redacted",
		"identity.verification_session.requires_input",
		"identity.verification_session.verified",
		"invoice.created",
		"invoice.deleted",
		"invoice.finalization_failed",
		"invoice.finalized",
		"invoice.marked_uncollectible",
		"invoice.sent",
		"invoice.upcoming",
		"invoice.updated",
		"invoice.voided",
		"invoiceitem.created",
		"invoiceitem.deleted",
		"invoiceitem.updated",
		"issuing_authorization.created",
		"issuing_authorization.request",
		"issuing_authorization.updated",
		"issuing_card.created",
		"issuing_card.updated",
		"issuing_cardholder.created",
		"issuing_cardholder.updated",
		"issuing_dispute.closed",
		"issuing_dispute.created",
		"issuing_dispute.funds_reinstated",
		"issuing_dispute.submitted",
		"issuing_dispute.updated",
		"issuing_transaction.created",
		"issuing_transaction.updated",
		"mandate.updated",
		"order.created",
		"order.payment_failed",
		"order.payment_succeeded",
		"order.updated",
		"order_return.created",
		"payment_intent.amount_capturable_updated",
		"payment_intent.canceled",
		"payment_intent.created",
		"payment_intent.payment_failed",
		"payment_intent.processing",
		"payment_intent.requires_action",
		"payment_intent.succeeded",
		"payment_link.created",
		"payment_link.updated",
		"payment_method.attached",
		"payment_method.automatically_updated",
		"payment_method.detached",
		"payment_method.updated",
		"payout.canceled",
		"payout.created",
		"payout.failed",
		"payout.paid",
		"payout.updated",
		"person.created",
		"person.deleted",
		"person.updated",
		"plan.created",
		"plan.deleted",
		"plan.updated",
		"price.created",
		"price.deleted",
		"price.updated",
		"product.created",
		"product.deleted",
		"product.updated",
		"promotion_code.created",
		"promotion_code.updated",
		"quote.accepted",
		"quote.canceled",
		"quote.created",
		"quote.finalized",
		"radar.early_fraud_warning.created",
		"radar.early_fraud_warning.updated",
		"recipient.created",
		"recipient.deleted",
		"recipient.updated",
		"reporting.report_run.failed",
		"reporting.report_run.succeeded",
		"reporting.report_type.updated",
		"review.closed",
		"review.opened",
		"setup_intent.canceled",
		"setup_intent.created",
		"setup_intent.requires_action",
		"setup_intent.setup_failed",
		"sigma.scheduled_query_run.created",
		"sku.created",
		"sku.deleted",
		"sku.updated",
		"source.canceled",
		"source.chargeable",
		"source.failed",
		"source.mandate_notification",
		"source.refund_attributes_required",
		"source.transaction.created",
		"source.transaction.updated",
		"subscription_schedule.aborted",
		"subscription_schedule.canceled",
		"subscription_schedule.completed",
		"subscription_schedule.created",
		"subscription_schedule.expiring",
		"subscription_schedule.released",
		"subscription_schedule.updated",
		"tax_rate.created",
		"tax_rate.updated",
		"topup.canceled",
		"topup.created",
		"topup.failed",
		"topup.reversed",
		"topup.succeeded",
		"transfer.created",
		"transfer.failed",
		"transfer.paid",
		"transfer.reversed",
		"transfer.updated":
		// no-op
	default:
		c.Warnf("Unrecognized event type: %s", event.Type)
	}
	_, err = tx.Exec(captureStripeEventQuery, event.Type, wasProcessed, eventBytes)
	return err
}
