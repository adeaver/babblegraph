package content

import (
	"babblegraph/util/geo"
	"babblegraph/wordsmith"
	"fmt"
	"strings"
	"time"
)

type TopicID string

func (t TopicID) Ptr() *TopicID {
	return &t
}

type dbTopic struct {
	ID             TopicID   `db:"_id"`
	CreatedAt      time.Time `db:"created_at"`
	LastModifiedAt time.Time `db:"last_modified_at"`
	Label          string    `db:"label"`
	IsActive       bool      `db:"is_active"`
}

func (d dbTopic) ToNonDB() Topic {
	return Topic{
		ID:       d.ID,
		Label:    d.Label,
		IsActive: d.IsActive,
	}
}

type Topic struct {
	ID       TopicID `json:"id"`
	Label    string  `json:"label"`
	IsActive bool    `json:"is_active"`
}

type TopicDisplayNameID string

type dbTopicDisplayName struct {
	ID             TopicDisplayNameID     `db:"_id"`
	CreatedAt      time.Time              `db:"created_at"`
	LastModifiedAt time.Time              `db:"last_modified_at"`
	TopicID        TopicID                `db:"topic_id"`
	LanguageCode   wordsmith.LanguageCode `db:"language_code"`
	Label          string                 `db:"label"`
	IsActive       bool                   `db:"is_active"`
}

func (d dbTopicDisplayName) ToNonDB() TopicDisplayName {
	return TopicDisplayName{
		ID:           d.ID,
		TopicID:      d.TopicID,
		LanguageCode: d.LanguageCode,
		Label:        d.Label,
		IsActive:     d.IsActive,
	}
}

type TopicDisplayName struct {
	ID           TopicDisplayNameID     `json:"id"`
	TopicID      TopicID                `json:"topic_id"`
	LanguageCode wordsmith.LanguageCode `json:"language_code"`
	Label        string                 `json:"label"`
	IsActive     bool                   `json:"is_active"`
}

type SourceID string

func (s SourceID) Ptr() *SourceID {
	return &s
}

type dbSource struct {
	ID                    SourceID               `db:"_id"`
	CreatedAt             time.Time              `db:"created_at"`
	LastModifiedAt        time.Time              `db:"last_modified_at"`
	Title                 string                 `db:"title"`
	URL                   string                 `db:"url"`
	URLIdentifier         string                 `db:"url_identifier"`
	Type                  SourceType             `db:"type"`
	Country               geo.CountryCode        `db:"country"`
	IngestStrategy        IngestStrategy         `db:"ingest_strategy"`
	LanguageCode          wordsmith.LanguageCode `db:"language_code"`
	IsActive              bool                   `db:"is_active"`
	ShouldUseURLAsSeedURL bool                   `db:"should_use_url_as_seed_url"`
	MonthlyAccessLimit    *int64                 `db:"monthly_access_limit"`
}

func (d dbSource) ToNonDB() Source {
	return Source{
		ID:                    d.ID,
		URL:                   d.URL,
		Title:                 d.Title,
		Type:                  d.Type,
		Country:               d.Country,
		IngestStrategy:        d.IngestStrategy,
		LanguageCode:          d.LanguageCode,
		IsActive:              d.IsActive,
		ShouldUseURLAsSeedURL: d.ShouldUseURLAsSeedURL,
		MonthlyAccessLimit:    d.MonthlyAccessLimit,
	}
}

type Source struct {
	ID                    SourceID               `json:"id"`
	Title                 string                 `json:"title"`
	URL                   string                 `json:"url"`
	Type                  SourceType             `json:"type"`
	Country               geo.CountryCode        `json:"country"`
	IngestStrategy        IngestStrategy         `json:"ingest_strategy"`
	LanguageCode          wordsmith.LanguageCode `json:"language_code"`
	IsActive              bool                   `json:"is_active"`
	ShouldUseURLAsSeedURL bool                   `json:"should_use_url_as_seed_url"`
	MonthlyAccessLimit    *int64                 `json:"monthly_access_limit"`
}

type SourceType string

const (
	SourceTypeNewsWebsite SourceType = "news-website"
)

func (s SourceType) Str() string {
	return string(s)
}

func (s SourceType) Ptr() *SourceType {
	return &s
}

func GetSourceTypeFromString(t string) (*SourceType, error) {
	switch strings.ToLower(t) {
	case SourceTypeNewsWebsite.Str():
		return SourceTypeNewsWebsite.Ptr(), nil
	default:
		return nil, fmt.Errorf("Unsupported source type: %s", t)
	}
}

type IngestStrategy string

const (
	IngestStrategyWebsiteHTML1 IngestStrategy = "website-html-1"
)

func (i IngestStrategy) Str() string {
	return string(i)
}

func (i IngestStrategy) Ptr() *IngestStrategy {
	return &i
}

func GetIngestStrategyFromString(i string) (*IngestStrategy, error) {
	switch strings.ToLower(i) {
	case IngestStrategyWebsiteHTML1.Str():
		return IngestStrategyWebsiteHTML1.Ptr(), nil
	default:
		return nil, fmt.Errorf("Unsupported ingest strategy type: %s", i)
	}
}

type SourceTopicMappingID string

type dbSourceTopicMapping struct {
	ID             SourceTopicMappingID `db:"_id"`
	CreatedAt      time.Time            `db:"created_at"`
	LastModifiedAt time.Time            `db:"last_modified_at"`
	RootID         SourceID             `db:"root_id"`
	TopicID        TopicID              `db:"topic_id"`
	IsActive       bool                 `db:"is_active"`
}

