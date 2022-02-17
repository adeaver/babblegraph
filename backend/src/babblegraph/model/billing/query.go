package billing

import (
	"babblegraph/config"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/env"
	"babblegraph/util/math/decimal"
	"babblegraph/util/ptr"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/sub"
)

const (
	getBillingInformationQuery             = "SELECT * FROM billing_information WHERE _id = $1"
	lookupBillingInformationForUserIDQuery = "SELECT * FROM billing_information WHERE user_id = $1"
	insertBillingInformationForUserIDQuery = "INSERT INTO billing_information (user_id, external_id_mapping_id) VALUES ($1, $2)"

	getExternalIDMappingByIDQuery = "SELECT * FROM billing_external_id_mapping WHERE _id = $1"
	insertExternalIDMappingQuery  = "INSERT INTO billing_external_id_mapping (id_type, external_id) VALUES ($1, $2) RETURNING _id"

	lookupPremiumNewsletterSubscriptionByIDQuery          = "SELECT * FROM billing_premium_newsletter_subscription WHERE _id = $1"
	lookupPremiumNewsletterSubscriptionQuery              = "SELECT * FROM billing_premium_newsletter_subscription WHERE billing_information_id = $1"
	lookupPremiumNewsletterNonTerminatedSubscriptionQuery = "SELECT * FROM billing_premium_newsletter_subscription WHERE billing_information_id = $1 AND is_terminated = FALSE"
	insertPremiumNewsletterSubscriptionQuery              = "INSERT INTO billing_premium_newsletter_subscription (_id, billing_information_id, external_id_mapping_id) VALUES ($1, $2, $3)"
	terminatePremiumNewsletterSubscriptionQuery           = "UPDATE billing_premium_newsletter_subscription SET is_terminated = TRUE WHERE _id = $1"

	insertPremiumNewsletterSubscriptionDebounceRecordQuery = "INSERT INTO billing_premium_newsletter_subscription_debounce_record (billing_information_id) VALUES ($1)"
	deletePremiumNewsletterSubscriptionDebounceRecordQuery = "DELETE FROM billing_premium_newsletter_subscription_debounce_record WHERE billing_information_id = $1"

	// TODO: add hold_until to these queries if need be.
	getAllPremiumNewsletterSyncRequestQuery = "SELECT * FROM billing_premium_newsletter_sync_request"
	// This model might need to change, but the current idea is that this should act as more of queue.
	// Where each subscription can only have whatever the latest update type is. As of right now, the only update type is that it makes a switch to active
	insertPremiumNewsletterSyncRequestQuery    = "INSERT INTO billing_premium_newsletter_sync_request VALUES (premium_newsletter_subscription_id, update_type) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	deletePremiumNewsletterSyncRequestQuery    = "DELETE FROM billing_premium_newsletter_sync_request WHERE premium_newsletter_subscription_id = $1"
	incrementPremiumNewsletterSyncRequestQuery = "UPDATE billing_premium_newsletter_sync_request SET attempt_number = attempt_number + 1 WHERE premium_newsletter_subscription_id = $1"
)

func GetOrCreateBillingInformationForUser(c ctx.LogContext, tx *sqlx.Tx, userID users.UserID) (*BillingInformation, error) {
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return nil, err
	case billingInformation != nil:
		externalID, err := getExternalIDMapping(tx, billingInformation.ExternalIDMappingID)
		if err != nil {
			return nil, err
		}
		out := &BillingInformation{
			UserID: &userID,
		}
		switch externalID.IDType {
		case externalIDTypeStripe:
			out.StripeCustomerID = ptr.String(externalID.ExternalID)
		default:
			return nil, fmt.Errorf("Unrecognized external ID type %s", externalID.IDType)
		}
		return out, nil
	default:
		stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
		user, err := users.GetUser(tx, userID)
		switch {
		case err != nil:
			return nil, err
		case user.Status != users.UserStatusVerified:
			return nil, fmt.Errorf("user is in the wrong state")
		}
		customerParams := &stripe.CustomerParams{
			Email: stripe.String(user.EmailAddress),
		}
		stripeCustomer, err := customer.New(customerParams)
		if err != nil {
			return nil, err
		}
		externalMappingID, err := insertExternalIDMapping(tx, stripeCustomer.ID)
		if err != nil {
			c.Warnf("Attempting to rollback customer with Stripe ID %s and Babblegraph User ID %s", stripeCustomer.ID, userID)
			if _, sErr := customer.Del(stripeCustomer.ID, &stripe.CustomerParams{}); sErr != nil {
				c.Errorf("Error rolling back customer ID %s in Stripe for user ID %s because of error %s", stripeCustomer.ID, userID, sErr.Error())
			}
			return nil, err
		}
		if err := insertBillingInformationForUserID(tx, userID, *externalMappingID); err != nil {
			c.Warnf("Attempting to rollback customer with Stripe ID %s and Babblegraph User ID %s", stripeCustomer.ID, userID)
			if _, sErr := customer.Del(stripeCustomer.ID, &stripe.CustomerParams{}); sErr != nil {
				c.Errorf("Error rolling back customer ID %s in Stripe for user ID %s because of error %s", stripeCustomer.ID, userID, sErr.Error())
			}
			return nil, err
		}
		return &BillingInformation{
			UserID:           &userID,
			StripeCustomerID: ptr.String(stripeCustomer.ID),
		}, nil
	}
}

