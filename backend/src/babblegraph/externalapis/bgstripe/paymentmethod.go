package bgstripe

import (
	"babblegraph/model/users"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
)

const (
	insertPaymentMethodQuery = "INSERT INTO bgstripe_payment_method (babblegraph_user_id, stripe_payment_method_id, card_type, last_four_digits, expiration_month, expiration_year) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING"

	getPaymentMethodForUserQuery = "SELECT * FROM bgstripe_payment_method WHERE babblegraph_user_id = $1"
)

func InsertPaymentMethod(tx *sqlx.Tx, userID users.UserID, paymentMethod stripe.PaymentMethod) error {
	if paymentMethod.Card == nil {
		return nil
	}
	if _, err := tx.Exec(insertPaymentMethodQuery, userID, paymentMethod.ID, paymentMethod.Card.Brand, paymentMethod.Card.Last4, paymentMethod.Card.ExpMonth, paymentMethod.Card.ExpYear); err != nil {
		return err
	}
	return nil
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
