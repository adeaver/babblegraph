package billing

import (
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
)

const (
	getBillingInformationQuery                 = "SELECT * FROM billing_information WHERE _id = $1"
	lookupBillingInformationForUserIDQuery     = "SELECT * FROM billing_information WHERE user_id = $1"
	lookupBillingInformationForExternalIDQuery = "SELECT * FROM billing_information WHERE external_id_mapping_id = $1"
	insertBillingInformationForUserIDQuery     = "INSERT INTO billing_information (user_id, external_id_mapping_id) VALUES ($1, $2)"
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

// TODO: maybe external ID type should be exported
func LookupBillingInformationByExternalID(tx *sqlx.Tx, externalID string) (*BillingInformation, error) {
	externalIDMapping, err := lookupExternalIDMappingByExternalID(tx, externalIDTypeStripe, externalID)
	switch {
	case err != nil:
		return nil, err
	case externalIDMapping != nil:
		var matches []dbBillingInformation
		err = tx.Select(&matches, lookupBillingInformationForExternalIDQuery, externalIDMapping.ID)
		switch {
		case err != nil:
			return nil, err
		case len(matches) == 0:
			return nil, nil
		case len(matches) > 1:
			return nil, fmt.Errorf("Expected at most one billing information for external id %s, but got %d", externalIDMapping.ID, len(matches))
		default:
			return &BillingInformation{
				UserID:           matches[0].UserID,
				StripeCustomerID: ptr.String(externalID),
			}, nil
		}
	default:
		// This is where we'd try another ID Type
		return nil, nil
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
