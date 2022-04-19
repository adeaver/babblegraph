package billing

import (
	"babblegraph/config"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/env"
	"babblegraph/util/math/decimal"
	"babblegraph/util/ptr"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
)

const (
	lookupPremiumNewsletterSubscriptionByIDQuery          = "SELECT * FROM billing_premium_newsletter_subscription WHERE _id = $1"
	lookupPremiumNewsletterSubscriptionQuery              = "SELECT * FROM billing_premium_newsletter_subscription WHERE billing_information_id = $1 ORDER BY created_at DESC"
	lookupPremiumNewsletterSubscriptionByExternalIDQuery  = "SELECT * FROM billing_premium_newsletter_subscription WHERE external_id_mapping_id = $1"
	lookupPremiumNewsletterNonTerminatedSubscriptionQuery = "SELECT * FROM billing_premium_newsletter_subscription WHERE billing_information_id = $1 AND is_terminated = FALSE"
	insertPremiumNewsletterSubscriptionQuery              = "INSERT INTO billing_premium_newsletter_subscription (_id, billing_information_id, external_id_mapping_id) VALUES ($1, $2, $3)"
	terminatePremiumNewsletterSubscriptionQuery           = "UPDATE billing_premium_newsletter_subscription SET is_terminated = TRUE, last_modified_at = timezone('utc', now()) WHERE _id = $1"

	insertPremiumNewsletterSubscriptionDebounceRecordQuery = "INSERT INTO billing_premium_newsletter_subscription_debounce_record (billing_information_id) VALUES ($1)"
	deletePremiumNewsletterSubscriptionDebounceRecordQuery = "DELETE FROM billing_premium_newsletter_subscription_debounce_record WHERE billing_information_id = $1"

	lookupNewsletterTrialForEmailAddressQuery = "SELECT * FROM billing_newsletter_subscription_trials WHERE email_address = $1 AND created_at > current_date - interval '2 years'"
	insertTrialRecordForUserQuery             = "INSERT INTO billing_newsletter_subscription_trials (email_address) VALUES ($1)"
)

func LookupPremiumNewsletterSubscriptionForUser(c ctx.LogContext, tx *sqlx.Tx, userID users.UserID) (*PremiumNewsletterSubscription, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return nil, err
	case billingInformation == nil:
		return nil, nil
	}
	return lookupActivePremiumNewsletterSubscriptionForUser(c, tx, *billingInformation)
}

func CreatePremiumNewsletterSubscriptionForUserWithID(c ctx.LogContext, tx *sqlx.Tx, userID users.UserID, subscriptionID PremiumNewsletterSubscriptionID) (*PremiumNewsletterSubscription, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return nil, err
	case billingInformation == nil:
		return nil, fmt.Errorf("Expected there to be a billing information for user %s, but none exists", userID)
	}
	if _, err := tx.Exec(insertPremiumNewsletterSubscriptionDebounceRecordQuery, billingInformation.ID); err != nil {
		return nil, err
	}
	stripeProductID, err := getStripeProductIDForEnvironment()
	if err != nil {
		return nil, err
	}
	subscriptionParams := &stripe.SubscriptionParams{
		PaymentBehavior: stripe.String("default_incomplete"),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripeProductID,
			},
		},
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	subscriptionParams.AddExpand("default_payment_method")
	externalID, err := getExternalIDMapping(tx, billingInformation.ExternalIDMappingID)
	if err != nil {
		return nil, err
	}
	switch externalID.IDType {
	case externalIDTypeStripe:
		subscriptionParams.Customer = ptr.String(externalID.ExternalID)
	default:
		return nil, fmt.Errorf("Unrecognized external ID type %s", externalID.IDType)
	}
	// There is no active subscription, but there may have been a previous one, so we need
	// to check the trial eligibility
	trialEligibilityDays, err := GetPremiumNewsletterSubscriptionTrialEligibilityForUser(tx, userID)
	if err != nil {
		return nil, err
	}
	if trialEligibilityDays != nil && *trialEligibilityDays > 0 {
		subscriptionParams.TrialPeriodDays = trialEligibilityDays
		if err := insertTrialRecordForUser(tx, userID); err != nil {
			return nil, err
		}
	}
	stripeSubscription, err := sub.New(subscriptionParams)
	if err != nil {
		return nil, err
	}
	if err := insertActivePremiumNewsletterSubscriptionForUser(tx, subscriptionID, billingInformation.ID, stripeSubscription); err != nil {
		c.Warnf("Attempting to rollback stripe subscription with Stripe ID %s and Babblegraph User ID %s", stripeSubscription.ID, userID)
		if _, sErr := sub.Cancel(stripeSubscription.ID, &stripe.SubscriptionCancelParams{}); sErr != nil {
			c.Errorf("Error rolling back subscription ID %s in Stripe for user ID %s because of error %s", stripeSubscription.ID, userID, sErr.Error())
		}
		return nil, err
	}
	return convertStripeSubscriptionToPremiumNewsletterSubscription(tx, stripeSubscription, nil)
}

