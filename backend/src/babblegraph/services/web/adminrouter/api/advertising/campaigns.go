package advertising

import (
	"babblegraph/model/admin"
	"babblegraph/model/advertising"
	"babblegraph/model/content"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/urlparser"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type getCampaignRequest struct {
	ID advertising.CampaignID `json:"id"`
}

type getCampaignResponse struct {
	Campaign advertising.Campaign `json:"campaign"`
}

func getCampaign(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getCampaignRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var campaign *advertising.Campaign
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		campaign, err = advertising.GetCampaignByID(tx, req.ID)
		return err
	}); err != nil {
		return nil, err
	}
	return getCampaignResponse{
		Campaign: *campaign,
	}, nil
}

type getAllCampaignsRequest struct{}

type getAllCampaignsResponse struct {
	Campaigns []advertising.Campaign `json:"campaigns"`
}

func getAllCampaigns(adminID admin.ID, r *router.Request) (interface{}, error) {
	var campaigns []advertising.Campaign
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		campaigns, err = advertising.GetAllCampaigns(tx)
		return err
	}); err != nil {
		return nil, err
	}
	return getAllCampaignsResponse{
		Campaigns: campaigns,
	}, nil
}

type insertCampaignRequest struct {
	VendorID              advertising.VendorID              `json:"vendor_id"`
	SourceID              advertising.AdvertisementSourceID `json:"source_id"`
	URL                   string                            `json:"url"`
	Name                  string                            `json:"name"`
	ShouldApplyToAllUsers bool                              `json:"should_apply_to_all_users"`
}

type insertCampaignResponse struct {
	ID advertising.CampaignID `json:"id"`
}

func insertCampaign(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req insertCampaignRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	parsedURL := urlparser.ParseURL(req.URL)
	if parsedURL == nil {
		return nil, fmt.Errorf("URL %s did not parse", req.URL)
	}
	var campaignID *advertising.CampaignID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		campaignID, err = advertising.InsertNewCampaign(tx, advertising.InsertNewCampaignInput{
			VendorID:              req.VendorID,
			SourceID:              req.SourceID,
			URL:                   *parsedURL,
			Name:                  req.Name,
			ShouldApplyToAllUsers: req.ShouldApplyToAllUsers,
		})
		return err
	}); err != nil {
		return nil, err
	}
	return insertCampaignResponse{
		ID: *campaignID,
	}, nil
}

type updateCampaignRequest struct {
	CampaignID            advertising.CampaignID            `json:"campaign_id"`
	VendorID              advertising.VendorID              `json:"vendor_id"`
	SourceID              advertising.AdvertisementSourceID `json:"source_id"`
	URL                   string                            `json:"url"`
	Name                  string                            `json:"name"`
	ShouldApplyToAllUsers bool                              `json:"should_apply_to_all_users"`
	IsActive              bool                              `json:"is_active"`
	RolloutPercentage     int64                             `json:"rollout_percentage"`
}

type updateCampaignResponse struct {
	Success bool `json:"success"`
}

func updateCampaign(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req updateCampaignRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	parsedURL := urlparser.ParseURL(req.URL)
	if parsedURL == nil {
		return nil, fmt.Errorf("URL %s did not parse", req.URL)
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return advertising.UpdateCampaign(tx, req.CampaignID, advertising.UpdateCampaignInput{
			VendorID:              req.VendorID,
			SourceID:              req.SourceID,
			URL:                   *parsedURL,
			Name:                  req.Name,
			ShouldApplyToAllUsers: req.ShouldApplyToAllUsers,
			IsActive:              req.IsActive,
			RolloutPercentage:     req.RolloutPercentage,
		})
	}); err != nil {
		return nil, err
	}
	return updateCampaignResponse{
		Success: true,
	}, nil
}

type getCampaignTopicMappingsRequest struct {
	CampaignID advertising.CampaignID `json:"campaign_id"`
}

type getCampaignTopicMappingsResponse struct {
	TopicIDs []content.TopicID `json:"topic_ids"`
}

func getCampaignTopicMappings(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getCampaignTopicMappingsRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var topicIDs []content.TopicID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		topicIDs, err = advertising.GetActiveTopicMappingsForCampaignID(tx, req.CampaignID)
		return err
	}); err != nil {
		return nil, err
	}
	return getCampaignTopicMappingsResponse{
		TopicIDs: topicIDs,
	}, nil
}

type updateCampaignTopicMappingsRequest struct {
	CampaignID          advertising.CampaignID `json:"campaign_id"`
	ActiveTopicMappings []content.TopicID      `json:"active_topic_mappings"`
}

type updateCampaignTopicMappingsResponse struct {
	Success bool `json:"success"`
}

func updateCampaignTopicMappings(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req updateCampaignTopicMappingsRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return advertising.UpsertActiveTopicMappingsForCampaignID(tx, req.CampaignID, req.ActiveTopicMappings)
	}); err != nil {
		return nil, err
	}
	return updateCampaignTopicMappingsResponse{
		Success: true,
	}, nil
}
