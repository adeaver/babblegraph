package advertising

import (
	"babblegraph/model/admin"
	"babblegraph/model/advertising"
	"babblegraph/services/web/router"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

type getAllVendorsRequest struct{}

type getAllVendorsResponse struct {
	Vendors []advertising.Vendor `json:"vendors"`
}

func getAllVendors(adminID admin.ID, r *router.Request) (interface{}, error) {
	var vendors []advertising.Vendor
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		vendors, err = advertising.GetAllVendors(tx)
		return err
	}); err != nil {
		return nil, err
	}
	return getAllVendorsResponse{
		Vendors: vendors,
	}, nil
}
