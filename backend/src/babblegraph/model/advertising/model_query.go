package advertising

import (
	"babblegraph/model/content"
	"babblegraph/util/urlparser"
	"babblegraph/wordsmith"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	getAllVendorsQuery  = "SELECT * FROM advertising_vendors"
	insertVendorQuery   = "INSERT INTO advertising_vendors (name, website_url, is_active) VALUES ($1, $2, $3) RETURNING _id"
	editVendorByIDQuery = "UPDATE advertising_vendors SET name = $1, website_url = $2, is_active = $3, last_modified_at = timezone('utc', now()) WHERE _id = $4"

	getAllSourcesQuery  = "SELECT * FROM advertising_sources"
	insertSourceQuery   = "INSERT INTO advertising_sources (name, url, type, is_active) VALUES ($1, $2, $3, $4)"
	editSourceByIDQuery = "UPDATE advertising_sources SET name = $1, url = $2, type = $3, is_active = $4, last_modified_at = timezone('utc', now()) WHERE _id = $5"

	getAllCampaignsQuery = "SELECT * FROM advertising_campaigns"
	getCampaignByIDQuery = "SELECT * FROM advertising_campaigns WHERE _id = $1"
	insertCampaignQuery  = `
    INSERT INTO
        advertising_campaigns (
            vendor_id,
            source_id,
            url,
            is_active,
            should_apply_to_all_users,
            name,
            expires_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7
        )`
	editCampaignByIDQuery = `
    UPDATE
        advertising_campaigns
    SET
        vendor_id = $1,
        source_id = $2,
        url = $3,
        is_active = $4,
        should_apply_to_all_users=$5,
        name = $6,
        expires_at = $7,
        last_modified_at = timezone('utc', now())
    WHERE _id = $8`

	getCampaignTopicMappingQuery                = "SELECT * FROM advertising_campaign_topic_mappings WHERE campaign_id = $1 AND is_active = TRUE"
	markAllCampaignTopicMappingsAsInActiveQuery = "UPDATE advertising_campaign_topic_mappings SET is_active=FALSE, last_modified_at = timezone('utc', now()) WHERE campaign_id = $1"
	upsertCampaignTopicMappingQuery             = `INSERT INTO
        advertising_campaign_topic_mappings (
            campaign_id, topic_id, is_active
        ) VALUES (
            $1, $2, $3
        ) ON CONFLICT (campaign_id, topic_id) DO UPDATE
        SET last_modified_at = timezone('utc', now()), is_active = $3`

	getAllAdvertisementsForCampaignID = "SELECT * FROM advertising_advertisement WHERE campaign_id = $1"
	insertAdvertisementQuery          = `INSERT INTO advertising_advertisements (
        language_code, campaign_id, title, image_url, description, is_active
    ) VALUES (
        $1, $2, $3, $4, $5, $6
    ) RETURNING _id`
	editAdvertisementQuery = `UPDATE advertising_advertisements SET
        language_code=$1,
        title=$2,
        image_url=$3,
        description=$4,
        is_active=$5
    WHERE
        _id = $6`
)

