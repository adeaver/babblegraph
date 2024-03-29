package content

import (
	"babblegraph/model/admin"
	"babblegraph/model/content"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/deref"
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
	Title                 *string `json:"title,omitempty"`
	URL                   string  `json:"url"`
	Type                  string  `json:"type"`
	IngestStrategy        string  `json:"ingest_strategy"`
	LanguageCode          string  `json:"language_code"`
	ShouldUseURLAsSeedURL bool    `json:"should_use_url_as_seed_url"`
	MonthlyAccessLimit    *int64  `json:"monthly_access_limit"`
	Country               string  `json:"country"`
}

type addSourceResponse struct {
	Title string           `json:"title"`
	ID    content.SourceID `json:"id"`
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
	title := deref.String(req.Title, u.Domain)
	url, err := urlparser.EnsureProtocol(u.URL)
	if err != nil {
		return nil, err
	}
	var sourceID *content.SourceID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sourceID, err = content.InsertSource(tx, content.InsertSourceInput{
			Title:                 title,
			LanguageCode:          *languageCode,
			Type:                  *sourceType,
			IngestStrategy:        *ingestStrategy,
			Country:               *countryCode,
			MonthlyAccessLimit:    req.MonthlyAccessLimit,
			URL:                   *url,
			ShouldUseURLAsSeedURL: req.ShouldUseURLAsSeedURL,
		})
		return err
	}); err != nil {
		return nil, err
	}
	return addSourceResponse{
		Title: title,
		ID:    *sourceID,
	}, nil
}

type updateSourceRequest struct {
	ID                    content.SourceID `json:"id"`
	Title                 string           `json:"title"`
	LanguageCode          string           `json:"language_code"`
	URL                   string           `json:"url"`
	Type                  string           `json:"type"`
	IngestStrategy        string           `json:"ingest_strategy"`
	IsActive              bool             `json:"is_active"`
	ShouldUseURLAsSeedURL bool             `json:"should_use_url_as_seed_url"`
	MonthlyAccessLimit    *int64           `json:"monthly_access_limit"`
	Country               string           `json:"country"`
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
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return nil, err
	}
	u := urlparser.ParseURL(req.URL)
	if u == nil {
		return nil, fmt.Errorf("Invalid URL")
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return content.UpdateSource(tx, req.ID, content.UpdateSourceInput{
			Title:                 req.Title,
			LanguageCode:          *languageCode,
			Type:                  *sourceType,
			IngestStrategy:        *ingestStrategy,
			Country:               *countryCode,
			MonthlyAccessLimit:    req.MonthlyAccessLimit,
			URL:                   u.URL,
			IsActive:              req.IsActive,
			ShouldUseURLAsSeedURL: req.ShouldUseURLAsSeedURL,
		})
	}); err != nil {
		return nil, err
	}
	return updateSourceResponse{
		Success: true,
	}, nil
}

type getAllSourceSeedsForSourceRequest struct {
	SourceID content.SourceID `json:"source_id"`
}

type getAllSourceSeedsForSourceResponse struct {
	SourceSeeds []content.SourceSeed `json:"source_seeds"`
}

func getAllSourceSeedsForSource(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getAllSourceSeedsForSourceRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var sourceSeeds []content.SourceSeed
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sourceSeeds, err = content.GetAllSourceSeedsForSource(tx, req.SourceID)
		return err
	}); err != nil {
		return nil, err
	}
	return getAllSourceSeedsForSourceResponse{
		SourceSeeds: sourceSeeds,
	}, nil
}

type addSourceSeedRequest struct {
	SourceID content.SourceID `json:"source_id"`
	URL      string           `json:"url"`
}

type addSourceSeedResponse struct {
	ID content.SourceSeedID `json:"id"`
}

func addSourceSeed(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req addSourceSeedRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	u := urlparser.ParseURL(req.URL)
	if u == nil {
		return nil, fmt.Errorf("Invalid URL")
	}
	var sourceSeedID *content.SourceSeedID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sourceSeedID, err = content.AddSourceSeed(tx, req.SourceID, *u, false)
		return err
	}); err != nil {
		return nil, err
	}
	return addSourceSeedResponse{
		ID: *sourceSeedID,
	}, nil
}