func GetPremiumNewsletterSubscriptionByID(c ctx.LogContext, tx *sqlx.Tx, id PremiumNewsletterSubscriptionID) (*PremiumNewsletterSubscription, error) {
	var matches []dbPremiumNewsletterSubscription
	err := tx.Select(&matches, lookupPremiumNewsletterSubscriptionByIDQuery, id)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0,
		len(matches) > 1:
		return nil, fmt.Errorf("Expecting exactly one premium newsletter subscription for id %s, but got %d", id, len(matches))
	default:
		return getStripeSubscriptionAndConvertSubscriptionForDBPremiumNewsletterSubscription(c, tx, matches[0], true)
	}
}

func UpdateSubscriptionAutoRenewForUser(tx *sqlx.Tx, userID users.UserID, isAutoRenewEnabled bool) error {
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return err
	case billingInformation == nil:
		return fmt.Errorf("Expected there to be a billing information for user %s, but none exists", userID)
	}
	dbPremiumNewsletterSubscription, err := lookupDBActivePremiumNewsletterSubscriptionForUser(tx, *billingInformation)
	switch {
	case err != nil:
		return err
	case dbPremiumNewsletterSubscription == nil:
		return nil
	case dbPremiumNewsletterSubscription != nil:
		externalID, err := getExternalIDMapping(tx, dbPremiumNewsletterSubscription.ExternalIDMappingID)
		if err != nil {
			return err
		}
		switch externalID.IDType {
		case externalIDTypeStripe:
			stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
			if _, err := sub.Update(externalID.ExternalID, &stripe.SubscriptionParams{
				CancelAtPeriodEnd: ptr.Bool(!isAutoRenewEnabled),
			}); err != nil {
				return err
			}
			return nil
		default:
			return fmt.Errorf("Invalid ID Type %s", externalID.IDType)
		}
	default:
		panic("unreachable")
	}
}

func CancelPremiumNewsletterSubscriptionForUser(c ctx.LogContext, tx *sqlx.Tx, userID users.UserID) error {
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return err
	case billingInformation == nil:
		return fmt.Errorf("Expected there to be a billing information for user %s, but none exists", userID)
	}
	dbPremiumNewsletterSubscription, err := lookupDBActivePremiumNewsletterSubscriptionForUser(tx, *billingInformation)
	switch {
	case err != nil:
		return err
	case dbPremiumNewsletterSubscription == nil:
		return nil
	case dbPremiumNewsletterSubscription != nil:
		if _, err := tx.Exec(terminatePremiumNewsletterSubscriptionQuery, dbPremiumNewsletterSubscription.ID); err != nil {
			return err
		}
		externalID, err := getExternalIDMapping(tx, dbPremiumNewsletterSubscription.ExternalIDMappingID)
		if err != nil {
			return err
		}
		switch externalID.IDType {
		case externalIDTypeStripe:
			stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
			if _, err := sub.Cancel(externalID.ExternalID, nil); err != nil {
				return err
			}
			return nil
		default:
			return fmt.Errorf("Invalid ID Type %s", externalID.IDType)
		}
	default:
		panic("unreachable")
	}
}

