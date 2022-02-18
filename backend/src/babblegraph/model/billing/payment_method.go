package billing

import (
	"babblegraph/model/users"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentmethod"
)

func GetPaymentMethodsForUser(tx *sqlx.Tx, userID users.UserID) ([]PaymentMethod, error) {
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
		switch externalID.IDType {
		case externalIDTypeStripe:
			return getStripePaymentMethodsForUser(tx, externalID.ExternalID)
		default:
			return nil, fmt.Errorf("Unrecognized external ID type %s", externalID.IDType)
		}
	}
}

// TODO: delete payment method
// TODO: update default payment method
// --> also needs a getDefaultPaymentMethodID for user, which should actually get Customer

func getStripePaymentMethodsForUser(tx *sqlx.Tx, stripeCustomerID string) ([]PaymentMethod, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	stripePaymentMethods := paymentmethod.List(&stripe.PaymentMethodListParams{
		Customer: ptr.String(stripeCustomerID),
	})
	var out []PaymentMethod
	for _, paymentMethod := range stripePaymentMethods.PaymentMethodList().Data {
		convertedPaymentMethod, err := convertStripePaymentMethod(paymentMethod)
		if err != nil {
			return nil, err
		}
		out = append(out, *convertedPaymentMethod)
	}
	return out, nil
}

func convertStripePaymentMethod(stripePaymentMethod *stripe.PaymentMethod) (*PaymentMethod, error) {
	if stripePaymentMethod.Card == nil {
		return nil, fmt.Errorf("Card is null, but expected only card type")
	}
	return &PaymentMethod{
		ExternalID:     stripePaymentMethod.ID,
		DisplayMask:    stripePaymentMethod.Card.Last4,
		CardExpiration: fmt.Sprintf("%2d/%d", stripePaymentMethod.Card.ExpMonth, stripePaymentMethod.Card.ExpYear),
		CardType:       getCardTypeForStripeBrand(stripePaymentMethod.Card.Brand),
	}, nil
}

func getCardTypeForStripeBrand(brand stripe.PaymentMethodCardBrand) CardType {
	switch brand {
	case stripe.PaymentMethodCardBrandAmex:
		return CardTypeAmex
	case stripe.PaymentMethodCardBrandDiscover:
		return CardTypeDiscover
	case stripe.PaymentMethodCardBrandMastercard:
		return CardTypeMastercard
	case stripe.PaymentMethodCardBrandVisa:
		return CardTypeVisa
	default:
		return CardTypeOther
	}
}