type updateSourceSeedRequest struct {
	SourceSeedID content.SourceSeedID `json:"source_seed_id"`
	URL          string               `json:"url"`
	IsActive     bool                 `json:"is_active"`
}

type updateSourceSeedResponse struct {
	Success bool `json:"success"`
}

func updateSourceSeed(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req updateSourceSeedRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	u := urlparser.ParseURL(req.URL)
	if u == nil {
		return nil, fmt.Errorf("Invalid URL")
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return content.UpdateSourceSeed(tx, req.SourceSeedID, *u, req.IsActive)
	}); err != nil {
		return nil, err
	}
	return updateSourceSeedResponse{
		Success: true,
	}, nil
}

type getSourceSourceSeedMappingsForSourceRequest struct {
	SourceID content.SourceID `json:"source_id"`
}

type getSourceSourceSeedMappingsForSourceResponse struct {
	SourceSeedMappings []content.SourceSeedTopicMapping `json:"source_seed_mappings"`
}

func getSourceSourceSeedMappingsForSource(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getSourceSourceSeedMappingsForSourceRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var sourceSeedMappings []content.SourceSeedTopicMapping
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		sourceSeeds, err := content.GetAllSourceSeedsForSource(tx, req.SourceID)
		if err != nil {
			return err
		}
		var sourceSeedIDs []content.SourceSeedID
		for _, s := range sourceSeeds {
			sourceSeedIDs = append(sourceSeedIDs, s.ID)
		}
		sourceSeedMappings, err = content.GetAllSourceSeedTopicMappings(tx, sourceSeedIDs)
		return err
	}); err != nil {
		return nil, err
	}
	return getSourceSourceSeedMappingsForSourceResponse{
		SourceSeedMappings: sourceSeedMappings,
	}, nil
}

type sourceSeedMappingsUpdate struct {
	SourceSeedID content.SourceSeedID `json:"source_seed_id"`
	TopicIDs     []content.TopicID    `json:"topic_ids"`
	IsActive     bool                 `json:"is_active"`
}

type upsertSourceSeedMappingsRequest struct {
	Updates []sourceSeedMappingsUpdate `json:"updates"`
}

type upsertSourceSeedMappingsResponse struct {
	Success bool `json:"success"`
}

func upsertSourceSeedMappings(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req upsertSourceSeedMappingsRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		for _, u := range req.Updates {
			for _, t := range u.TopicIDs {
				if err := content.UpsertSourceSeedMapping(tx, u.SourceSeedID, t, u.IsActive); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return upsertSourceSeedMappingsResponse{
		Success: true,
	}, nil
}

type getSourceFilterForSourceIDRequest struct {
	SourceID content.SourceID `json:"source_id"`
}

type getSourceFilterForSourceIDResponse struct {
	SourceFilter *content.SourceFilter `json:"source_filter"`
}

func getSourceFilterForSourceID(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getSourceFilterForSourceIDRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var sourceFilter *content.SourceFilter
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sourceFilter, err = content.LookupSourceFilterForSource(tx, req.SourceID)
		return err
	}); err != nil {
		return nil, err
	}
	return getSourceFilterForSourceIDResponse{
		SourceFilter: sourceFilter,
	}, nil
}

type upsertSourceFilterForSourceRequest struct {
	SourceID            content.SourceID `json:"source_id"`
	IsActive            bool             `json:"is_active"`
	UseLDJSONValidation *bool            `json:"use_ld_json_validation,omitempty"`
	PaywallClasses      []string         `json:"paywall_classes"`
	PaywallIDs          []string         `json:"paywall_ids"`
}

type upsertSourceFilterForSourceResponse struct {
	SourceFilterID content.SourceFilterID `json:"source_filter_id"`
}

func upsertSourceFilterForSource(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req upsertSourceFilterForSourceRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var sourceFilterID *content.SourceFilterID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		sourceFilterID, err = content.UpsertSourceFilterForSource(tx, req.SourceID, content.UpsertSourceFilterForSourceInput{
			IsActive:            req.IsActive,
			UseLDJSONValidation: req.UseLDJSONValidation,
			PaywallClasses:      req.PaywallClasses,
			PaywallIDs:          req.PaywallIDs,
		})
		return err
	}); err != nil {
		return nil, err
	}
	return upsertSourceFilterForSourceResponse{
		SourceFilterID: *sourceFilterID,
	}, nil
}