func lookupActivePremiumNewsletterSubscriptionForUser(c ctx.LogContext, tx *sqlx.Tx, billingInformation dbBillingInformation) (*PremiumNewsletterSubscription, error) {
	// There are three possible scenarios for this function:
	// The database returns no active subscriptions - in which case we assume that there are no active subscriptions in the provider
	// The database returns an active subscription, but the payment provider returns a non-active subscription
	// --> in this case we need to delete any debounce records and update the db record, this function will return nil, nil
	// The database returns an active subscription, which maps to an active subscription in the provider, in which case we're good
	dbPremiumNewsletterSubscription, err := lookupDBActivePremiumNewsletterSubscriptionForUser(tx, billingInformation)
	switch {
	case err != nil:
		return nil, err
	case dbPremiumNewsletterSubscription == nil:
		return nil, nil
	case dbPremiumNewsletterSubscription != nil:
		return getStripeSubscriptionAndConvertSubscriptionForDBPremiumNewsletterSubscription(c, tx, *dbPremiumNewsletterSubscription, false)
	default:
		panic("unreachable")
	}
}

func getStripeSubscriptionAndConvertSubscriptionForDBPremiumNewsletterSubscription(c ctx.LogContext, tx *sqlx.Tx, dbPremiumNewsletterSubscription dbPremiumNewsletterSubscription, shouldReturnIfTerminated bool) (*PremiumNewsletterSubscription, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	var premiumNewsletterSubscription *PremiumNewsletterSubscription
	externalID, err := getExternalIDMapping(tx, dbPremiumNewsletterSubscription.ExternalIDMappingID)
	if err != nil {
		return nil, err
	}
	switch externalID.IDType {
	case externalIDTypeStripe:
		subscriptionParams := &stripe.SubscriptionParams{}
		subscriptionParams.AddExpand("latest_invoice.payment_intent")
		subscriptionParams.AddExpand("default_payment_method")
		stripeSubscription, err := sub.Get(externalID.ExternalID, subscriptionParams)
		if err != nil {
			return nil, err
		}
		premiumNewsletterSubscription, err = convertStripeSubscriptionToPremiumNewsletterSubscription(tx, stripeSubscription, &dbPremiumNewsletterSubscription)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unrecognized external ID type %s", externalID.IDType)
	}
	switch premiumNewsletterSubscription.PaymentState {
	case PaymentStateCreatedUnpaid,
		PaymentStateTrialNoPaymentMethod,
		PaymentStateTrialPaymentMethodAdded,
		PaymentStateActive,
		PaymentStateErrored:
		return premiumNewsletterSubscription, nil
	case PaymentStateTerminated:
		c.Infof("Provider subscription with external ID mapping %s is terminated, but corresponding Babblegraph subscription %s is not. Resolving...", dbPremiumNewsletterSubscription.ExternalIDMappingID, dbPremiumNewsletterSubscription.ID)
		if _, err := tx.Exec(terminatePremiumNewsletterSubscriptionQuery, dbPremiumNewsletterSubscription.ID); err != nil {
			return nil, err
		}
		if _, err := tx.Exec(deletePremiumNewsletterSubscriptionDebounceRecordQuery, dbPremiumNewsletterSubscription.BillingInformationID); err != nil {
			return nil, err
		}
		if shouldReturnIfTerminated {
			return premiumNewsletterSubscription, nil
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("Unrecognized payment state %d", premiumNewsletterSubscription.PaymentState)
	}
}

func lookupDBActivePremiumNewsletterSubscriptionForUser(tx *sqlx.Tx, billingInformation dbBillingInformation) (*dbPremiumNewsletterSubscription, error) {
	var matches []dbPremiumNewsletterSubscription
	err := tx.Select(&matches, lookupPremiumNewsletterNonTerminatedSubscriptionQuery, billingInformation.ID)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0:
		return nil, nil
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most one active premium newsletter subscription for user %s but got %d", *billingInformation.UserID, len(matches))
	default:
		return &matches[0], nil
	}
}

