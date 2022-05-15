package index

import (
	"babblegraph/model/billing"
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
	"babblegraph/util/cache"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func getPromotionCodeRoutes() []router.IndexRoute {
	c := ctx.GetDefaultLogContext()
	var promotionCodes []billing.PromotionCode
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		promotionCodes, err = billing.GetAllURLPromotionCodes(c, tx)
		return err
	}); err != nil {
		c.Errorf("Error making promotion routes: %s", err.Error())
	}
	var out []router.IndexRoute
	for _, promoCode := range promotionCodes {
		promoCode := promoCode
		c.Debugf("Registering route for code %s", promoCode.Code)
		out = append(out, router.IndexRoute{
			Path: router.IndexPath{
				Text: fmt.Sprintf("/%s", strings.ToLower(promoCode.Code)),
			},
			Handler: routermiddleware.WithNoBodyRequestLogger(makePromotionCodeHandler(promoCode.Code)),
		})
	}
	return out
}

func makePromotionCodeHandler(code string) func(r *router.Request) (interface{}, error) {
	return func(r *router.Request) (interface{}, error) {
		r.Debugf("In route for %s", code)
		var promotionCode billing.PromotionCode
		if err := cache.WithCache(billing.GetPromotionCodeCacheKey(code), &promotionCode, 3*time.Hour, func() (interface{}, error) {
			var promo *billing.PromotionCode
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				var err error
				promo, err = billing.LookupPromotionByCode(tx, code)
				return err
			}); err != nil {
				return nil, err
			}
			return promo, nil
		}); err != nil {
			return nil, err
		}
		r.Debugf("Got promotion %+v", promotionCode)
		if err := routermiddleware.SetPromotionCodeIfActive(r, promotionCode); err != nil {
			return nil, err
		}
		return ptr.String(env.GetAbsoluteURLForEnvironment("")), nil
	}
}
