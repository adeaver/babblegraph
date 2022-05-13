package billing

import (
	"babblegraph/model/admin"
	"babblegraph/model/billing"
	"babblegraph/model/useraccounts"
	"babblegraph/model/users"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/email"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "billing",
	Routes: []router.Route{
		{
			Path: "get_billing_information_for_email_address_1",
			Handler: middleware.WithPermission(
				admin.PermissionManageBilling,
				getBillingInformationForEmailAddress,
			),
		}, {
			Path: "force_sync_for_user_1",
			Handler: middleware.WithPermission(
				admin.PermissionManageBilling,
				forceSyncForUser,
			),
		}, {
			Path: "create_promotion_code_1",
			Handler: middleware.WithPermission(
				admin.PermissionManageBilling,
				createPromotionCode,
			),
		},
	},
}

type getBillingInformationForEmailAddressRequest struct {
	EmailAddress string `json:"email_address"`
}

type getBillingInformationForEmailAddressResponse struct {
	BillingInformation *billing.UserBillingInformation `json:"billing_information"`
	UserAccountStatus  *useraccounts.SubscriptionLevel `json:"user_account_status"`
}

func getBillingInformationForEmailAddress(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getBillingInformationForEmailAddressRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	var userBillingInformation *billing.UserBillingInformation
	var userAccountStatus *useraccounts.SubscriptionLevel
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userBillingInformation, err = billing.GetBillingInformationForEmailAddress(r, tx, formattedEmailAddress)
		switch {
		case err != nil:
			return err
		case userBillingInformation == nil:
			return nil
		default:
			userAccountStatus, err = useraccounts.LookupSubscriptionLevelForUser(tx, userBillingInformation.UserID)
			return err
		}
	}); err != nil {
		return nil, err
	}
	return getBillingInformationForEmailAddressResponse{
		BillingInformation: userBillingInformation,
		UserAccountStatus:  userAccountStatus,
	}, nil
}

type forceSyncForUserRequest struct {
	UserID users.UserID `json:"user_id"`
}

type forceSyncForUserResponse struct {
	Success bool `json:"success"`
}

func forceSyncForUser(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req forceSyncForUserRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		premiumNewsletterSubscription, err := billing.LookupPremiumNewsletterSubscriptionForUser(r, tx, req.UserID)
		if err != nil {
			return err
		}
		return billing.SyncUserAccountWithPremiumNewsletterSubscription(tx, req.UserID, premiumNewsletterSubscription)
	}); err != nil {
		return nil, err
	}
	return forceSyncForUserResponse{
		Success: true,
	}, nil
}

type createPromotionCodeRequest struct {
	Code           string           `json:"code"`
	Discount       billing.Discount `json:"discount"`
	MaxRedemptions *int64           `json:"max_redemptions,omitempty"`
	PromotionType  string           `json:"promotion_type"`
}

type createPromotionCodeResponse struct {
	PromotionCode *billing.PromotionCode `json:"promotion_code"`
}

func createPromotionCode(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req createPromotionCodeRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	promotionType, err := billing.GetPromotionTypeForString(req.PromotionType)
	if err != nil {
		return nil, err
	}
	var promotionCode *billing.PromotionCode
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		promotionCode, err = billing.CreatePromotionCode(r, tx, billing.CreatePromotionCodeInput{
			Code:           req.Code,
			Discount:       req.Discount,
			MaxRedemptions: req.MaxRedemptions,
			PromotionType:  *promotionType,
		})
		return err
	}); err != nil {
		return nil, err
	}
	return createPromotionCodeResponse{
		PromotionCode: promotionCode,
	}, nil
}