func insertActivePremiumNewsletterSubscriptionForUser(tx *sqlx.Tx, premiumNewsletterSubscriptionID PremiumNewsletterSubscriptionID, billingInformationID BillingInformationID, stripeSubscription *stripe.Subscription) error {
	externalIDMappingID, err := insertExternalIDMapping(tx, stripeSubscription.ID)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(insertPremiumNewsletterSubscriptionQuery, premiumNewsletterSubscriptionID, billingInformationID, externalIDMappingID); err != nil {
		return err
	}
	return nil
}

func insertTrialRecordForUser(tx *sqlx.Tx, userID users.UserID) error {
	user, err := users.GetUser(tx, userID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(insertTrialRecordForUserQuery, user.EmailAddress)
	return err
}

func GetPremiumNewsletterSubscriptionTrialEligibilityForUser(tx *sqlx.Tx, userID users.UserID) (*int64, error) {
	subscriptionLevel, err := useraccounts.LookupSubscriptionLevelForUser(tx, userID)
	switch {
	case err != nil:
		return nil, err
	case subscriptionLevel == nil:
		user, err := users.GetUser(tx, userID)
		if err != nil {
			return nil, err
		}
		return getNewsletterSubscriptionTrialEligibility(tx, user.EmailAddress)
	case *subscriptionLevel == useraccounts.SubscriptionLevelLegacy,
		*subscriptionLevel == useraccounts.SubscriptionLevelPremium,
		*subscriptionLevel == useraccounts.SubscriptionLevelBetaPremium:
		return ptr.Int64(0), nil
	default:
		return nil, fmt.Errorf("Unreachable")
	}
}

func getNewsletterSubscriptionTrialEligibility(tx *sqlx.Tx, emailAddress string) (*int64, error) {
	var matches []dbNewsletterSubscriptionTrial
	if err := tx.Select(&matches, lookupNewsletterTrialForEmailAddressQuery, emailAddress); err != nil {
		return nil, err
	}
	trialEligibilityDays := config.PremiumNewsletterSubscriptionTrialLengthDays
	for _, m := range matches {
		lengthOfPreviousTrialInHours := decimal.FromInt64(int64(time.Now().Sub(m.CreatedAt) / time.Hour))
		lengthOfPreviousTrialsInDaysRounded := lengthOfPreviousTrialInHours.Divide(decimal.FromInt64(24)).ToInt64Rounded()
		trialEligibilityDays -= lengthOfPreviousTrialsInDaysRounded
	}
	if trialEligibilityDays < 0 {
		trialEligibilityDays = 0
	}
	return ptr.Int64(trialEligibilityDays), nil
}

// Webhook functions

func lookupPremiumNewsletterSubscriptionForStripeID(c ctx.LogContext, tx *sqlx.Tx, stripeID string) (*PremiumNewsletterSubscription, error) {
	externalIDMapping, err := lookupExternalIDMappingByExternalID(tx, externalIDTypeStripe, stripeID)
	switch {
	case err != nil:
		return nil, err
	case externalIDMapping == nil:
		return nil, nil
	default:
		externalIDMappingID := externalIDMapping.ID
		var matches []dbPremiumNewsletterSubscription
		err := tx.Select(&matches, lookupPremiumNewsletterSubscriptionByExternalIDQuery, externalIDMappingID)
		switch {
		case err != nil:
			return nil, err
		case len(matches) == 0:
			return nil, nil
		case len(matches) > 1:
			return nil, fmt.Errorf("Expected at most one subscription for external ID %s, but got %d", externalIDMappingID, len(matches))
		default:
			return getStripeSubscriptionAndConvertSubscriptionForDBPremiumNewsletterSubscription(c, tx, matches[0], true)
		}
	}
}