func GetAllVendors(tx *sqlx.Tx) ([]Vendor, error) {
	var matches []dbVendor
	if err := tx.Select(&matches, getAllVendorsQuery); err != nil {
		return nil, err
	}
	var out []Vendor
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func InsertVendor(tx *sqlx.Tx, name string, u urlparser.ParsedURL, isActive bool) (*VendorID, error) {
	rows, err := tx.Queryx(insertVendorQuery, name, u.URL, isActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vendorID VendorID
	for rows.Next() {
		if err := rows.Scan(&vendorID); err != nil {
			return nil, err
		}
	}
	return &vendorID, nil
}

func EditVendor(tx *sqlx.Tx, id VendorID, name string, u urlparser.ParsedURL, isActive bool) error {
	if _, err := tx.Exec(editVendorByIDQuery, name, u.URL, isActive, id); err != nil {
		return err
	}
	return nil
}

func GetAllSources(tx *sqlx.Tx) ([]AdvertisementSource, error) {
	var matches []dbAdvertisementSource
	if err := tx.Select(&matches, getAllSourcesQuery); err != nil {
		return nil, err
	}
	var out []AdvertisementSource
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func InsertSource(tx *sqlx.Tx, name string, u urlparser.ParsedURL, sourceType AdvertisementSourceType, isActive bool) (*AdvertisementSourceID, error) {
	rows, err := tx.Queryx(insertSourceQuery, name, u.URL, sourceType, isActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sourceID AdvertisementSourceID
	for rows.Next() {
		if err := rows.Scan(&sourceID); err != nil {
			return nil, err
		}
	}
	return &sourceID, nil
}

type EditSourceInput struct {
	Name       string
	URL        urlparser.ParsedURL
	SourceType AdvertisementSourceType
	IsActive   bool
}

func EditSource(tx *sqlx.Tx, id AdvertisementSourceID, input EditSourceInput) error {
	if _, err := tx.Exec(editSourceByIDQuery, input.Name, input.URL.URL, input.SourceType, input.IsActive, id); err != nil {
		return err
	}
	return nil
}

func GetAllCampaigns(tx *sqlx.Tx) ([]Campaign, error) {
	var matches []dbCampaign
	if err := tx.Select(&matches, getAllCampaignsQuery); err != nil {
		return nil, err
	}
	var out []Campaign
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

func GetCampaignByID(tx *sqlx.Tx, id CampaignID) (*Campaign, error) {
	var matches []dbCampaign
	err := tx.Select(&matches, getCampaignByIDQuery, id)
	switch {
	case err != nil:
		return nil, err
	case len(matches) != 1:
		return nil, fmt.Errorf("Expected exactly one campaign for id %s, but got %d", id, len(matches))
	default:
		out := matches[0].ToNonDB()
		return &out, nil
	}
}

type InsertNewCampaignInput struct {
	VendorID              VendorID
	SourceID              AdvertisementSourceID
	URL                   urlparser.ParsedURL
	IsActive              bool
	Name                  string
	ShouldApplyToAllUsers bool
	ExpiresAt             *time.Time
}

func InsertNewCampaign(tx *sqlx.Tx, input InsertNewCampaignInput) (*CampaignID, error) {
	rows, err := tx.Queryx(insertCampaignQuery, input.VendorID, input.SourceID, input.URL.URL, input.IsActive, input.ShouldApplyToAllUsers, input.Name, input.ExpiresAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var campaignID CampaignID
	for rows.Next() {
		if err := rows.Scan(&campaignID); err != nil {
			return nil, err
		}
	}
	return &campaignID, nil
}

type UpdateCampaignInput struct {
	VendorID              VendorID
	SourceID              AdvertisementSourceID
	URL                   urlparser.ParsedURL
	IsActive              bool
	Name                  string
	ShouldApplyToAllUsers bool
	ExpiresAt             *time.Time
}

func UpdateCampaign(tx *sqlx.Tx, id CampaignID, input UpdateCampaignInput) error {
	if _, err := tx.Exec(editCampaignByIDQuery, input.VendorID, input.SourceID, input.URL.URL, input.IsActive, input.ShouldApplyToAllUsers, input.Name, input.ExpiresAt, id); err != nil {
		return err
	}
	return nil
}

func GetActiveTopicMappingsForCampaignID(tx *sqlx.Tx, campaignID CampaignID) ([]content.TopicID, error) {
	var matches []dbCampaignTopicMapping
	if err := tx.Select(&matches, getCampaignTopicMappingQuery, campaignID); err != nil {
		return nil, err
	}
	var out []content.TopicID
	for _, m := range matches {
		out = append(out, m.TopicID)
	}
	return out, nil
}

func UpsertActiveTopicMappingsForCampaignID(tx *sqlx.Tx, campaignID CampaignID, topics []content.TopicID) error {
	if _, err := tx.Exec(markAllCampaignTopicMappingsAsInActiveQuery, campaignID); err != nil {
		return err
	}
	for _, t := range topics {
		if _, err := tx.Exec(upsertCampaignTopicMappingQuery, campaignID, t, true); err != nil {
			return err
		}
	}
	return nil
}

func GetAllAdvertisementsForCampaignID(tx *sqlx.Tx, campaignID CampaignID) ([]Advertisement, error) {
	var matches []dbAdvertisement
	if err := tx.Select(&matches, getAllAdvertisementsForCampaignID, campaignID); err != nil {
		return nil, err
	}
	var out []Advertisement
	for _, m := range matches {
		out = append(out, m.ToNonDB())
	}
	return out, nil
}

type InsertNewAdvertisementInput struct {
	LanguageCode wordsmith.LanguageCode
	CampaignID   CampaignID
	Title        string
	Description  string
	ImageURL     string
	IsActive     bool
}

func InsertNewAdvertisement(tx *sqlx.Tx, input InsertNewAdvertisementInput) (*AdvertisementID, error) {
	rows, err := tx.Queryx(insertAdvertisementQuery, input.LanguageCode, input.CampaignID, input.Title, input.ImageURL, input.Description, input.IsActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var advertisementID AdvertisementID
	for rows.Next() {
		if err := rows.Scan(&advertisementID); err != nil {
			return nil, err
		}
	}
	return &advertisementID, nil
}

type UpdateAdvertisementInput struct {
	LanguageCode wordsmith.LanguageCode
	Title        string
	Description  string
	ImageURL     string
	IsActive     bool
}

func UpdateAdvertisement(tx *sqlx.Tx, id AdvertisementID, input UpdateAdvertisementInput) error {
	if _, err := tx.Exec(editAdvertisementQuery, input.LanguageCode, input.Title, input.ImageURL, input.Description, input.IsActive, id); err != nil {
		return err
	}
	return nil
}
