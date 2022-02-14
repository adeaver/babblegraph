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
	lookupBillingInformationForUserIDQuery = "SELECT * FROM billing_information WHERE user_id = $1"
	insertBillingInformationForUserIDQuery = "INSERT INTO billing_information (user_id, external_id_mapping_id) VALUES ($1, $2)"

	getExternalIDMappingByIDQuery = "SELECT * FROM billing_external_id_mapping WHERE _id = $1"
	insertExternalIDMappingQuery  = "INSERT INTO billing_external_id_mapping (id_type, external_id) VALUES ($1, $2) RETURNING _id"
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
