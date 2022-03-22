package advertising

import (
	"babblegraph/model/admin"
	"babblegraph/model/advertising"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"fmt"

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

type insertVendorRequest struct {
	Name       string `json:"name"`
	WebsiteURL string `json:"website_url"`
}

type insertVendorResponse struct {
	ID advertising.VendorID `json:"id"`
}

func insertVendor(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req insertVendorRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	parsedURL := urlparser.ParseURL(req.WebsiteURL)
	if parsedURL == nil {
		return nil, fmt.Errorf("Invalid URL")
	}
	var id *advertising.VendorID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		id, err = advertising.InsertVendor(tx, req.Name, *parsedURL, false)
		return err
	}); err != nil {
		return nil, err
	}
	return insertVendorResponse{
		ID: *id,
	}, nil
}

type editVendorRequest struct {
	ID         advertising.VendorID `json:"id"`
	Name       string               `json:"name"`
	WebsiteURL string               `json:"website_url"`
	IsActive   bool                 `json:"is_active"`
}

type editVendorResponse struct {
	Success bool `json:"success"`
}

func editVendor(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req editVendorRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	parsedURL := urlparser.ParseURL(req.WebsiteURL)
	if parsedURL == nil {
		return nil, fmt.Errorf("Invalid URL")
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return advertising.EditVendor(tx, req.ID, req.Name, *parsedURL, req.IsActive)
	}); err != nil {
		return nil, err
	}
	return editVendorResponse{
		Success: true,
	}, nil
}