func (d dbSourceTopicMapping) ToNonDB() SourceTopicMapping {
	return SourceTopicMapping{
		ID:       d.ID,
		RootID:   d.RootID,
		TopicID:  d.TopicID,
		IsActive: d.IsActive,
	}
}

type SourceTopicMapping struct {
	ID       SourceTopicMappingID `json:"id"`
	RootID   SourceID             `json:"root_id"`
	TopicID  TopicID              `json:"topic_id"`
	IsActive bool                 `json:"is_active"`
}

type SourceSeedID string

func (s SourceSeedID) Ptr() *SourceSeedID {
	return &s
}

type dbSourceSeed struct {
	ID             SourceSeedID `db:"_id"`
	CreatedAt      time.Time    `db:"created_at"`
	LastModifiedAt time.Time    `db:"last_modified_at"`
	RootID         SourceID     `db:"root_id"`
	URL            string       `db:"url"`
	URLIdentifier  string       `db:"url_identifier"`
	URLParams      *string      `db:"url_params"`
	IsActive       bool         `db:"is_active"`
}

func (d dbSourceSeed) ToNonDB() SourceSeed {
	return SourceSeed{
		ID:       d.ID,
		RootID:   d.RootID,
		URL:      d.URL,
		IsActive: d.IsActive,
	}
}

type SourceSeed struct {
	ID       SourceSeedID `json:"id"`
	RootID   SourceID     `json:"root_id"`
	URL      string       `json:"url"`
	IsActive bool         `json:"is_active"`
}

type SourceFilterID string

const paywallFilterDelimiter = "#"

type dbSourceFilter struct {
	ID                  SourceFilterID `db:"_id"`
	CreatedAt           time.Time      `db:"created_at"`
	LastModifiedAt      time.Time      `db:"last_modified_at"`
	RootID              SourceID       `db:"root_id"`
	IsActive            bool           `db:"is_active"`
	UseLDJSONValidation *bool          `db:"use_ld_json_validation"`
	PaywallClasses      *string        `db:"paywall_classes"`
	PaywallIDs          *string        `db:"paywall_ids"`
}

func (d dbSourceFilter) ToNonDB() SourceFilter {
	var paywallClasses, paywallIDs []string
	if d.PaywallClasses != nil {
		paywallClasses = strings.Split(*d.PaywallClasses, paywallFilterDelimiter)
	}
	if d.PaywallIDs != nil {
		paywallIDs = strings.Split(*d.PaywallIDs, paywallFilterDelimiter)
	}
	return SourceFilter{
		ID:                  d.ID,
		RootID:              d.RootID,
		IsActive:            d.IsActive,
		UseLDJSONValidation: d.UseLDJSONValidation,
		PaywallClasses:      paywallClasses,
		PaywallIDs:          paywallIDs,
	}
}

type SourceFilter struct {
	ID                  SourceFilterID `json:"id"`
	RootID              SourceID       `json:"root_id"`
	IsActive            bool           `json:"is_active"`
	UseLDJSONValidation *bool          `json:"use_ld_json_validation"`
	PaywallClasses      []string       `json:"paywall_classes"`
	PaywallIDs          []string       `json:"paywall_ids"`
}

type SourceSeedTopicMappingID string

func (s SourceSeedTopicMappingID) Ptr() *SourceSeedTopicMappingID {
	return &s
}

type dbSourceSeedTopicMapping struct {
	ID             SourceSeedTopicMappingID `db:"_id"`
	CreatedAt      time.Time                `db:"created_at"`
	LastModifiedAt time.Time                `db:"last_modified_at"`
	SourceSeedID   SourceSeedID             `db:"source_seed_id"`
	TopicID        TopicID                  `db:"topic_id"`
	IsActive       bool                     `db:"is_active"`
}

func (d dbSourceSeedTopicMapping) ToNonDB() SourceSeedTopicMapping {
	return SourceSeedTopicMapping{
		ID:           d.ID,
		SourceSeedID: d.SourceSeedID,
		TopicID:      d.TopicID,
		IsActive:     d.IsActive,
	}
}

type SourceSeedTopicMapping struct {
	ID           SourceSeedTopicMappingID `json:"id"`
	SourceSeedID SourceSeedID             `json:"source_seed_id"`
	TopicID      TopicID                  `json:"topic_id"`
	IsActive     bool                     `json:"is_active"`
}

// This ID represents the union between either a topic-source mapping or a topic-source seed mapping
// since sources and source seeds can both contribute to topic mapping
type TopicMappingID string

func (t TopicMappingID) Ptr() *TopicMappingID {
	return &t
}

type MakeTopicMappingIDInput struct {
	SourceSeedTopicMappingID *SourceSeedTopicMappingID
	SourceTopicMappingID     *SourceTopicMappingID
}

func MustMakeTopicMappingID(input MakeTopicMappingIDInput) TopicMappingID {
	switch {
	case input.SourceSeedTopicMappingID != nil:
		return TopicMappingID(fmt.Sprintf("sourceseed-%s", *input.SourceSeedTopicMappingID))
	case input.SourceTopicMappingID != nil:
		return TopicMappingID(fmt.Sprintf("source-%s", *input.SourceTopicMappingID))
	default:
		panic("Neither source topic mapping or source seed topic mapping id are non-null")
	}
}
