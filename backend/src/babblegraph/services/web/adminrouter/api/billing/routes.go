package billing

import (
	"babblegraph/model/admin"
	"babblegraph/model/billing"
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
		},
	},
}

type getBillingInformationForEmailAddressRequest struct {
	EmailAddress string `json:"email_address"`
}

type getBillingInformationForEmailAddressResponse struct {
	BillingInformation *billing.UserBillingInformation `json:"billing_information"`
}

func getBillingInformationForEmailAddress(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getBillingInformationForEmailAddressRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	formattedEmailAddress := email.FormatEmailAddress(req.EmailAddress)
	var userBillingInformation *billing.UserBillingInformation
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		userBillingInformation, err = billing.GetBillingInformationForEmailAddress(r, tx, formattedEmailAddress)
		return err
	}); err != nil {
		return nil, err
	}
	return getBillingInformationForEmailAddressResponse{
		BillingInformation: userBillingInformation,
	}, nil
}