func getBillingInformation(tx *sqlx.Tx, id BillingInformationID) (*dbBillingInformation, error) {
	var matches []dbBillingInformation
	err := tx.Select(&matches, getBillingInformationQuery, id)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0,
		len(matches) > 1:
		return nil, fmt.Errorf("Expected exactly one match for billing information id %s, but got %d", id, len(matches))
	default:
		return &matches[0], nil
	}
}

func lookupBillingInformationForUserID(tx *sqlx.Tx, userID users.UserID) (*dbBillingInformation, error) {
	var matches []dbBillingInformation
	err := tx.Select(&matches, lookupBillingInformationForUserIDQuery, userID)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0:
		return nil, nil
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most one billing information for user ID %s, but got %d", userID, len(matches))
	}
	return &matches[0], nil
}

func insertBillingInformationForUserID(tx *sqlx.Tx, userID users.UserID, externalID externalIDMappingID) error {
	if _, err := tx.Exec(insertBillingInformationForUserIDQuery, userID, externalID); err != nil {
		return err
	}
	return nil
}

func LookupPremiumNewsletterSubscriptionForUser(c ctx.LogContext, tx *sqlx.Tx, userID users.UserID) (*PremiumNewsletterSubscription, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return nil, err
	case billingInformation == nil:
		return nil, fmt.Errorf("Expected there to be a billing information for user %s, but none exists", userID)
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

func GetPremiumNewsletterSubscriptionTrialEligibilityForUser(tx *sqlx.Tx, userID users.UserID) (*int64, error) {
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return nil, err
	case billingInformation == nil:
		return ptr.Int64(config.PremiumNewsletterSubscriptionTrialLengthDays), nil
	default:
		var matches []dbPremiumNewsletterSubscription
		err := tx.Select(&matches, lookupPremiumNewsletterSubscriptionQuery, billingInformation.ID)
		switch {
		case err != nil:
			return nil, err
		case len(matches) == 0:
			return ptr.Int64(config.PremiumNewsletterSubscriptionTrialLengthDays), nil
		default:
			var oldestMatch *time.Time
			for _, m := range matches {
				if oldestMatch == nil || oldestMatch.After(m.CreatedAt) {
					oldestMatch = &m.CreatedAt
				}
			}
			hoursSinceOldestTrialStarted := decimal.FromInt64(int64(time.Now().Sub(*oldestMatch) / time.Hour))
			roundedDaysSinceOldestTrialStarted := hoursSinceOldestTrialStarted.Divide(decimal.FromInt64(24)).ToInt64Rounded()
			return &roundedDaysSinceOldestTrialStarted, nil
		}
	}
}

func InsertPremiumNewsletterSyncRequest(tx *sqlx.Tx, id PremiumNewsletterSubscriptionID, updateType PremiumNewsletterSubscriptionUpdateType) error {
	if _, err := tx.Exec(insertPremiumNewsletterSyncRequestQuery, id, updateType); err != nil {
		return err
	}
	return nil
}

func GetPremiumNewsletterSyncRequests(tx *sqlx.Tx) (map[PremiumNewsletterSubscriptionID]PremiumNewsletterSubscriptionUpdateType, error) {
	var matches []dbPremiumNewsletterSubscriptionSyncRequest
	if err := tx.Select(&matches, getAllPremiumNewsletterSyncRequestQuery); err != nil {
		return nil, err
	}
	out := make(map[PremiumNewsletterSubscriptionID]PremiumNewsletterSubscriptionUpdateType)
	for _, m := range matches {
		out[m.PremiumNewsletterSubscriptionID] = m.UpdateType
	}
	return out, nil
}

func MarkPremiumNewsletterSyncRequestDone(tx *sqlx.Tx, id PremiumNewsletterSubscriptionID) error {
	if _, err := tx.Exec(deletePremiumNewsletterSyncRequestQuery, id); err != nil {
		return err
	}
	return nil
}

func MarkPremiumNewsletterSyncRequestForRetry(tx *sqlx.Tx, id PremiumNewsletterSubscriptionID) error {
	if _, err := tx.Exec(incrementPremiumNewsletterSyncRequestQuery, id); err != nil {
		return err
	}
	return nil
}

func getExternalIDMapping(tx *sqlx.Tx, id externalIDMappingID) (*dbExternalIDMapping, error) {
	var matches []dbExternalIDMapping
	err := tx.Select(&matches, getExternalIDMappingByIDQuery, id)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0,
		len(matches) > 1:
		return nil, fmt.Errorf("Expected exactly one external id mapping for id %s, but got %d", id, len(matches))
	default:
		return &matches[0], nil
	}
}

func insertExternalIDMapping(tx *sqlx.Tx, externalID string) (*externalIDMappingID, error) {
	rows, err := tx.Query(insertExternalIDMappingQuery, externalIDTypeStripe, externalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id externalIDMappingID
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
	}
	return &id, nil
}
