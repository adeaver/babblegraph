package advertising

import (
	"babblegraph/model/content"
	"babblegraph/model/email"
	"babblegraph/wordsmith"
	"fmt"
	"strings"
	"time"
)

type VendorID string

type dbVendor struct {
	CreatedAt      time.Time `db:"created_at"`
	LastModifiedAt time.Time `db:"last_modified_at"`
	ID             VendorID  `db:"_id"`
	IsActive       bool      `db:"is_active"`
	Name           string    `db:"name"`
	WebsiteURL     string    `db:"website_url"`
}

func (d dbVendor) ToNonDB() Vendor {
	return Vendor{
		ID:         d.ID,
		IsActive:   d.IsActive,
		Name:       d.Name,
		WebsiteURL: d.WebsiteURL,
	}
}

type Vendor struct {
	ID         VendorID `json:"id"`
	IsActive   bool     `json:"is_active"`
	Name       string   `json:"name"`
	WebsiteURL string   `json:"website_url"`
}

type AdvertisementSourceID string

type dbAdvertisementSource struct {
	CreatedAt      time.Time               `db:"created_at"`
	LastModifiedAt time.Time               `db:"last_modified_at"`
	ID             AdvertisementSourceID   `db:"_id"`
	Name           string                  `db:"name"`
	URL            string                  `db:"url"`
	Type           AdvertisementSourceType `db:"type"`
	IsActive       bool                    `db:"is_active"`
}

func (d dbAdvertisementSource) ToNonDB() AdvertisementSource {
	return AdvertisementSource{
		ID:       d.ID,
		Name:     d.Name,
		URL:      d.URL,
		Type:     d.Type,
		IsActive: d.IsActive,
	}
}

type AdvertisementSource struct {
	ID       AdvertisementSourceID   `json:"id"`
	Name     string                  `json:"name"`
	URL      string                  `json:"url"`
	Type     AdvertisementSourceType `json:"type"`
	IsActive bool                    `json:"is_active"`
}

type AdvertisementSourceType string

const (
	AdvertisementSourceTypeAffiliate AdvertisementSourceType = "affiliate"
)

func (a AdvertisementSourceType) Ptr() *AdvertisementSourceType {
	return &a
}

func (a AdvertisementSourceType) Str() string {
	return string(a)
}

func GetSourceTypeFromString(s string) (*AdvertisementSourceType, error) {
	switch strings.ToLower(s) {
	case AdvertisementSourceTypeAffiliate.Str():
		return AdvertisementSourceTypeAffiliate.Ptr(), nil
	default:
		return nil, fmt.Errorf("Unsupported source type %s", s)
	}
}

type CampaignID string

type dbCampaign struct {
	CreatedAt      time.Time             `db:"created_at"`
	LastModifiedAt time.Time             `db:"last_modified_at"`
	ID             CampaignID            `db:"_id"`
	VendorID       VendorID              `db:"vendor_id"`
	Source         AdvertisementSourceID `db:"source_id"`
	URL            string                `db:"url"`
	IsActive       bool                  `db:"is_active"`
	Name           string                `db:"name"`
	ExpiresAt      *time.Time            `db:"expires_at"`
}

type CampaignTopicMappingID string

type dbCampaignTopicMapping struct {
	CreatedAt      time.Time              `db:"created_at"`
	LastModifiedAt time.Time              `db:"last_modified_at"`
	ID             CampaignTopicMappingID `db:"_id"`
	CampaignID     CampaignID             `db:"campaign_id"`
	TopicID        content.TopicID        `db:"topic_id"`
	IsActive       bool                   `db:"is_active"`
}

type AdvertisementID string

type dbAdvertisement struct {
	CreatedAt      time.Time              `db:"created_at"`
	LastModifiedAt time.Time              `db:"last_modified_at"`
	ID             AdvertisementID        `db:"_id"`
	LanguageCode   wordsmith.LanguageCode `db:"language_code"`
	CampaignID     CampaignID             `db:"campaign_id"`
	Title          string                 `db:"title"`
	ImageURL       string                 `db:"image_url"`
	Description    string                 `db:"description"`
	IsActive       bool                   `db:"is_active"`
}

type UserAdvertisementID string

type dbUserAdvertisement struct {
	CreatedAt       time.Time           `db:"created_at"`
	LastModifiedAt  time.Time           `db:"last_modified_at"`
	ID              UserAdvertisementID `db:"_id"`
	AdvertisementID AdvertisementID     `db:"advertisement_id"`
	CampaignID      CampaignID          `db:"campaign_id"`
	EmailRecordID   email.ID            `db:"email_record_id"`
}

type AdvertisementClickID string

type dbAdvertisementClick struct {
	CreatedAt           time.Time            `db:"created_at"`
	LastModifiedAt      time.Time            `db:"last_modified_at"`
	ID                  AdvertisementClickID `db:"_id"`
	UserAdvertisementID UserAdvertisementID  `db:"user_advertisement_id"`
}
