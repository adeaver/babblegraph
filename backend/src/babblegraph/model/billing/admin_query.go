package billing

import (
	"babblegraph/model/users"
	"babblegraph/util/ctx"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/coupon"
	"github.com/stripe/stripe-go/v72/promotioncode"
)

const (
	createPromotionCodeQuery = "INSERT INTO billing_promotion_codes (type, external_id_mapping_id) VALUES ($1, $2)"
)

type UserBillingInformation struct {
	UserID         users.UserID                    `json:"user_id"`
	ExternalIDType string                          `json:"external_id_type"`
	Subscriptions  []PremiumNewsletterSubscription `json:"subscriptions"`
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
		return &out, nil
	default:
		return nil, fmt.Errorf("Unsupported external ID type %s", externalID.IDType)
	}
}

type CreatePromotionCodeInput struct {
	Code           string
	Discount       Discount
	MaxRedemptions *int64
	PromotionType  PromotionType
}

func CreatePromotionCode(c ctx.LogContext, tx *sqlx.Tx, input CreatePromotionCodeInput) (*PromotionCode, error) {
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	var stripeCoupon *stripe.Coupon
	var stripePromotionCode *stripe.PromotionCode
	var err error
	defer func() {
		if err != nil {
			if stripeCoupon != nil {
				if _, err := coupon.Del(stripeCoupon.ID, nil); err != nil {
					c.Errorf("Error rolling back coupon %s", stripeCoupon.ID)
				}
			}
			if stripePromotionCode != nil {
				if _, err := promotioncode.Update(stripePromotionCode.ID, &stripe.PromotionCodeParams{
					Active: ptr.Bool(false),
				}); err != nil {
					c.Errorf("Error rolling back promotion code %s", stripePromotionCode.ID)
				}
			}
		}
	}()
	couponParams := &stripe.CouponParams{
		Currency: ptr.String("USD"),
	}
	switch {
	case input.Discount.AmountOffCents != nil:
		couponParams.AmountOff = input.Discount.AmountOffCents
	case input.Discount.PercentOffBPS != nil:
		couponParams.PercentOff = stripe.Float64(float64(*input.Discount.PercentOffBPS) / 100.0)
	default:
		return nil, fmt.Errorf("Must specify either amount off or percent off")
	}
	couponParams.MaxRedemptions = input.MaxRedemptions
	stripeCoupon, err = coupon.New(couponParams)
	if err != nil {
		return nil, err
	}
	stripePromotionCode, err = promotioncode.New(&stripe.PromotionCodeParams{
		Coupon: stripe.String(stripeCoupon.ID),
		Code:   ptr.String(input.Code),
		Active: ptr.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	externalIDMappingID, err := insertExternalIDMapping(tx, stripePromotionCode.ID)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(createPromotionCodeQuery, input.PromotionType, externalIDMappingID); err != nil {
		return nil, err
	}
	return nil, nil
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
