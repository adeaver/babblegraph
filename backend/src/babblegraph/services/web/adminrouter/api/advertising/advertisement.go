package advertising

import (
	"babblegraph/model/admin"
	"babblegraph/model/advertising"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/wordsmith"

	"github.com/jmoiron/sqlx"
)

type getAllAdvertisementsForCampaignRequest struct {
	CampaignID advertising.CampaignID `json:"campaign_id"`
}

type getAllAdvertisementsForCampaignResponse struct {
	Advertisements  []advertising.Advertisement `json:"advertisements"`
	CampaignMetrics advertising.CampaignMetrics `json:"campaign_metrics"`
}

func getAllAdvertisementsForCampaign(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getAllAdvertisementsForCampaignRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var advertisements []advertising.Advertisement
	var campaignMetrics *advertising.CampaignMetrics
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		advertisements, err = advertising.GetAllAdvertisementsForCampaignID(tx, req.CampaignID)
		if err != nil {
			return err
		}
		campaignMetrics, err = advertising.GetAdvertisementMetrics(tx, req.CampaignID)
		return err
	}); err != nil {
		return nil, err
	}
	return getAllAdvertisementsForCampaignResponse{
		Advertisements:  advertisements,
		CampaignMetrics: *campaignMetrics,
	}, nil
}

type insertAdvertisementRequest struct {
	LanguageCode string                 `json:"language_code"`
	CampaignID   advertising.CampaignID `json:"campaign_id"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	ImageURL     string                 `json:"image_url"`
}

type insertAdvertisementResponse struct {
	ID advertising.AdvertisementID `json:"id"`
}

func insertAdvertisement(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req insertAdvertisementRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return nil, err
	}
	var advertisementID *advertising.AdvertisementID
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		advertisementID, err = advertising.InsertNewAdvertisement(tx, advertising.InsertNewAdvertisementInput{
			LanguageCode: *languageCode,
			CampaignID:   req.CampaignID,
			Title:        req.Title,
			Description:  req.Description,
			ImageURL:     req.ImageURL,
		})
		return err
	}); err != nil {
		return nil, err
	}
	return insertAdvertisementResponse{
		ID: *advertisementID,
	}, nil
}

type updateAdvertisementRequest struct {
	ID           advertising.AdvertisementID `json:"id"`
	LanguageCode string                      `json:"language_code"`
	Title        string                      `json:"title"`
	Description  string                      `json:"description"`
	ImageURL     string                      `json:"image_url"`
	IsActive     bool                        `json:"is_active"`
}

type updateAdvertisementResponse struct {
	Success bool `json:"success"`
}

func updateAdvertisement(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req updateAdvertisementRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	languageCode, err := wordsmith.GetLanguageCodeFromString(req.LanguageCode)
	if err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return advertising.UpdateAdvertisement(tx, req.ID, advertising.UpdateAdvertisementInput{
			LanguageCode: *languageCode,
			Title:        req.Title,
			Description:  req.Description,
			ImageURL:     req.ImageURL,
			IsActive:     req.IsActive,
		})
	}); err != nil {
		return nil, err
	}
	return updateAdvertisementResponse{
		Success: true,
	}, nil
}
