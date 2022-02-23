package billing

import (
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserBillingInformation struct {
	UserID            users.UserID                    `json:"user_id"`
	UserAccountStatus *useraccounts.SubscriptionLevel `json:"user_account_status"`
	ExternalIDType    string                          `json:"external_id_type"`
	Subscriptions     []PremiumNewsletterSubscription `json:"subscriptions"`
}

func GetBillingInformationForEmailAddress(c ctx.LogContext, tx *sqlx.Tx, emailAddress string) (*UserBillingInformation, error) {
	user, err := users.LookupUserByEmailAddress(tx, emailAddress)
	switch {
	case err != nil:
		return nil, err
	case user == nil:
		c.Infof("User doesn't exist")
		return nil, nil
	}
	billingInformation, err := lookupBillingInformationForUserID(tx, user.ID)
	switch {
	case err != nil:
		return nil, err
	case billingInformation == nil:
		c.Infof("User %s does not have billing information", user.ID)
		return nil, nil
	}
	out := UserBillingInformation{
		UserID: user.ID,
	}
	externalID, err := getExternalIDMapping(tx, billingInformation.ExternalIDMappingID)
	if err != nil {
		return nil, err
	}
	switch externalID.IDType {
	case externalIDTypeStripe:
		out.ExternalIDType = string(externalIDTypeStripe)
		out.Subscriptions, err = getAllStripeSubscriptionsForUser(c, tx, billingInformation.ID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unsupported external ID type %s", externalID.IDType)
	}
	out.UserAccountStatus, err = useraccounts.LookupSubscriptionLevelForUser(tx, user.ID)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func getAllStripeSubscriptionsForUser(c ctx.LogContext, tx *sqlx.Tx, billingInformationID BillingInformationID) ([]PremiumNewsletterSubscription, error) {
	var matches []dbPremiumNewsletterSubscription
	err := tx.Select(&matches, lookupPremiumNewsletterSubscriptionQuery, billingInformationID)
	switch {
	case err != nil:
		return nil, err
	case len(matches) == 0:
		return nil, nil
	}
	var out []PremiumNewsletterSubscription
	for _, m := range matches {
		premiumNewsletterSubscription, err := getStripeSubscriptionAndConvertSubscriptionForDBPremiumNewsletterSubscription(c, tx, m, true)
		if err != nil {
			return nil, err
		}
		out = append(out, *premiumNewsletterSubscription)
	}
	return out, nil
}
