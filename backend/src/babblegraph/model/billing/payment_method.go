package billing

import (
	"babblegraph/model/users"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
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
			return getStripePaymentMethodsForUser(externalID.ExternalID)
		default:
			return nil, fmt.Errorf("Unrecognized external ID type %s", externalID.IDType)
		}
	}
}

func MarkPaymentMethodAsDefaultForUser(tx *sqlx.Tx, userID users.UserID, externalPaymentMethodID string) error {
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return err
	case billingInformation == nil:
		return fmt.Errorf("Expected billing information for user %s, but got none", userID)
	default:
		externalID, err := getExternalIDMapping(tx, billingInformation.ExternalIDMappingID)
		if err != nil {
			return err
		}
		switch externalID.IDType {
		case externalIDTypeStripe:
			return markPaymentMethodAsDefaultForUser(externalID.ExternalID, externalPaymentMethodID)
		default:
			return fmt.Errorf("Unrecognized external ID type %s", externalID.IDType)
		}
	}
}

func DeletePaymentMethodForUser(tx *sqlx.Tx, userID users.UserID, externalPaymentMethodID string) error {
	billingInformation, err := lookupBillingInformationForUserID(tx, userID)
	switch {
	case err != nil:
		return err
	case billingInformation == nil:
		return fmt.Errorf("Expected billing information for user %s, but got none", userID)
	default:
		externalID, err := getExternalIDMapping(tx, billingInformation.ExternalIDMappingID)
		if err != nil {
			return err
		}
		switch externalID.IDType {
		case externalIDTypeStripe:
			return deletePaymentMethod(externalID.ExternalID, externalPaymentMethodID)
		default:
			return fmt.Errorf("Unrecognized external ID type %s", externalID.IDType)
		}
	}
}

func getStripeDefaultPaymentMethodID(stripeCustomerID string) (*string, error) {
	customer, err := customer.Get(stripeCustomerID, &stripe.CustomerParams{})
	switch {
	case err != nil:
		return nil, err
	case customer.InvoiceSettings == nil:
		return nil, nil
	case customer.InvoiceSettings.DefaultPaymentMethod == nil:
		return nil, nil
	default:
		return ptr.String(customer.InvoiceSettings.DefaultPaymentMethod.ID), nil
	}
}

func getStripePaymentMethodsForUser(stripeCustomerID string) ([]PaymentMethod, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	stripePaymentMethods := paymentmethod.List(&stripe.PaymentMethodListParams{
		Customer: ptr.String(stripeCustomerID),
		Type:     ptr.String("card"),
	})
	defaultPaymentMethodID, err := getStripeDefaultPaymentMethodID(stripeCustomerID)
	if err != nil {
		return nil, err
	}
	var out []PaymentMethod
	for _, paymentMethod := range stripePaymentMethods.PaymentMethodList().Data {
		convertedPaymentMethod, err := convertStripePaymentMethod(defaultPaymentMethodID, paymentMethod)
		if err != nil {
			return nil, err
		}
		out = append(out, *convertedPaymentMethod)
	}
	return out, nil
}

func markPaymentMethodAsDefaultForUser(stripeCustomerID, paymentMethodID string) error {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	if _, err := customer.Update(stripeCustomerID, &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: ptr.String(paymentMethodID),
		},
	}); err != nil {
		return err
	}
	return nil
}

func deletePaymentMethod(stripeCustomerID, paymentMethodID string) error {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	if _, err := paymentmethod.Detach(paymentMethodID, nil); err != nil {
		return err
	}
	return nil
}

func convertStripePaymentMethod(defaultPaymentMethodID *string, stripePaymentMethod *stripe.PaymentMethod) (*PaymentMethod, error) {
	if stripePaymentMethod.Card == nil {
		return nil, fmt.Errorf("Card is null, but expected only card type")
	}
	return &PaymentMethod{
		ExternalID:     stripePaymentMethod.ID,
		DisplayMask:    stripePaymentMethod.Card.Last4,
		CardExpiration: fmt.Sprintf("%2d/%d", stripePaymentMethod.Card.ExpMonth, stripePaymentMethod.Card.ExpYear),
		CardType:       getCardTypeForStripeBrand(stripePaymentMethod.Card.Brand),
		IsDefault:      defaultPaymentMethodID != nil && stripePaymentMethod.ID == *defaultPaymentMethodID,
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
