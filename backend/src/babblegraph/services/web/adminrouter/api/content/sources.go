package content

import (
	"babblegraph/model/admin"
	"babblegraph/model/content"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/geo"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type getAllSourcesRequest struct{}

type getAllSourcesResponse struct {
	Sources []content.Source `json:"sources"`
}

func getAllSources(adminID admin.ID, r *router.Request) (interface{}, error) {
	var sources []content.Source
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sources, err = content.GetAllSources(tx)
		return err
	}); err != nil {
		return nil, err
	}
	return getAllSourcesResponse{
		Sources: sources,
	}, nil
}

type getSourceByIDRequest struct {
	ID content.SourceID `json:"id"`
}

type getSourceByIDResponse struct {
	Source content.Source `json:"source"`
}

func getSourceByID(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getSourceByIDRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var source *content.Source
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		source, err = content.GetSource(tx, req.ID)
		return err
	}); err != nil {
		return nil, err
	}
	return getSourceByIDResponse{
		Source: *source,
	}, nil
}

type addSourceRequest struct {
	URL                string `json:"url"`
	Type               string `json:"type"`
	IngestStrategy     string `json:"ingest_strategy"`
	LanguageCode       string `json:"language_code"`
	MonthlyAccessLimit *int64 `json:"monthly_access_limit"`
	Country            string `json:"country"`
}

type addSourceResponse struct {
	ID content.SourceID `json:"id"`
}

func addSource(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req addSourceRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	countryCode, err := geo.GetCountryCodeFromString(req.Country)
	if err != nil {
		return nil, err
	}
	ingestStrategy, err := content.GetIngestStrategyFromString(req.IngestStrategy)
	if err != nil {
		return nil, err
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return nil, err
	}
	sourceType, err := content.GetSourceTypeFromString(req.Type)
	if err != nil {
		return nil, err
	}
	u := urlparser.ParseURL(req.URL)
	if u == nil {
		return nil, fmt.Errorf("Invalid URL")
	}
	var sourceID *content.SourceID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sourceID, err = content.InsertSource(tx, content.InsertSourceInput{
			LanguageCode:       *languageCode,
			Type:               *sourceType,
			IngestStrategy:     *ingestStrategy,
			Country:            *countryCode,
			MonthlyAccessLimit: req.MonthlyAccessLimit,
			URL:                u.URL,
		})
		return err
	}); err != nil {
		return nil, err
	}
	return addSourceResponse{
		ID: *sourceID,
	}, nil
}

type updateSourceRequest struct {
	ID                 content.SourceID `json:"id"`
	URL                string           `json:"url"`
	Type               string           `json:"type"`
	IngestStrategy     string           `json:"ingest_strategy"`
	IsActive           bool             `json:"is_active"`
	MonthlyAccessLimit *int64           `json:"monthly_access_limit"`
	Country            string           `json:"country"`
}

type updateSourceResponse struct {
	Success bool `json:"success"`
}

func updateSource(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req updateSourceRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	countryCode, err := geo.GetCountryCodeFromString(req.Country)
	if err != nil {
		return nil, err
	}
	ingestStrategy, err := content.GetIngestStrategyFromString(req.IngestStrategy)
	if err != nil {
		return nil, err
	}
	sourceType, err := content.GetSourceTypeFromString(req.Type)
	if err != nil {
		return nil, err
	}
	u := urlparser.ParseURL(req.URL)
	if u == nil {
		return nil, fmt.Errorf("Invalid URL")
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return content.UpdateSource(tx, req.ID, content.UpdateSourceInput{
			Type:               *sourceType,
			IngestStrategy:     *ingestStrategy,
			Country:            *countryCode,
			MonthlyAccessLimit: req.MonthlyAccessLimit,
			URL:                u.URL,
			IsActive:           req.IsActive,
		})
	}); err != nil {
		return nil, err
	}
	return updateSourceResponse{
		Success: true,
	}, nil
}
