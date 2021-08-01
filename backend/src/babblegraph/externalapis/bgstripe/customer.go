package bgstripe

import (
	"babblegraph/model/users"
	"babblegraph/util/env"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
)

const (
	getStripeCustomerForUserQuery = "SELECT * FROM bgstripe_customer WHERE babblegraph_user_id = $1"
	createCustomerForUserQuery    = "INSERT INTO bgstripe_customer (babblegraph_user_id, stripe_customer_id) VALUES ($1, $2)"
)

func CreateCustomerForUser(tx *sqlx.Tx, userID users.UserID) (*CustomerID, error) {
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
	if _, err := tx.Exec(createCustomerForUserQuery, userID, stripeCustomer.ID); err != nil {
		log.Println(fmt.Sprintf("Attempting to roll back customer: %s", stripeCustomer.ID))
		if _, sErr := customer.Del(stripeCustomer.ID, &stripe.CustomerParams{}); sErr != nil {
			formattedSErr := fmt.Errorf("Error rolling back customer ID %s in Stripe, for user %s: %s. Original error: %s", stripeCustomer.ID, userID, sErr, err)
			log.Println(formattedSErr.Error())
			sentry.CaptureException(formattedSErr)
		}
		return nil, err
	}
	asCustomerID := CustomerID(stripeCustomer.ID)
	return &asCustomerID, nil
}

func getStripeCustomerForUserID(tx *sqlx.Tx, userID users.UserID) (*dbStripeCustomer, error) {
	var matches []dbStripeCustomer
	if err := tx.Select(&matches, getStripeCustomerForUserQuery, userID); err != nil {
		return nil, err
	}
	switch {
	case len(matches) == 0:
		return nil, fmt.Errorf("no matches found for user ID: %s", userID)
	case len(matches) == 1:
		m := matches[0]
		return &m, nil
	default:
		return nil, fmt.Errorf("expected 1 stripe customer match for user ID %s, but got %d", userID, len(matches))
	}
}
