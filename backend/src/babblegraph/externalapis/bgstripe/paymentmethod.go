package bgstripe

import (
	"babblegraph/model/users"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/setupintent"
)

const (
	insertPaymentMethodQuery = "INSERT INTO bgstripe_payment_method (babblegraph_user_id, stripe_payment_method_id, card_type, last_four_digits, expiration_month, expiration_year) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING"

	getPaymentMethodByIDQuery    = "SELECT * FROM bgstripe_payment_method WHERE stripe_payment_method_id = $1"
	getPaymentMethodForUserQuery = "SELECT * FROM bgstripe_payment_method WHERE babblegraph_user_id = $1"
)

type AddPaymentMethodCredentials struct {
	ClientSecret  string
	SetupIntentID string
}

func GetAddPaymentMethodCredentialsForUser(tx *sqlx.Tx, userID users.UserID) (*AddPaymentMethodCredentials, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	customer, err := getStripeCustomerForUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	params := &stripe.SetupIntentParams{
		Confirm:  ptr.Bool(true),
		Customer: stripe.String(string(customer.StripeCustomerID)),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		Usage: ptr.String("off_session"),
	}
	si, err := setupintent.New(params)
	if err != nil {
		return nil, err
	}
	return &AddPaymentMethodCredentials{
		ClientSecret:  si.ClientSecret,
		SetupIntentID: si.ID,
	}, nil
}

func FindStripePaymentMethodAndInsert(tx *sqlx.Tx, userID users.UserID, stripePaymentMethodID PaymentMethodID) (*PaymentMethod, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	paymentMethod, err := paymentmethod.Get(string(stripePaymentMethodID), nil)
	if err != nil {
		return nil, err
	}
	if err := InsertPaymentMethod(tx, userID, paymentMethod); err != nil {
		return nil, err
	}
	paymentMethod, err := LookupPaymentMethod(tx, stripePaymentMethodID)
	switch {
	case err != nil:
		return nil, err
	case paymentMethod == nil:
		return nil, fmt.Errorf("No payment method found")
	default:
		return paymentMethod, nil
	}

}

func InsertPaymentMethod(tx *sqlx.Tx, userID users.UserID, paymentMethod stripe.PaymentMethod) error {
	if paymentMethod.Card == nil {
		return nil
	}
	if paymentMethod.Customer == nil {
		return fmt.Errorf("Payment method not associated with customer")
	}
	if err := verifyCustomerIDForUser(tx, userID, CustomerID(paymentMethod.Customer.ID)); err != nil {
		return err
	}
	if _, err := tx.Exec(insertPaymentMethodQuery, userID, paymentMethod.ID, paymentMethod.Card.Brand, paymentMethod.Card.Last4, paymentMethod.Card.ExpMonth, paymentMethod.Card.ExpYear); err != nil {
		return err
	}
	return nil
}

func LookupPaymentMethod(tx *sqlx.Tx, paymentMethodID PaymentMethodID) (*PaymentMethod, error) {
	var matches []dbStripePaymentMethod
	if err := tx.Select(&matches, getPaymentMethodByIDQuery, paymentMethodID); err != nil {
		return nil, err
	}
	switch {
	case len(matches) == 0:
		return nil, nil
	case len(matches) == 1:
		m := matches[0]
		customer, err := getStripeCustomerForUserID(tx, m.BabblegraphUserID)
		if err != nil {
			return nil, err
		}
		return &PaymentMethod{
			StripePaymentMethodID: m.StripePaymentMethodID,
			CardType:              m.CardType,
			LastFourDigits:        m.LastFourDigits,
			ExpirationMonth:       m.ExpirationMonth,
			ExpirationYear:        m.ExpirationYear,
			IsDefault:             customer.DefaultPaymentMethodID == nil && *customer.DefaultPaymentMethodID == m.StripePaymentMethodID,
		}, nil
	default:
		return nil, fmt.Errorf("Expected at most 1 match, but got %d", len(matches))
	}
}

func GetPaymentMethodsForUser(tx *sqlx.Tx, userID users.UserID) ([]PaymentMethod, error) {
	var matches []dbStripePaymentMethod
	if err := tx.Select(&matches, getPaymentMethodForUserQuery, userID); err != nil {
		return nil, err
	}
	customer, err := getStripeCustomerForUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	var out []PaymentMethod
	for _, m := range matches {
		out = append(out, PaymentMethod{
			StripePaymentMethodID: m.StripePaymentMethodID,
			CardType:              m.CardType,
			LastFourDigits:        m.LastFourDigits,
			ExpirationMonth:       m.ExpirationMonth,
			ExpirationYear:        m.ExpirationYear,
			IsDefault:             customer.DefaultPaymentMethodID == nil && *customer.DefaultPaymentMethodID == m.StripePaymentMethodID,
		})
	}
	return out, nil
}

func verifyCustomerIDForUser(tx *sqlx.Tx, userID users.UserID, customerID CustomerID) error {
	customer, err := getStripeCustomerForUserID(tx, userID)
	if err != nil {
		return err
	}
	if customer.StripeCustomerID != customerID {
		return fmt.Errorf("Incorrect customer")
	}
	return nil
}
