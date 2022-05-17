package routermiddleware

import (
	"babblegraph/model/billing"
	"babblegraph/services/web/router"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/encrypt"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	PromotionCodeCookieName = "bbgpromo"

	promotionCodeValidity = 30 * 24 * time.Hour
)

type WithMaybePromotionHandler func(promotionCode *billing.PromotionCode, r *router.Request) (interface{}, error)

func WithMaybePromotion(handler WithMaybePromotionHandler) func(r *router.Request) (interface{}, error) {
	return func(r *router.Request) (interface{}, error) {
		var promotionCode *billing.PromotionCode
		promotionCookieTokenValue := lookupPromotionCodeCookie(r)
		if promotionCookieTokenValue != nil {
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				promotionCode, err = billing.LookupPromotionByCode(tx, promotionCookieTokenValue.Code)
				return err
			}); err != nil {
				return nil, err
			}
		}
		return handler(promotionCode, r)
	}
}

func SetPromotionCodeIfActive(r *router.Request, promotionCode billing.PromotionCode) error {
	promoToken, err := encodePromotionCode(promotionCode)
	if err != nil {
		return err
	}
	if promotionCode.IsActive {
		if lookupPromotionCodeCookie(r) == nil {
			r.RespondWithCookie(&http.Cookie{
				Name:     PromotionCodeCookieName,
				Value:    *promoToken,
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Now().Add(30 * 24 * time.Hour),
			})
		}
	}
	return nil
}

func lookupPromotionCodeCookie(r *router.Request) *promotionCodeTokenValue {
	for _, cookie := range r.GetCookies() {
		if cookie.Name == PromotionCodeCookieName {
			tokenValue := decodePromotionCode(r, cookie.Value)
			if tokenValue != nil {
				if time.Unix(tokenValue.SetAtUnixTimestamp, 0).Add(promotionCodeValidity).Before(time.Now()) {
					return nil
				}
				return tokenValue
			}
		}
	}
	return nil
}

type promotionCodeTokenValue struct {
	SetAtUnixTimestamp int64  `json:"set_at_unix_timestamp"`
	Code               string `json:"code"`
}

func encodePromotionCode(promotionCode billing.PromotionCode) (*string, error) {
	return encrypt.GetToken(encrypt.TokenPair{
		Key: PromotionCodeCookieName,
		Value: promotionCodeTokenValue{
			SetAtUnixTimestamp: time.Now().Unix(),
			Code:               promotionCode.Code,
		},
	})
}

func decodePromotionCode(c ctx.LogContext, tokenStr string) *promotionCodeTokenValue {
	var out *promotionCodeTokenValue
	if err := encrypt.WithDecodedToken(tokenStr, func(tokenPair encrypt.TokenPair) error {
		if tokenPair.Key != PromotionCodeCookieName {
			return fmt.Errorf("Incorrect key for cookie value: %s", tokenPair.Key)
		}
		val, ok := tokenPair.Value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Value did not correctly parse as a map, was %T", tokenPair.Value)
		}
		setAtUnixTimestampInterface, ok := val["set_at_unix_timestamp"]
		if !ok {
			return fmt.Errorf("Map %+v did not have set_at_unix_timestamp", val)
		}
		setAtUnixTimestamp, ok := setAtUnixTimestampInterface.(float64)
		if !ok {
			return fmt.Errorf("Set at unix timestamp was a %T not a float64", setAtUnixTimestampInterface)
		}
		codeInterface, ok := val["code"]
		if !ok {
			return fmt.Errorf("Map %+v did not have code", val)
		}
		code, ok := codeInterface.(string)
		if !ok {
			return fmt.Errorf("Code was a %T not a string", codeInterface)
		}
		out = &promotionCodeTokenValue{
			SetAtUnixTimestamp: int64(setAtUnixTimestamp),
			Code:               code,
		}
		return nil
	}); err != nil {
		c.Debugf("Error decoding token %s", err.Error())
	}
	return out
}
