package billing

import (
	"babblegraph/util/ctx"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/promotioncode"
)

const (
	lookupAllPromotionCodes    = "SELECT * FROM billing_promotion_codes"
	lookupPromotionByCodeQuery = "SELECT * FROM billing_promotion_codes WHERE code = $1"
)

func GetPromotionCodeCacheKey(code string) string {
	return fmt.Sprintf("billing_promotion_%s", code)
}

func GetAllURLPromotionCodes(c ctx.LogContext, tx *sqlx.Tx) ([]PromotionCode, error) {
	var matches []dbPromotionCode
	if err := tx.Select(&matches, lookupAllPromotionCodes); err != nil {
		return nil, err
	}
	var externalIDs []externalIDMappingID
	dbPromotionCodeByExternalIDs := make(map[externalIDMappingID]dbPromotionCode)
	for _, m := range matches {
		if m.Type != PromotionTypeURL {
			continue
		}
		externalIDs = append(externalIDs, m.ExternalIDMappingID)
		dbPromotionCodeByExternalIDs[m.ExternalIDMappingID] = m
	}
	if len(externalIDs) == 0 {
		c.Infof("No promotion codes found in database")
		return nil, nil
	}
	externalIDMappings, err := getManyExternalIDMapping(tx, externalIDs)
	if err != nil {
		return nil, err
	}
	dbPromotionCodesByID := make(map[string]dbPromotionCode)
	for _, externalIDMapping := range externalIDMappings {
		promotionCode, ok := dbPromotionCodeByExternalIDs[externalIDMapping.ID]
		if !ok {
			continue
		}
		externalIDWithType := fmt.Sprintf("%s_%s", externalIDMapping.IDType, externalIDMapping.ExternalID)
		dbPromotionCodesByID[externalIDWithType] = promotionCode
	}
	// Get Stripe Codes
	stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
	promotionCodeListParams := &stripe.PromotionCodeListParams{}
	promotionCodeListParams.AddExpand("data.coupon")
	stripePromotionCodesIterator := promotioncode.List(promotionCodeListParams)
	if stripePromotionCodesIterator.Err() != nil {
		return nil, stripePromotionCodesIterator.Err()
	}
	var out []PromotionCode
	for stripePromotionCodesIterator.Next() {
		stripePromotionCode, ok := stripePromotionCodesIterator.Current().(*stripe.PromotionCode)
		if !ok {
			c.Debugf("Could not convert to promotion code")
			continue
		}
		externalIDWithType := fmt.Sprintf("%s_%s", externalIDTypeStripe, stripePromotionCode.ID)
		dbPromotionCode, ok := dbPromotionCodesByID[externalIDWithType]
		if !ok {
			continue
		}
		promotionCode, err := mergePromotionCode(dbPromotionCode, stripePromotionCode)
		if err != nil {
			c.Warnf("Encountered an error while merging promotion code %s: %s", externalIDWithType, err.Error())
			continue
		}
		out = append(out, *promotionCode)
	}
	return out, nil
}

func LookupPromotionByCode(tx *sqlx.Tx, code string) (*PromotionCode, error) {
	var matches []dbPromotionCode
	if err := tx.Select(&matches, lookupPromotionByCodeQuery, code); err != nil {
		return nil, err
	}
	switch {
	case len(matches) == 0:
		return nil, nil
	case len(matches) > 1:
		return nil, fmt.Errorf("Expected at most one match for code %s but got %d", code, len(matches))
	default:
		externalIDMapping, err := getExternalIDMapping(tx, matches[0].ExternalIDMappingID)
		switch {
		case err != nil:
			return nil, err
		case externalIDMapping.IDType == externalIDTypeStripe:
			stripe.Key = env.MustEnvironmentVariable("STRIPE_KEY")
			searchParams := &stripe.PromotionCodeParams{}
			stripePromotionCode, err := promotioncode.Get(externalIDMapping.ExternalID, searchParams)
			if err != nil {
				return nil, err
			}
			return mergePromotionCode(matches[0], stripePromotionCode)
		default:
			return nil, fmt.Errorf("Unrecognized external ID type %s", externalIDMapping.IDType)
		}
	}
}

func mergePromotionCode(db dbPromotionCode, stripePromotionCode *stripe.PromotionCode) (*PromotionCode, error) {
	if stripePromotionCode.Coupon == nil {
		return nil, fmt.Errorf("Promotion code has no coupon")
	}
	isActive := stripePromotionCode.Coupon.Valid && stripePromotionCode.Active
	var percentOffBPS *int64
	if stripePromotionCode.Coupon.PercentOff > 1.0 {
		percentOffBPS = ptr.Int64(int64(stripePromotionCode.Coupon.PercentOff * 100))
	}
	return &PromotionCode{
		IsActive: isActive,
		Code:     db.Code,
		Type:     db.Type,
		Discount: Discount{
			AmountOffCents: ptr.Int64(stripePromotionCode.Coupon.AmountOff),
			PercentOffBPS:  percentOffBPS,
		},
	}, nil
}
