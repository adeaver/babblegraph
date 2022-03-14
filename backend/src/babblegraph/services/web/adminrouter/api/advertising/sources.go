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

type getAllSourcesRequest struct{}

type getAllSourcesResponse struct {
	Sources []advertising.AdvertisementSource `json:"sources"`
}

func getAllSources(adminID admin.ID, r *router.Request) (interface{}, error) {
	var sources []advertising.AdvertisementSource
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sources, err = advertising.GetAllSources(tx)
		return err
	}); err != nil {
		return nil, err
	}
	return getAllSourcesResponse{
		Sources: sources,
	}, nil
}

type insertSourceRequest struct {
	SourceType string `json:"source_type"`
	Name       string `json:"name"`
	WebsiteURL string `json:"website_url"`
}

type insertSourceResponse struct {
	ID advertising.AdvertisementSourceID `json:"id"`
}

func insertSource(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req insertSourceRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	sourceType, err := advertising.GetSourceTypeFromString(req.SourceType)
	if err != nil {
		return nil, err
	}
	parsedURL := urlparser.ParseURL(req.WebsiteURL)
	if parsedURL == nil {
		return nil, err
	}
	var id *advertising.AdvertisementSourceID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		id, err = advertising.InsertSource(tx, req.Name, *parsedURL, *sourceType, false)
		return err
	}); err != nil {
		return nil, err
	}
	return insertSourceResponse{
		ID: *id,
	}, nil
}

type editSourceRequest struct {
	ID         advertising.AdvertisementSourceID `json:"id"`
	SourceType string                            `json:"source_type"`
	Name       string                            `json:"name"`
	WebsiteURL string                            `json:"website_url"`
	IsActive   bool                              `json:"is_active"`
}

type editSourceResponse struct {
	Success bool `json:"success"`
}

func editSource(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req editSourceRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	sourceType, err := advertising.GetSourceTypeFromString(req.SourceType)
	if err != nil {
		return nil, err
	}
	parsedURL := urlparser.ParseURL(req.WebsiteURL)
	if parsedURL == nil {
		return nil, fmt.Errorf("URL did not parse correctly")
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return advertising.EditSource(tx, req.ID, advertising.EditSourceInput{
			Name:       req.Name,
			URL:        *parsedURL,
			SourceType: *sourceType,
			IsActive:   req.IsActive,
		})
	}); err != nil {
		return nil, err
	}
	return editSourceResponse{
		Success: true,
	}, nil
}
