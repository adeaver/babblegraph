package billing

import (
	"babblegraph/model/users"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/setupintent"
)

// This file is for any stripe specific methods

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
